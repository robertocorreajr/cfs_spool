package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

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

// cmdWriteTag grava os dados na tag RFID
func cmdWriteTag(args []string) {
	fs := flag.NewFlagSet("write-tag", flag.ExitOnError)

	// campos ASCII obrigatórios
	var batch, date, supplier, material, color, length, serial, reserve string
	fs.StringVar(&batch, "batch", "", "Batch (3 chars): Exemplo: 1A5")
	fs.StringVar(&date, "date", "", "Date (5 chars): Formato YYMDD, exemplo: 24120 para Janeiro 2024")  
	fs.StringVar(&supplier, "supplier", "", "Supplier (4 chars): 0276=Creality, 0000=Genérico")
	fs.StringVar(&material, "material", "", "Material (5 chars): Exemplo: 04001 para CR-PLA")
	fs.StringVar(&color, "color", "", "Color (7 chars): Exemplo: 077BB41 para verde")
	fs.StringVar(&length, "length", "", "Length (4 chars): Exemplo: 0330 para 330cm")
	fs.StringVar(&serial, "serial", "", "Serial (6 chars): Exemplo: 000001")
	fs.StringVar(&reserve, "reserve", "00000000000000", "Reserve (14 chars): Padrão são zeros")

	// opções da key atual (key que está na tag agora)
	var keyType, currentKey string
	var newKey, newKeyType string
	var debug, verify bool
	fs.StringVar(&keyType, "current-key-type", "A", "Tipo da key atual: A ou B")
	fs.StringVar(&currentKey, "current-key", "", "Key atual da tag (12-hex). Se vazio, tenta derivar do UID")
	fs.StringVar(&newKeyType, "new-key-type", "", "Tipo da nova key (A ou B). Se vazio, mantém a atual")
	fs.StringVar(&newKey, "new-key", "", "Nova key (12-hex). Se vazio, mantém a atual")
	fs.BoolVar(&debug, "debug", false, "Mostrar informações de debug")
	fs.BoolVar(&verify, "verify", true, "Verificar escrita lendo os blocos novamente")

	_ = fs.Parse(args)

	// Validar campos obrigatórios
	if batch == "" || date == "" || supplier == "" || material == "" || 
		color == "" || length == "" || serial == "" {
		fmt.Println("❌ Erro: Todos os campos são obrigatórios!")
		fmt.Println("\nExemplo de uso:")
		fmt.Println("./cfs-spool write-tag \\")
		fmt.Println("  -batch=A20 \\")
		fmt.Println("  -date=25158 \\")  
		fmt.Println("  -supplier=0276 \\")
		fmt.Println("  -material=04001 \\")
		fmt.Println("  -color=077BB41 \\")
		fmt.Println("  -length=0330 \\")
		fmt.Println("  -serial=000001")
		return
	}

	// Validar código do fornecedor
	validSuppliers := map[string]string{
		"0276": "Creality",
		"0000": "Genérico",
	}
	if _, isValid := validSuppliers[supplier]; !isValid {
		fmt.Printf("❌ Erro: Código de fornecedor '%s' é inválido!\n", supplier)
		fmt.Println("\nCódigos válidos:")
		fmt.Println("  • 0276 = Creality")
		fmt.Println("  • 0000 = Genérico")
		return
	}

	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║         GRAVADOR DE TAG CREALITY        ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// 1. Conectar ao leitor
	rdr, err := rfid.Open()
	if err != nil {
		fmt.Printf("❌ Erro ao conectar leitor: %v\n", err)
		return
	}
	defer rdr.Close()

	// 2. Ler UID
	uid, err := rdr.UID()
	if err != nil {
		fmt.Printf("❌ Erro ao ler UID: %v\n", err)
		return
	}
	fmt.Printf("🆔 UID da tag: %s\n", uid)

	// 3. Determinar key atual
	var useKey string
	if currentKey != "" {
		useKey = strings.ToUpper(currentKey)
		fmt.Printf("🔑 Usando key fornecida: %s\n", useKey)
	} else {
		// Derivar key do UID
		derivedKey, err := creality.DeriveS1KeyFromUID(uid)
		if err != nil {
			fmt.Printf("❌ Erro ao derivar key do UID: %v\n", err)
			return
		}
		useKey = derivedKey
		fmt.Printf("🔑 Key derivada do UID: %s\n", useKey)
	}

	// 4. Criar estrutura de campos
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

	// 5. Validar tamanhos dos campos
	payload, err := fields.ASCIIConcat()
	if err != nil {
		fmt.Printf("❌ Erro na validação dos campos: %v\n", err)
		fmt.Println("\n📋 Tamanhos corretos:")
		fmt.Println("  • Batch: 3 chars")
		fmt.Println("  • Date: 5 chars (YYMDD)")  
		fmt.Println("  • Supplier: 4 chars")
		fmt.Println("  • Material: 5 chars")
		fmt.Println("  • Color: 7 chars")
		fmt.Println("  • Length: 4 chars")
		fmt.Println("  • Serial: 6 chars")
		fmt.Println("  • Reserve: 14 chars")
		return
	}

	if debug {
		fmt.Printf("\n📋 Payload ASCII (38 bytes): %s\n", payload)
	}

	// 6. Criptografar dados
	fmt.Println("🔐 Criptografando dados...")
	b4, b5, b6, err := creality.EncryptPayloadToBlocks(payload)
	if err != nil {
		fmt.Printf("❌ Erro na criptografia: %v\n", err)
		return
	}

	fmt.Printf("\n📊 DADOS RAW QUE SERÃO GRAVADOS:\n")
	fmt.Printf("  Payload ASCII: %s (%d bytes)\n", payload, len(payload))
	fmt.Printf("  Bloco 4 (hex): %s\n", b4)
	fmt.Printf("  Bloco 5 (hex): %s\n", b5)
	fmt.Printf("  Bloco 6 (hex): %s\n", b6)

	if debug {
		fmt.Printf("📊 Bloco 4: %s\n", b4)
		fmt.Printf("📊 Bloco 5: %s\n", b5)
		fmt.Printf("📊 Bloco 6: %s\n", b6)
	}

	// 7. Mostrar prévia dos dados  
	fmt.Println("\n╔══════════════════════════════════════════╗")
	fmt.Println("║          PRÉVIA DOS DADOS                ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Printf("📦 Lote:        %s\n", fields.Batch)
	fmt.Printf("📅 Data:        %s\n", fields.FormatDate())
	fmt.Printf("🏭 Fornecedor:  %s\n", fields.Supplier)
	fmt.Printf("🧪 Material:    %s\n", fields.GetMaterialName())
	fmt.Printf("🎨 Cor:         %s\n", fields.FormatColor())
	fmt.Printf("📏 Comprimento: %s\n", fields.FormatLength())
	fmt.Printf("🔢 Serial:      %s\n", fields.Serial)

	// 8. Confirmação do usuário
	fmt.Println("\n⚠️  ATENÇÃO: Esta operação irá SOBRESCREVER os dados da tag!")
	fmt.Print("Deseja continuar? (S/n): ")
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "" && strings.ToLower(confirm) != "s" && strings.ToLower(confirm) != "sim" {
		fmt.Println("❌ Operação cancelada pelo usuário.")
		return
	}

	// 9. Escrever blocos usando método JavaScript
	fmt.Println("\n🔧 Gravando dados na tag...")
	
	// Preparar lista de blocos
	blocksToWrite := []string{b4, b5, b6}
	
	// Escrever como tag não criptografada primeiro (chave padrão)
	err = rdr.WriteTagCFS(uid, blocksToWrite, false)
	if err != nil {
		fmt.Printf("❌ Erro na escrita: %v\n", err)
		fmt.Println("\n💡 Dicas:")
		fmt.Println("  • Verifique se a tag está corretamente posicionada")
		fmt.Println("  • Confirme se é uma tag MIFARE Classic")
		fmt.Println("  • Tente com uma tag nova/zerada")
		return
	}

	fmt.Println("✅ Dados gravados com sucesso!")

	// 10. Verificação (se habilitada)
	if verify {
		fmt.Println("\n🔍 Verificando escrita...")
		
		// Aguardar um pouco para a tag processar
		time.Sleep(time.Millisecond * 100)
		
		// Tentar ler com a mesma chave que foi usada para gravar
		success := true
		for i, expectedBlock := range blocksToWrite {
			block := byte(4 + i)
			readData, err := rdr.TryReadBlock(block, rfid.KeyTypeA, useKey)
			if err != nil {
				fmt.Printf("❌ Verificação falhou no bloco %d: %v\n", block, err)
				success = false
			} else if strings.EqualFold(readData, expectedBlock) {
				if debug {
					fmt.Printf("✅ Bloco %d verificado\n", block)
				}
			} else {
				fmt.Printf("❌ Dados não conferem no bloco %d!\n", block)
				fmt.Printf("    Esperado: %s\n", expectedBlock)
				fmt.Printf("    Lido:     %s\n", readData)
				success = false
			}
		}
		
		if success {
			fmt.Println("✅ Verificação: Todos os dados foram gravados corretamente!")
		} else {
			fmt.Println("⚠️  Alguns blocos podem não ter sido gravados corretamente.")
		}
	}

	// 11. Atualizar keys se solicitado  
	if newKey != "" {
		fmt.Println("\n🔑 Atualizando keys de acesso...")
		
		var newKeyTypeByte byte = rfid.KeyTypeA // Por padrão usa tipo A
		if newKeyType != "" {
			if strings.ToUpper(newKeyType) == "B" {
				newKeyTypeByte = rfid.KeyTypeB
			} else {
				newKeyTypeByte = rfid.KeyTypeA
			}
		}
		
		// TODO: Implementar escrita no bloco 7 (trailer do setor)
		// Bloco 7 contém: KeyA(6) + AccessBits(3) + GPB(1) + KeyB(6)
		fmt.Printf("⚠️  Atualização de keys ainda não implementada\n")
		fmt.Printf("    Nova key seria: %s (tipo %s)\n", newKey, string(rune(newKeyTypeByte)))
	}

	fmt.Println("\n╔══════════════════════════════════════════╗")
	fmt.Println("║           GRAVAÇÃO CONCLUÍDA             ║")
	fmt.Println("╚══════════════════════════════════════════╝")
}
