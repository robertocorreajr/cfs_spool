package creality

import (
	"crypto/aes"
	"encoding/hex"
	"errors"
	"strings"
)

// Chaves fixas (mesmas do script JS)
const (
	keyUID     = "q3bu^t1nqfZ(pf$1"
	keyPayload = "H@CFkRnz@KAtBJp2"
)

// DeriveS1KeyFromUID gera a key A/B do setor 1.
func DeriveS1KeyFromUID(uidHex string) (string, error) {
	uidHex = strings.ToUpper(uidHex)
	if len(uidHex) != 8 {
		return "", errors.New("UID deve ter 4 bytes (8 hex)")
	}
	// repete o UID 4× → 16 bytes
	repeated, _ := hex.DecodeString(uidHex + uidHex + uidHex + uidHex)
	block, _ := aes.NewCipher([]byte(keyUID))
	out := make([]byte, 16)
	block.Encrypt(out, repeated)                              // AES-ECB (1 bloco)
	return strings.ToUpper(hex.EncodeToString(out)[:12]), nil // 6 bytes
}

func EncryptPayloadToBlocks(ascii48 string) (b4, b5, b6 string, err error) {
	// 48 ASCII bytes → precisa ser múltiplo de 16
	if len(ascii48) != 48 {
		return "", "", "", errors.New("payload ASCII deve ter 48 bytes")
	}

	key := []byte(keyPayload) // 16 bytes
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", "", err
	}

	plain := []byte(ascii48)
	out := make([]byte, len(plain))

	// ECB manual: cifra bloco-a-bloco (NoPadding)
	for i := 0; i < len(plain); i += 16 {
		block.Encrypt(out[i:i+16], plain[i:i+16])
	}

	hexStr := strings.ToUpper(hex.EncodeToString(out)) // 96 hex

	// split em 3 blocos de 32 hex
	b4 = hexStr[0:32]
	b5 = hexStr[32:64]
	b6 = hexStr[64:96]
	return
}

// DecryptBlocks descriptografa os blocos concatenados (96 hex chars) para ASCII
func DecryptBlocks(hexPayload string) (string, error) {
	hexPayload = strings.ToUpper(hexPayload)
	if len(hexPayload) != 96 {
		return "", errors.New("payload hex deve ter 96 chars (48 bytes)")
	}

	// Decodificar hex para bytes
	cipherBytes, err := hex.DecodeString(hexPayload)
	if err != nil {
		return "", err
	}

	key := []byte(keyPayload) // 16 bytes
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	out := make([]byte, len(cipherBytes))

	// ECB manual: descriptografa bloco-a-bloco (NoPadding)
	for i := 0; i < len(cipherBytes); i += 16 {
		block.Decrypt(out[i:i+16], cipherBytes[i:i+16])
	}

	return string(out), nil
}
