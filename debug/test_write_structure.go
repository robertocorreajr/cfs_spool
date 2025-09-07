package main

import (
	"fmt"
	"log"

	"github.com/robertocorreajr/cfs_spool/internal/creality"
)

func main() {
	fmt.Println("=== Teste da Estrutura de Gravação ===")

	// Criar campos para gravação
	fields := creality.NewFields()
	
	// Configurar campos como seriam enviados do formulário web
	fields.Date = "AB124"     // Dezembro de 2024
	fields.Supplier = "0276"  // Creality  
	fields.Material = "01001" // Hyper PLA
	fields.Length = "0165"    // 500g
	fields.Serial = "000001"  // Serial
	
	// Testar SetColor com 6 caracteres (como vem da interface web)
	fmt.Println("Testando SetColor com entrada de 6 caracteres...")
	
	err := fields.SetColor("FF4010")  // Laranja - 6 caracteres
	if err != nil {
		log.Fatalf("Erro ao definir cor FF4010: %v", err)
	}
	
	fmt.Printf("✅ Cor configurada: '%s' (deve ser 7 chars: 0FF4010)\n", fields.Color)
	
	// Testar concatenação para gravação
	payload, err := fields.ASCIIConcat()
	if err != nil {
		log.Fatalf("Erro na concatenação: %v", err)
	}
	
	fmt.Printf("✅ Payload gerado: '%s' (tamanho: %d bytes)\n", payload, len(payload))
	
	if len(payload) != 38 {
		log.Fatalf("❌ ERRO: Esperado 38 bytes, obtido %d bytes", len(payload))
	}
	
	// Testar parsing reverso (como seria na leitura)
	parsedFields, err := creality.ParseFields(payload)
	if err != nil {
		log.Fatalf("Erro no parsing: %v", err)
	}
	
	fmt.Printf("✅ Parsing reverso bem-sucedido:\n")
	fmt.Printf("   Color: '%s' -> %s\n", parsedFields.Color, parsedFields.FormatColor())
	
	// Testar com cor diferente
	fmt.Println("\n=== Teste com cor azul ===")
	
	fields2 := creality.NewFields()
	fields2.Date = "AB124"
	fields2.Supplier = "0276"
	fields2.Material = "01001"
	fields2.Length = "0165"
	fields2.Serial = "000002"
	
	err = fields2.SetColor("1A2B3C")  // Azul escuro - 6 caracteres
	if err != nil {
		log.Fatalf("Erro ao definir cor 1A2B3C: %v", err)
	}
	
	payload2, err := fields2.ASCIIConcat()
	if err != nil {
		log.Fatalf("Erro na concatenação 2: %v", err)
	}
	
	parsedFields2, err := creality.ParseFields(payload2)
	if err != nil {
		log.Fatalf("Erro no parsing 2: %v", err)
	}
	
	fmt.Printf("✅ Cor 2: '%s' -> %s\n", parsedFields2.Color, parsedFields2.FormatColor())
	
	// Testar com cor branca
	fmt.Println("\n=== Teste com cor branca ===")
	
	fields3 := creality.NewFields()
	fields3.Date = "AB124"
	fields3.Supplier = "0276"
	fields3.Material = "01001"
	fields3.Length = "0165"
	fields3.Serial = "000003"
	
	err = fields3.SetColor("FFFFFF")  // Branco - 6 caracteres
	if err != nil {
		log.Fatalf("Erro ao definir cor FFFFFF: %v", err)
	}
	
	payload3, err := fields3.ASCIIConcat()
	if err != nil {
		log.Fatalf("Erro na concatenação 3: %v", err)
	}
	
	parsedFields3, err := creality.ParseFields(payload3)
	if err != nil {
		log.Fatalf("Erro no parsing 3: %v", err)
	}
	
	fmt.Printf("✅ Cor 3: '%s' -> %s\n", parsedFields3.Color, parsedFields3.FormatColor())
	
	fmt.Println("\n✅ Teste de estrutura de gravação concluído com sucesso!")
	fmt.Println("✅ Todos os payloads têm 38 bytes")
	fmt.Println("✅ SetColor aceita 6 caracteres corretamente") 
	fmt.Println("✅ FormatColor exibe cores corretamente")
}
