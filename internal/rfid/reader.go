package rfid

// Leitura/escrita de MIFARE Classic (PC/SC + ACR122U).
//
//   rdr, _ := Open()
//   defer rdr.Close()
//   uid := rdr.UID()
//   rdr.WriteBlock(4, keyTypeB, "FFFFFFFFFFFF", data32Hex)

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/ebfe/scard"
)

const (
	KeyTypeA = byte(0x60)
	KeyTypeB = byte(0x61)
)

// Reader mantém conexão PC/SC aberta.
type Reader struct {
	ctx  *scard.Context
	card *scard.Card
}

// Open conecta no 1º leitor encontrado (ACR122…).
func Open() (*Reader, error) {
	ctx, err := scard.EstablishContext()
	if err != nil {
		return nil, err
	}
	readers, err := ctx.ListReaders()
	if err != nil || len(readers) == 0 {
		return nil, errors.New("nenhum leitor PC/SC")
	}
	card, err := ctx.Connect(readers[0], scard.ShareShared, scard.ProtocolAny)
	if err != nil {
		return nil, err
	}
	return &Reader{ctx: ctx, card: card}, nil
}

func (r *Reader) Close() {
	if r.card != nil {
		r.card.Disconnect(scard.LeaveCard)
	}
	if r.ctx != nil {
		r.ctx.Release()
	}
}

// UID retorna 4 bytes em hex.
func (r *Reader) UID() (string, error) {
	// APDU “Get Data” para ACR122
	resp, err := r.transmit([]byte{0xFF, 0xCA, 0x00, 0x00, 0x00})
	if err != nil {
		return "", err
	}
	if len(resp) < 2 || resp[len(resp)-2] != 0x90 {
		return "", errors.New("falha UID")
	}
	return hex.EncodeToString(resp[:len(resp)-2]), nil
}

// Authenticate bloco com key (12 hex).
func (r *Reader) auth(block byte, keyType byte, keyHex string) error {
	key, _ := hex.DecodeString(keyHex)
	if len(key) != 6 {
		return errors.New("key deve ter 12 hex")
	}
	cmd := append([]byte{0xFF, 0x86, 0, 0, 5,
		1,       // version
		block,   // bloco
		keyType, // 0x60 A / 0x61 B
		0},      // key slot
		key...) // 6 bytes
	resp, err := r.transmit(cmd)
	if err != nil {
		return err
	}
	if len(resp) < 2 || resp[len(resp)-2] != 0x90 {
		return errors.New("authenticate falhou")
	}
	return nil
}

// WriteBlock grava bloco (4‐15…) com 32 hex (16 bytes).
func (r *Reader) WriteBlock(block byte, keyType byte, keyHex, dataHex string) error {
	if err := r.auth(block, keyType, keyHex); err != nil {
		return err
	}
	data, _ := hex.DecodeString(dataHex)
	if len(data) != 16 {
		return errors.New("bloco precisa de 32 hex")
	}
	cmd := append([]byte{0xFF, 0xD6, 0x00, block, 16}, data...)
	resp, err := r.transmit(cmd)
	if err != nil {
		return err
	}
	if resp[len(resp)-2] != 0x90 {
		return errors.New("write falhou")
	}
	return nil
}

// ReadRange lê n blocos consecutivos; devolve slice de 32-hex strings.
func (r *Reader) ReadRange(start byte, count int, keyType byte, keyHex string) ([]string, error) {
	out := make([]string, 0, count)
	for i := 0; i < count; i++ {
		blk := start + byte(i)
		if err := r.auth(blk, keyType, keyHex); err != nil {
			return nil, err
		}
		resp, err := r.transmit([]byte{0xFF, 0xB0, 0x00, blk, 16})
		if err != nil || resp[len(resp)-2] != 0x90 {
			return nil, fmt.Errorf("read %d falhou", blk)
		}
		out = append(out, strings.ToUpper(hex.EncodeToString(resp[:16])))
	}
	return out, nil
}

func (r *Reader) transmit(cmd []byte) ([]byte, error) {
	return r.card.Transmit(cmd)
}

