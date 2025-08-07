package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/robertocorreajr/cfs_spool/internal/creality"
	"github.com/robertocorreajr/cfs_spool/internal/rfid"
)

// cmdReadTag lê uma tag RFID e decodifica o conteúdo
func cmdReadTag(args []string) {
	fs := flag.NewFlagSet("read-tag", flag.ExitOnError)
	var keyType, currentKey string
	var debug bool
	fs.StringVar(&keyType, "type", "A", "A|B da key para leitura (default A)")
	fs.StringVar(&currentKey, "key", "", "Key específica (12-hex). Se vazio, deriva do UID")
	fs.BoolVar(&debug, "debug", false, "Mostrar informações de debug")
	_ = fs.Parse(args)

	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║          LEITOR DE TAG CREALITY          ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// 1. Conectar ao leitor
	rdr, err := rfid.Open()
	if err != nil {
		fmt.Printf("Erro ao conectar leitor: %v\n", err)
		return
	}
	defer rdr.Close()

	// 2. Ler UID
	uid, err := rdr.UID()
	if err != nil {
		fmt.Printf("Erro ao ler UID: %v\n", err)
		return
	}
	fmt.Printf("UID: %s\n", uid)

	// 2.5. Diagnóstico (se habilitado)
	if debug {
		fmt.Println("\n--- Diagnóstico básico ---")
		fmt.Println("Tentando reconectar...")
		rdr.Close()
		rdr, err = rfid.Open()
		if err != nil {
			fmt.Printf("Erro ao reconectar: %v\n", err)
			return
		}
		defer rdr.Close()
		
		uid2, err := rdr.UID()
		if err != nil {
			fmt.Printf("Erro ao ler UID novamente: %v\n", err)
		} else {
			fmt.Printf("UID confirmado: %s\n", uid2)
		}
		
		fmt.Println("Testando comunicação básica...")
		err = rdr.TestBasicRead()
		if err != nil {
			fmt.Printf("Teste básico falhou: %v\n", err)
		}
	}

	// 3. Preparar lista de chaves para testar
	var testKeys []string
	
	if currentKey != "" {
		// Usar apenas a chave específica fornecida
		testKeys = []string{strings.ToUpper(currentKey)}
		fmt.Printf("Usando chave fornecida: %s\n", currentKey)
	} else {
		// Derivar chave S1 do UID e adicionar chaves comuns
		s1Key, err := creality.DeriveS1KeyFromUID(uid)
		if err != nil {
			fmt.Printf("Erro ao derivar chave S1: %v\n", err)
			return
		}
		fmt.Printf("Chave S1 derivada do UID: %s\n", s1Key)
		
		testKeys = []string{
			s1Key,            // Chave derivada (primeiro)
			"FFFFFFFFFFFF",   // Padrão MIFARE
			"000000000000",   // Zeros
			"A0A1A2A3A4A5",   // MAD key
			"B0B1B2B3B4B5",   // MAD key 2  
			"D3F7D3F7D3F7",   // NDEF key
			"714C5C886E97",   // Transport keys
			"587EE5F9350F",
			"A0478CC39091",
			"533CB6C723F6",
			"8FD0A4F256E9",
		}
	}

	// 4. Tentar ler blocos 4-6 com diferentes chaves
	fmt.Println("\n--- Tentando ler blocos 4-6 ---")
	
	var blocks []string
	var found bool
	var successKey string
	
	for i, testKey := range testKeys {
		if len(testKeys) > 1 {
			if i == 0 && currentKey == "" {
				fmt.Printf("Tentativa %d: Chave derivada (%s)\n", i+1, testKey)
			} else {
				fmt.Printf("Tentativa %d: %s\n", i+1, testKey)
			}
		}
		
		// Tentar com ambos os tipos de chave (A e B)
		for _, kt := range []byte{rfid.KeyTypeA, rfid.KeyTypeB} {
			keyName := map[byte]string{rfid.KeyTypeA: "A", rfid.KeyTypeB: "B"}[kt]
			
			// Primeiro tentar método padrão
			result, err := rdr.ReadRange(4, 3, kt, testKey)
			if err == nil {
				fmt.Printf("✅ SUCESSO! Chave %s tipo %s (método padrão)\n", testKey, keyName)
				blocks = result
				successKey = testKey
				found = true
				break
			}
			
			// Se falhou, tentar método alternativo
			result, err = rdr.ReadRangeAlternative(4, 3, kt, testKey)
			if err == nil {
				fmt.Printf("✅ SUCESSO! Chave %s tipo %s (método alternativo)\n", testKey, keyName)
				blocks = result
				successKey = testKey
				found = true
				break
			}
			
			if debug {
				fmt.Printf("   Ambos métodos falharam com tipo %s: %v\n", keyName, err)
			}
		}
		if found {
			break
		}
	}
	
	if !found {
		fmt.Println("❌ Nenhuma chave funcionou.")
		fmt.Println("\nSugestões:")
		fmt.Println("1. Verifique se a tag está corretamente posicionada no leitor")
		fmt.Println("2. Use Proxmark3 para força bruta: 'hf mf autopwn'")
		fmt.Println("3. Teste com o app Android 'RFID for CFS'")
		fmt.Println("4. Verifique se é realmente uma tag Creality/MIFARE Classic")
		return
	}

	// 5. Mostrar blocos lidos
	fmt.Println("\n--- Blocos lidos ---")
	for i, block := range blocks {
		fmt.Printf("Bloco %d: %s\n", 4+i, block)
	}

	// 6. Tentar decodificar como tag Creality
	fmt.Println("\n--- Decodificação Creality ---")
	payloadHex := strings.Join(blocks, "")
	fmt.Printf("Payload hex (96 chars): %s\n", payloadHex)
	
	ascii, err := creality.DecryptBlocks(payloadHex)
	if err != nil {
		fmt.Printf("Erro na descriptografia: %v\n", err)
		fmt.Printf("Dados brutos (hex): %s\n", payloadHex)
	} else {
		fmt.Printf("ASCII descriptografado: %s\n", ascii)
		
		fields, err := creality.ParseFields(ascii)
		if err != nil {
			fmt.Printf("Erro no parsing dos campos: %v\n", err)
			fmt.Printf("ASCII bruto: %s\n", ascii)
		} else {
			fmt.Println("\n╔══════════════════════════════════════════╗")
			fmt.Println("║           INFORMAÇÕES DA TAG             ║")
			fmt.Println("╚══════════════════════════════════════════╝")
			
			fmt.Printf("📦 Lote:        %s\n", fields.Batch)
			fmt.Printf("📅 Data:        %s\n", fields.FormatDate())
			fmt.Printf("🏭 Fornecedor:  %s\n", fields.Supplier)
			fmt.Printf("🧪 Material:    %s\n", fields.GetMaterialName())
			fmt.Printf("🎨 Cor:         %s\n", fields.FormatColor())
			fmt.Printf("📏 Comprimento: %s\n", fields.FormatLength())
			fmt.Printf("🔢 Serial:      %s\n", fields.Serial)
			
			if fields.Reserve != "00000000000000" {
				fmt.Printf("💾 Reservado:   %s\n", fields.Reserve)
			}
			
			// Mostrar dados técnicos apenas em modo debug
			if debug {
				fmt.Println("\n╔══════════════════════════════════════════╗")
				fmt.Println("║            DADOS TÉCNICOS                ║")
				fmt.Println("╚══════════════════════════════════════════╝")
				fmt.Printf("🆔 UID:         %s\n", uid)
				fmt.Printf("🔑 Chave usada: %s\n", strings.Join([]string{successKey[0:4], successKey[4:8], successKey[8:12]}, " "))
				fmt.Printf("📊 Payload:     %s\n", payloadHex)
			}
		}
	}

	fmt.Println("\n╔══════════════════════════════════════════╗")
	fmt.Println("║            LEITURA CONCLUÍDA             ║")
	fmt.Println("╚══════════════════════════════════════════╝")
}

