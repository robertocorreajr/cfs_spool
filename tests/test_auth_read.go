package main

import (
	"fmt"
	"log"

	"github.com/robertocorreajr/cfs_spool/internal/rfid"
)

func main() {
	fmt.Println("=== Teste de Leitura com Autenticação ===")
	
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

	// Testar leitura com métodos de autenticação disponíveis
	fmt.Println("\n=== Teste com TryReadBlock ===")
	
	// Chaves padrão mais comuns
	defaultKeys := []string{
		"FFFFFFFFFFFF", // Chave padrão
		"000000000000", // Chave vazia
		"A0A1A2A3A4A5", // Chave comum
		"B0B1B2B3B4B5", // Chave comum
		"4D3A99C351DD", // Chave NFC padrão
		"1A982C7E459A", // Chave comum
	}

	// Obter chave derivada do UID
	derivedKey := rdr.DeriveKeyFromUID(uid)
	fmt.Printf("Chave derivada do UID: %s\n", derivedKey)
	
	// Adicionar chave derivada às chaves de teste
	testKeys := append([]string{derivedKey}, defaultKeys...)

	// Testar leitura dos blocos 4, 5, 6 (setor 1)
	for _, block := range []byte{4, 5, 6} {
		fmt.Printf("\n--- Bloco %d ---\n", block)
		
		success := false
		for keyIndex, key := range testKeys {
			fmt.Printf("Tentando chave %d (%s)...\n", keyIndex, key)
			
			// Tentar com Key Type A
			data, err := rdr.TryReadBlock(block, rfid.KeyTypeA, key)
			if err == nil {
				fmt.Printf("✅ Sucesso com Key A!\n")
				fmt.Printf("Dados: %s\n", data)
				success = true
				break
			}
			
			// Tentar com Key Type B
			data, err = rdr.TryReadBlock(block, rfid.KeyTypeB, key)
			if err == nil {
				fmt.Printf("✅ Sucesso com Key B!\n")
				fmt.Printf("Dados: %s\n", data)
				success = true
				break
			}
		}
		
		if !success {
			fmt.Printf("❌ Não foi possível ler bloco %d\n", block)
		}
	}
	
	// Testar leitura do trailer (bloco 7)
	fmt.Printf("\n--- Trailer (Bloco 7) ---\n")
	fmt.Println("Nota: O trailer pode não ser legível por questões de segurança")
	
	for keyIndex, key := range testKeys {
		fmt.Printf("Tentando chave %d (%s)...\n", keyIndex, key)
		
		data, err := rdr.TryReadBlock(7, rfid.KeyTypeA, key)
		if err == nil {
			fmt.Printf("✅ Trailer lido com Key A!\n")
			fmt.Printf("Dados: %s\n", data)
			break
		}
		
		data, err = rdr.TryReadBlock(7, rfid.KeyTypeB, key)
		if err == nil {
			fmt.Printf("✅ Trailer lido com Key B!\n")
			fmt.Printf("Dados: %s\n", data)
			break
		}
	}
}
