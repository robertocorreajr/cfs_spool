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
	"github.com/robertocorreajr/cfs_spool/internal/creality"
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

// WriteBlockAlternative tenta escrever um bloco com métodos alternativos
func (r *Reader) WriteBlockAlternative(block byte, keyType byte, keyHex, dataHex string) error {
	data, err := hex.DecodeString(dataHex)
	if err != nil || len(data) != 16 {
		return errors.New("bloco precisa de 32 hex válidos")
	}
	
	// Método 1: Autenticação + escrita normal
	if err := r.auth(block, keyType, keyHex); err == nil {
		cmd := append([]byte{0xFF, 0xD6, 0x00, block, 16}, data...)
		resp, err := r.transmit(cmd)
		if err == nil && len(resp) >= 2 && resp[len(resp)-2] == 0x90 {
			return nil
		}
	}

	// Método 2: Load key + authenticate + write
	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil || len(keyBytes) != 6 {
		return errors.New("key deve ter 12 hex")
	}
	
	// Load key no slot 0
	cmd := append([]byte{0xFF, 0x82, 0x00, 0x00, 0x06}, keyBytes...)
	resp, err := r.transmit(cmd)
	if err != nil || len(resp) < 2 || resp[len(resp)-2] != 0x90 {
		return fmt.Errorf("falha ao carregar key: %v", err)
	}
	
	// Authenticate usando key slot 0
	authCmd := []byte{0xFF, 0x86, 0x00, 0x00, 0x05, 0x01, 0x00, block, keyType, 0x00}
	resp, err = r.transmit(authCmd)
	if err != nil || len(resp) < 2 || resp[len(resp)-2] != 0x90 {
		return fmt.Errorf("falha na autenticação: %v", err)
	}
	
	// Write block
	cmd = append([]byte{0xFF, 0xD6, 0x00, block, 16}, data...)
	resp, err = r.transmit(cmd)
	if err != nil || len(resp) < 2 || resp[len(resp)-2] != 0x90 {
		return fmt.Errorf("falha na escrita: %v", err)
	}
	
	return nil
}

// WriteRange escreve múltiplos blocos consecutivos
func (r *Reader) WriteRange(start byte, blocks []string, keyType byte, keyHex string) error {
	for i, blockData := range blocks {
		block := start + byte(i)
		
		// Tentar primeiro método normal, depois alternativo
		err := r.WriteBlock(block, keyType, keyHex, blockData)
		if err != nil {
			err = r.WriteBlockAlternative(block, keyType, keyHex, blockData)
			if err != nil {
				return fmt.Errorf("falha ao escrever bloco %d: %v", block, err)
			}
		}
	}
	return nil
}

// WriteTagCFS escreve dados CFS nos blocos 4, 5, 6 usando o padrão JavaScript
func (r *Reader) WriteTagCFS(uid string, blocksToWrite []string, encrypted bool) error {
	// Primeiro tentar determinar se é tag nova ou usada
	// Testar autenticação com FFFFFFFFFFFF no setor 1
	var key string
	var isNewTag bool
	
	// Derivar a chave do UID primeiro
	derivedKey := r.DeriveKeyFromUID(uid)
	fmt.Printf("🔑 Chave derivada do UID %s: %s\n", uid, derivedKey)

	// Primeiro tentar autenticar com a chave derivada
	err := r.testAuthentication(4, derivedKey)
	if err == nil {
		key = derivedKey
		isNewTag = false
		fmt.Printf("🔄 Tag detectada como USADA (usando key derivada: %s)\n", key)
	} else {
		// Se falhar, tentar com a chave padrão
		err = r.testAuthentication(4, "FFFFFFFFFFFF")
		if err == nil {
			key = "FFFFFFFFFFFF"
			isNewTag = true
			fmt.Println("🆕 Tag detectada como NOVA (usando FFFFFFFFFFFF)")
		} else {
			// Se ambas falharem, usar a chave derivada mesmo assim
			// porque a tag pode estar usando uma chave antiga
			key = derivedKey
			isNewTag = false
			fmt.Printf("⚠️ Autenticação falhou, usando key derivada: %s\n", key)
		}
	}
	
	// Escrever blocos 4, 5, 6
	blocks := []byte{4, 5, 6}
	for i, blockNum := range blocks {
		if i >= len(blocksToWrite) {
			break
		}
		
		err := r.WriteBlockDirectly(blockNum, key, blocksToWrite[i], uid)
		if err != nil {
			return fmt.Errorf("erro ao escrever bloco %d: %v", blockNum, err)
		}
		
		fmt.Printf("✅ Bloco %d escrito com sucesso\n", blockNum)
	}
	
	// Para tags novas, atualizar o trailer (bloco 7) com key derivada
	// IMPORTANTE: A impressora Creality só reconhece tags com key derivada no trailer
	if isNewTag {
		derivedKey := r.DeriveKeyFromUID(uid)
		fmt.Printf("🔑 Atualizando trailer para compatibilidade Creality (key: %s)\n", derivedKey)
		
		// Access bits FF0780: permite leitura e escrita com KeyA ou KeyB
		// GPB 69: padrão Creality
		trailer := derivedKey + "FF078069" + derivedKey // KeyA + Access + GPB + KeyB

		fmt.Printf("🔑 Trailer que será gravado: %s\n", trailer)
		
		err := r.WriteBlockDirectly(7, key, trailer, uid) // Usar key atual (FFFFFFFFFFFF) para escrever
		if err != nil {
			return fmt.Errorf("erro ao escrever trailer: %v", err)
		}
		fmt.Println("✅ Trailer atualizado - tag compatível com impressora Creality")
	}
	
	return nil
}