// ReadBlockDirect lê um bloco sem autenticação (útil para bloco 0)
func (r *Reader) ReadBlockDirect(block byte) (string, error) {
	resp, err := r.transmit([]byte{0xFF, 0xB0, 0x00, block, 16})
	if err != nil {
		return "", err
	}
	if len(resp) < 18 || resp[len(resp)-2] != 0x90 {
		return "", fmt.Errorf("read direto bloco %d falhou", block)
	}
	return strings.ToUpper(hex.EncodeToString(resp[:16])), nil
}

// TestBasicRead tenta diferentes métodos de leitura para diagnóstico
func (r *Reader) TestBasicRead() error {
	// Método 1: Get Data (obtém UID/ATQA/SAK)
	resp, err := r.transmit([]byte{0xFF, 0xCA, 0x00, 0x00, 0x00})
	if err != nil {
		return fmt.Errorf("método 1 falhou: %v", err)
	}
	if len(resp) >= 2 && resp[len(resp)-2] == 0x90 {
		fmt.Printf("✓ Método 1 OK: %s\n", hex.EncodeToString(resp[:len(resp)-2]))
	}

	// Método 2: Load Key padrão
	defaultKey := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	cmd := append([]byte{0xFF, 0x82, 0x00, 0x00, 0x06}, defaultKey...)
	resp, err = r.transmit(cmd)
	if err == nil && len(resp) >= 2 && resp[len(resp)-2] == 0x90 {
		fmt.Printf("✓ Load Key OK\n")
		
		// Método 3: Authenticate bloco 4 com key A
		authCmd := []byte{0xFF, 0x86, 0x00, 0x00, 0x05, 0x01, 0x00, 0x04, 0x60, 0x00}
		resp, err = r.transmit(authCmd)
		if err == nil && len(resp) >= 2 && resp[len(resp)-2] == 0x90 {
			fmt.Printf("✓ Auth bloco 4 com key A padrão OK\n")
			return nil
		}
	}

	return fmt.Errorf("todos os métodos falharam")
}

// TryReadBlock tenta ler um bloco específico com diferentes abordagens
func (r *Reader) TryReadBlock(block byte, keyType byte, keyHex string) (string, error) {
	// Método 1: Autenticação + leitura normal
	if err := r.auth(block, keyType, keyHex); err == nil {
		resp, err := r.transmit([]byte{0xFF, 0xB0, 0x00, block, 16})
		if err == nil && len(resp) >= 18 && resp[len(resp)-2] == 0x90 {
			return strings.ToUpper(hex.EncodeToString(resp[:16])), nil
		}
	}

	// Método 2: Load key específica + authenticate + read
	keyBytes, _ := hex.DecodeString(keyHex)
	if len(keyBytes) == 6 {
		// Load key no slot 0
		cmd := append([]byte{0xFF, 0x82, 0x00, 0x00, 0x06}, keyBytes...)
		resp, err := r.transmit(cmd)
		if err == nil && len(resp) >= 2 && resp[len(resp)-2] == 0x90 {
			// Authenticate usando key slot 0
			authCmd := []byte{0xFF, 0x86, 0x00, 0x00, 0x05, 0x01, 0x00, block, keyType, 0x00}
			resp, err = r.transmit(authCmd)
			if err == nil && len(resp) >= 2 && resp[len(resp)-2] == 0x90 {
				// Read block
				resp, err = r.transmit([]byte{0xFF, 0xB0, 0x00, block, 16})
				if err == nil && len(resp) >= 18 && resp[len(resp)-2] == 0x90 {
					return strings.ToUpper(hex.EncodeToString(resp[:16])), nil
				}
			}
		}
	}

	return "", fmt.Errorf("não foi possível ler bloco %d", block)
}

// ReadRangeAlternative versão alternativa que tenta diferentes métodos
func (r *Reader) ReadRangeAlternative(start byte, count int, keyType byte, keyHex string) ([]string, error) {
	var blocks []string
	
	// Tentar ler bloco por bloco com diferentes métodos
	for i := 0; i < count; i++ {
		block := start + byte(i)
		data, err := r.TryReadBlock(block, keyType, keyHex)
		if err != nil {
			return nil, fmt.Errorf("falha no bloco %d: %v", block, err)
		}
		blocks = append(blocks, data)
	}
	
	return blocks, nil
}
