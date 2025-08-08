package main

import (
	"fmt"
	"log"

	"github.com/robertocorreajr/cfs_spool/internal/creality"
)

func main() {
	fmt.Println("=== Teste da Nova Estrutura de Campos ===")

	// Criar nova instância com valores fixos
	fields := creality.NewFields()
	
	// Definir outros campos
	fields.Date = "AB124"     // A B1 24
	fields.Supplier = "0276"  // Creality
	fields.Material = "01001" // Hyper PLA
	fields.Length = "0165"    // 165cm (500g)
	fields.Serial = "000001"  // Serial number
	
	// Testar SetColor
	err := fields.SetColor("FFFFF")
	if err != nil {
		log.Fatalf("Erro ao definir cor: %v", err)
	}
	
	fmt.Printf("Campos configurados:\n")
	fmt.Printf("  Batch: '%s' (fixo, tamanho: %d)\n", fields.Batch, len(fields.Batch))
	fmt.Printf("  Date: '%s' (tamanho: %d)\n", fields.Date, len(fields.Date))
	fmt.Printf("  Supplier: '%s' (tamanho: %d)\n", fields.Supplier, len(fields.Supplier))
	fmt.Printf("  Material: '%s' (tamanho: %d)\n", fields.Material, len(fields.Material))
	fmt.Printf("  Color: '%s' (tamanho: %d)\n", fields.Color, len(fields.Color))
	fmt.Printf("  Length: '%s' (tamanho: %d)\n", fields.Length, len(fields.Length))
	fmt.Printf("  Serial: '%s' (tamanho: %d)\n", fields.Serial, len(fields.Serial))
	fmt.Printf("  Reserve: '%s' (fixo, tamanho: %d)\n", fields.Reserve, len(fields.Reserve))
	
	// Testar concatenação ASCII
	payload, err := fields.ASCIIConcat()
	if err != nil {
		log.Fatalf("Erro na concatenação: %v", err)
	}
	
	fmt.Printf("\nPayload ASCII gerado: '%s'\n", payload)
	fmt.Printf("Tamanho total: %d bytes\n", len(payload))
	
	// Verificar se é exatamente 38 bytes
	if len(payload) != 38 {
		log.Fatalf("ERRO: Esperado 38 bytes, obtido %d bytes", len(payload))
	}
	
	// Testar parsing reverso
	parsedFields, err := creality.ParseFields(payload)
	if err != nil {
		log.Fatalf("Erro no parsing: %v", err)
	}
	
	fmt.Printf("\nCampos após parsing:\n")
	fmt.Printf("  Batch: '%s' (fixo)\n", parsedFields.Batch)
	fmt.Printf("  Date: '%s' -> %s\n", parsedFields.Date, parsedFields.FormatDate())
	fmt.Printf("  Supplier: '%s' -> %s\n", parsedFields.Supplier, parsedFields.GetSupplierName())
	fmt.Printf("  Material: '%s' -> %s\n", parsedFields.Material, parsedFields.GetMaterialName())
	fmt.Printf("  Color: '%s' -> %s\n", parsedFields.Color, parsedFields.FormatColor())
	fmt.Printf("  Length: '%s' -> %s\n", parsedFields.Length, parsedFields.FormatLength())
	fmt.Printf("  Serial: '%s'\n", parsedFields.Serial)
	fmt.Printf("  Reserve: '%s' (fixo)\n", parsedFields.Reserve)
	
	// Verificar se os valores fixos estão corretos
	if parsedFields.Batch != "A2" {
		log.Fatalf("ERRO: Batch deveria ser 'A2', obtido '%s'", parsedFields.Batch)
	}
	
	if parsedFields.Reserve != "000000" {
		log.Fatalf("ERRO: Reserve deveria ser '000000', obtido '%s'", parsedFields.Reserve)
	}
	
	if parsedFields.Color[0] != '0' {
		log.Fatalf("ERRO: Color deveria começar com '0', obtido '%s'", parsedFields.Color)
	}
	
	fmt.Println("\n✅ Teste concluído com sucesso!")
	fmt.Println("✅ Todos os campos estão no formato correto")
	fmt.Println("✅ Valores fixos estão sendo aplicados corretamente")
	fmt.Println("✅ Total de 38 bytes ASCII confirmado")
}