// cmdWriteTag grava os dados (por enquanto só imprime o que faria)
func cmdWriteTag(args []string) {
	fs := flag.NewFlagSet("write-tag", flag.ExitOnError)

	// campos ASCII obrigatórios
	var batch, date, supplier, material, color, length, serial, reserve string
	fs.StringVar(&batch, "batch", "", "Batch (3)")
	fs.StringVar(&date, "date", "", "Date (5)")
	fs.StringVar(&supplier, "supplier", "", "Supplier (4)")
	fs.StringVar(&material, "material", "", "Material (5)")
	fs.StringVar(&color, "color", "", "Color (7)")
	fs.StringVar(&length, "length", "", "Length (4)")
	fs.StringVar(&serial, "serial", "", "Serial (6)")
	fs.StringVar(&reserve, "reserve", "", "Reserve (14)")

	// opções da key atual
	var keyType, currentKey string
	fs.StringVar(&keyType, "type", "B", "A|B da key atual (default B)")
	fs.StringVar(&currentKey, "currentkey", "FFFFFFFFFFFF", "Key atual (12-hex)")

	_ = fs.Parse(args)

	fields := creality.Fields{
		Batch:    batch,
		Date:     date,
		Supplier: supplier,
		Material: material,
		Color:    color,
		Length:   length,
		Serial:   serial,
		Reserve:  reserve,
	}
	payload, err := fields.ASCIIConcat()
	dieIf(err)

	fmt.Println("≡ (mock) programação de tag ≡")
	fmt.Printf("KeyType=%s  CurrentKey=%s\n", keyType, currentKey)
	fmt.Printf("ASCII payload: %s\n", payload)
	fmt.Println(">> TODO: abrir leitor, gravar blocos 4-7, validar leitura")
}