// WriteBlockDirectly escreve um bloco usando Load Key + Authenticate + Write
func (r *Reader) WriteBlockDirectly(block byte, keyHex, dataHex string, uid ...string) error {
	// Re-selecionar o cartão para garantir estado limpo antes de autenticar
	r.UID()

	// Montar lista de chaves para tentar (sem duplicatas)
	keys := []string{keyHex}
	seen := map[string]bool{keyHex: true}

	// Derivar chave do UID se disponível (usa uid passado pelo chamador)
	uidStr := ""
	if len(uid) > 0 {
		uidStr = uid[0]
	}
	if uidStr != "" {
		derivedKey := r.DeriveKeyFromUID(uidStr)
		if !seen[derivedKey] {
			keys = append(keys, derivedKey)
			seen[derivedKey] = true
		}
	}

	// Chave padrão como fallback
	if !seen["FFFFFFFFFFFF"] {
		keys = append(keys, "FFFFFFFFFFFF")
	}

	var lastErr error
	for _, key := range keys {
		// 1. Load Key no slot 0
		keyBytes, err := hex.DecodeString(key)
		if err != nil || len(keyBytes) != 6 {
			lastErr = errors.New("key deve ter 12 hex chars válidos")
			continue
		}
		
		cmd := append([]byte{0xFF, 0x82, 0x00, 0x00, 0x06}, keyBytes...)
		resp, err := r.transmit(cmd)
		if err != nil || len(resp) < 2 || resp[len(resp)-2] != 0x90 {
			lastErr = fmt.Errorf("falha ao carregar key para bloco %d", block)
			continue
		}
		
		// 2. Tentar authenticate + write com cada KeyType
		// KeyB primeiro: access bits podem exigir KeyB para escrita
		var written bool
		for _, keyType := range []byte{KeyTypeB, KeyTypeA} {
			authCmd := []byte{0xFF, 0x86, 0x00, 0x00, 0x05, 0x01, 0x00, block, keyType, 0x00}
			resp, err = r.transmit(authCmd)
			if err != nil || len(resp) < 2 || resp[len(resp)-2] != 0x90 {
				// Auth falhou — re-selecionar cartão e recarregar key antes de tentar próximo tipo
				r.UID()
				reloadCmd := append([]byte{0xFF, 0x82, 0x00, 0x00, 0x06}, keyBytes...)
				r.transmit(reloadCmd)
				continue
			}

			// Auth OK — tentar escrever
			data, err := hex.DecodeString(dataHex)
			if err != nil || len(data) != 16 {
				return errors.New("dados devem ter 32 hex chars")
			}

			writeCmd := append([]byte{0xFF, 0xD6, 0x00, block, 16}, data...)
			resp, err = r.transmit(writeCmd)
			if err == nil && len(resp) >= 2 && resp[len(resp)-2] == 0x90 {
				written = true
				break
			}
			// Escrita falhou — re-selecionar e recarregar key antes de tentar próximo tipo
			r.UID()
			reloadCmd := append([]byte{0xFF, 0x82, 0x00, 0x00, 0x06}, keyBytes...)
			r.transmit(reloadCmd)
		}

		if written {
			return nil
		}
		lastErr = fmt.Errorf("falha na escrita do bloco %d", block)
	}
	
	return lastErr // Retorna o último erro se todas as tentativas falharam
}

// DeriveKeyFromUID deriva a chave do UID (implementação baseada no JS)
func DeriveKeyFromUID(uid string) (string, error) {
	if len(uid) != 8 {
		return "", errors.New("UID deve ter 8 hex chars")
	}
	
	// Esta função precisa ser importada do pacote creality
	// Por agora, retornamos a chave fixa
	return "FFFFFFFFFFFF", nil
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

// testAuthentication testa se uma key funciona para um bloco específico
func (r *Reader) testAuthentication(block byte, keyHex string) error {
	// 1. Load Key no slot 0
	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil || len(keyBytes) != 6 {
		return errors.New("key inválida")
	}
	
	cmd := append([]byte{0xFF, 0x82, 0x00, 0x00, 0x06}, keyBytes...)
	resp, err := r.transmit(cmd)
	if err != nil || len(resp) < 2 || resp[len(resp)-2] != 0x90 {
		return errors.New("falha ao carregar key")
	}
	
	// 2. Authenticate bloco
	authCmd := []byte{0xFF, 0x86, 0x00, 0x00, 0x05, 0x01, 0x00, block, KeyTypeA, 0x00}
	resp, err = r.transmit(authCmd)
	if err != nil || len(resp) < 2 || resp[len(resp)-2] != 0x90 {
		return errors.New("falha na autenticação")
	}
	
	return nil
}

// DeriveKeyFromUID deriva uma key do UID usando o algoritmo Creality
func (r *Reader) DeriveKeyFromUID(uid string) string {
	derivedKey, err := creality.DeriveS1KeyFromUID(uid)
	if err != nil {
		// Em caso de erro, retornar a chave padrão
		return "FFFFFFFFFFFF"
	}
	return derivedKey
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
