package main

import (
	"fmt"
	"log"

	"github.com/robertocorreajr/cfs_spool/internal/rfid"
	"github.com/robertocorreajr/cfs_spool/internal/creality"
)

func main() {
	fmt.Println("=== Teste de Leitura e Decodificação CFS ===")
	
	// Conectar ao leitor
	rdr, err := rfid.Open()
	if err != nil {
		log.Fatalf("Erro ao conectar: %v", err)
	}
	defer rdr.Close()

	// Ler UID
	uid, err := rdr.UID()
	if err != nil {
		log.Fatalf("Erro ao ler UID: %v", err)
	}
	fmt.Printf("UID: %s\n", uid)

	// Obter chave derivada
	derivedKey := rdr.DeriveKeyFromUID(uid)
	fmt.Printf("Chave derivada: %s\n", derivedKey)

	// Ler blocos 4, 5, 6
	fmt.Println("\n=== Lendo Dados Brutos ===")
	var blocks []string
	
	for _, block := range []byte{4, 5, 6} {
		data, err := rdr.TryReadBlock(block, rfid.KeyTypeA, derivedKey)
		if err != nil {
			log.Fatalf("Erro ao ler bloco %d: %v", block, err)
		}
		fmt.Printf("Bloco %d: %s\n", block, data)
		blocks = append(blocks, data)
	}

	// Decodificar dados CFS
	fmt.Println("\n=== Decodificando Dados CFS ===")
	
	// Concatenar blocos
	allData := blocks[0] + blocks[1] + blocks[2]
	fmt.Printf("Dados concatenados: %s\n", allData)
	
	// Tentar descriptografar
	decrypted, err := creality.DecryptBlocks(allData)
	if err != nil {
		fmt.Printf("Erro na descriptografia: %v\n", err)
		fmt.Println("Tentando interpretar dados sem descriptografia...")
		decrypted = allData
	} else {
		fmt.Printf("Dados descriptografados: %s\n", decrypted)
	}
	
	// Converter hex para ASCII se necessário
	ascii38 := decrypted
	if len(ascii38) == 76 { // Se ainda é hex, converter para ASCII (38*2=76)
		ascii38, err = hexToASCII(decrypted)
		if err != nil {
			fmt.Printf("Erro ao converter hex para ASCII: %v\n", err)
			return
		}
		fmt.Printf("ASCII convertido: %s\n", ascii38)
	}
	
	// Parse dos campos (novo formato de 38 bytes)
	fmt.Println("\n=== Interpretação dos Campos ===")
	
	if len(ascii38) >= 38 {
		fields, err := creality.ParseFields(ascii38[:38])
		if err != nil {
			fmt.Printf("Erro ao fazer parse: %v\n", err)
		} else {
			fmt.Printf("Lote (fixo): %s\n", fields.Batch)
			fmt.Printf("Data: %s (%s)\n", fields.Date, fields.FormatDate())
			fmt.Printf("Fornecedor: %s (%s)\n", fields.Supplier, fields.GetSupplierName())
			fmt.Printf("Material: %s (%s)\n", fields.Material, fields.GetMaterialName())
			fmt.Printf("Cor: %s (%s)\n", fields.Color, fields.FormatColor())
			fmt.Printf("Comprimento: %s (%s)\n", fields.Length, fields.FormatLength())
			fmt.Printf("Serial: %s\n", fields.Serial)
			fmt.Printf("Reserva (fixo): %s\n", fields.Reserve)
		}
	}
	
	// Mostrar dados hex em formato legível
	fmt.Println("\n=== Dados Hex Formatados ===")
	for i := 0; i < len(allData); i += 32 {
		end := i + 32
		if end > len(allData) {
			end = len(allData)
		}
		line := allData[i:end]
		fmt.Printf("Offset %02X: %s\n", i/2, line)
	}
}

// hexToASCII converte string hex para ASCII
func hexToASCII(hexStr string) (string, error) {
	if len(hexStr)%2 != 0 {
		return "", fmt.Errorf("string hex deve ter número par de caracteres")
	}
	
	var result []byte
	for i := 0; i < len(hexStr); i += 2 {
		b := 0
		_, err := fmt.Sscanf(hexStr[i:i+2], "%02x", &b)
		if err != nil {
			return "", err
		}
		result = append(result, byte(b))
	}
	
	return string(result), nil
}
