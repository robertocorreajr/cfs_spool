package main

import (
	"fmt"
	"log"
	"github.com/robertocorreajr/cfs_spool/internal/rfid"
	"github.com/robertocorreajr/cfs_spool/internal/creality"
)

func main() {
	fmt.Println("=== TESTE DE LEITURA PARA DIAGNÓSTICO ===")
	
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
	fmt.Printf("🆔 UID da tag: %s\n", uid)

	// Derivar chave do UID
	derivedKey, err := creality.DeriveS1KeyFromUID(uid)
	if err != nil {
		log.Fatalf("Erro ao derivar chave: %v", err)
	}
	fmt.Printf("🔑 Chave derivada: %s\n", derivedKey)

	// Lista de chaves para testar
	testKeys := []string{
		"FFFFFFFFFFFF", // Chave padrão
		derivedKey,     // Chave derivada do UID
		"000000000000", // Chave zero
	}

	fmt.Println("\n🔍 Testando leitura dos blocos de dados...")
	
	// Testar leitura de cada bloco (4, 5, 6) com diferentes chaves
	for _, block := range []byte{4, 5, 6} {
		fmt.Printf("\n📖 Testando bloco %d:\n", block)
		
		found := false
		for _, key := range testKeys {
			fmt.Printf("   Tentando chave %s... ", key)
			
			data, err := rdr.TryReadBlock(block, rfid.KeyTypeA, key)
			if err != nil {
				fmt.Printf("❌ Falhou: %v\n", err)
				continue
			}
			
			fmt.Printf("✅ Sucesso: %s\n", data)
			found = true
			break
		}
		
		if !found {
			fmt.Printf("   ⚠️  Nenhuma chave funcionou para o bloco %d\n", block)
		}
	}

	// Tentar ler e descriptografar dados completos
	fmt.Println("\n🔐 Tentando ler dados completos da tag...")
	
	var blocks []string
	var readSuccess bool

	// Tentar com each key sequencialmente para todos os blocos
	for _, key := range testKeys {
		fmt.Printf("Tentando ler todos os blocos com chave %s...\n", key)
		blocks = nil
		allBlocksRead := true
		
		for block := byte(4); block <= 6; block++ {
			data, err := rdr.TryReadBlock(block, rfid.KeyTypeA, key)
			if err != nil {
				fmt.Printf("   Bloco %d falhou: %v\n", block, err)
				allBlocksRead = false
				break
			}
			blocks = append(blocks, data)
		}
		
		if allBlocksRead {
			fmt.Printf("✅ Todos os blocos lidos com chave %s!\n", key)
			readSuccess = true
			break
		}
	}

	if !readSuccess {
		fmt.Println("❌ Não foi possível ler todos os blocos da tag")
		return
	}

	// Tentar descriptografar
	fmt.Println("\n🔓 Tentando descriptografar dados...")
	hexData := blocks[0] + blocks[1] + blocks[2]
	fmt.Printf("Dados hex concatenados: %s\n", hexData)
	
	decrypted, err := creality.DecryptBlocks(hexData)
	if err != nil {
		fmt.Printf("❌ Erro na descriptografia: %v\n", err)
		fmt.Println("   (Isso pode indicar que a tag não contém dados CFS válidos)")
	} else {
		fmt.Printf("✅ Dados descriptografados: %s\n", decrypted)
		
		// Tentar parser os campos
		fields, err := creality.ParseFields(decrypted)
		if err != nil {
			fmt.Printf("⚠️  Dados descriptografados mas parser falhou: %v\n", err)
		} else {
			fmt.Println("✅ Campos parseados com sucesso:")
			fmt.Printf("   Batch: %s\n", fields.Batch)
			fmt.Printf("   Date: %s (%s)\n", fields.Date, fields.FormatDate())
			fmt.Printf("   Supplier: %s (%s)\n", fields.Supplier, fields.GetSupplierName())
			fmt.Printf("   Material: %s (%s)\n", fields.Material, fields.GetMaterialName())
			fmt.Printf("   Color: %s (%s)\n", fields.Color, fields.FormatColor())
			fmt.Printf("   Length: %s (%s)\n", fields.Length, fields.FormatLength())
			fmt.Printf("   Serial: %s\n", fields.Serial)
		}
	}

	fmt.Println("\n=== DIAGNÓSTICO COMPLETO ===")
}
