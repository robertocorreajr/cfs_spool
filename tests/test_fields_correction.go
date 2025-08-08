package main

import (
	"fmt"
	"log"

	"github.com/robertocorreajr/cfs_spool/internal/creality"
)

func main() {
	fmt.Println("=== Teste das correções dos campos ===")

	// Criar nova instância com valores fixos
	fields := creality.NewFields()
	
	// Definir outros campos
	fields.Date = "AB124"     // A B1 24
	fields.Supplier = "0276"  // Creality
	fields.Material = "01001" // Hyper PLA
	fields.Length = "0165"    // 165cm (500g)
	fields.Serial = "000001"  // Serial number
	
	// Testar SetColor
	err := fields.SetColor("FFFFFF")  // 6 caracteres para cor branca
	if err != nil {
		log.Fatalf("Erro ao definir cor: %v", err)
	}
	
	fmt.Printf("Campos configurados:\n")
	fmt.Printf("  Batch: %s (fixo, tamanho: %d)\n", fields.Batch, len(fields.Batch))
	fmt.Printf("  Date: %s (tamanho: %d)\n", fields.Date, len(fields.Date))
	fmt.Printf("  Supplier: %s (tamanho: %d)\n", fields.Supplier, len(fields.Supplier))
	fmt.Printf("  Material: %s (tamanho: %d)\n", fields.Material, len(fields.Material))
	fmt.Printf("  Color: %s (tamanho: %d)\n", fields.Color, len(fields.Color))
	fmt.Printf("  Length: %s (tamanho: %d)\n", fields.Length, len(fields.Length))
	fmt.Printf("  Serial: %s (tamanho: %d)\n", fields.Serial, len(fields.Serial))
	fmt.Printf("  Reserve: %s (fixo, tamanho: %d)\n", fields.Reserve, len(fields.Reserve))
	
	// Testar concatenação ASCII
	payload, err := fields.ASCIIConcat()
	if err != nil {
		log.Fatalf("Erro na concatenação: %v", err)
	}
	
	fmt.Printf("\nPayload ASCII gerado: %s\n", payload)
	fmt.Printf("Tamanho total: %d bytes\n", len(payload))
	
	// Testar parsing reverso
	parsedFields, err := creality.ParseFields(payload)
	if err != nil {
		log.Fatalf("Erro no parsing: %v", err)
	}
	
	fmt.Printf("\nCampos após parsing:\n")
	fmt.Printf("  Batch: %s\n", parsedFields.Batch)
	fmt.Printf("  Date: %s\n", parsedFields.Date)
	fmt.Printf("  Supplier: %s (%s)\n", parsedFields.Supplier, parsedFields.GetSupplierName())
	fmt.Printf("  Material: %s (%s)\n", parsedFields.Material, parsedFields.GetMaterialName())
	fmt.Printf("  Color: %s (%s)\n", parsedFields.Color, parsedFields.FormatColor())
	fmt.Printf("  Length: %s (%s)\n", parsedFields.Length, parsedFields.FormatLength())
	fmt.Printf("  Serial: %s\n", parsedFields.Serial)
	fmt.Printf("  Reserve: %s\n", parsedFields.Reserve)
	
	// Testar com cor diferente
	fmt.Printf("\n=== Teste com cor diferente ===\n")
	fields2 := creality.NewFields()
	fields2.Date = "9A224"     // 9 A2 24
	fields2.Supplier = "0276"  // Creality
	fields2.Material = "01001" // Hyper PLA
	fields2.Length = "0165"    // 165cm (500g)
	fields2.Serial = "000001"  // Serial number
	
	err = fields2.SetColor("1C2E3F")  // Cor azul escura (6 caracteres)
	if err != nil {
		log.Fatalf("Erro ao definir cor 2: %v", err)
	}
	
	payload2, err := fields2.ASCIIConcat()
	if err != nil {
		log.Fatalf("Erro na concatenação 2: %v", err)
	}
	
	fmt.Printf("Payload 2: %s (tamanho: %d)\n", payload2, len(payload2))
	
	parsedFields2, err := creality.ParseFields(payload2)
	if err != nil {
		log.Fatalf("Erro no parsing 2: %v", err)
	}
	
	fmt.Printf("Color 2: %s (%s)\n", parsedFields2.Color, parsedFields2.FormatColor())
	
	fmt.Println("\n=== Teste concluído com sucesso! ===")
}
