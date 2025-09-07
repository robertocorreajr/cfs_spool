package main

import (
	"fmt"
	"log"

	"github.com/robertocorreajr/cfs_spool/internal/creality"
)

func main() {
	fmt.Println("=== Teste de Compatibilidade de Formatos ===")

	// Teste 1: Dados de 38 bytes (formato novo)
	fmt.Println("\n1. Teste com 38 bytes (formato novo):")
	data38 := "A2AB1240276010010FFFFF0165000001000000"
	fmt.Printf("   Dados: '%s' (%d bytes)\n", data38, len(data38))
	
	fields38, err := creality.ParseFieldsCompat(data38)
	if err != nil {
		log.Printf("   ❌ Erro: %v", err)
	} else {
		fmt.Printf("   ✅ Parseado: Batch='%s', Color='%s', Reserve='%s'\n", 
			fields38.Batch, fields38.Color, fields38.Reserve)
	}

	// Teste 2: Dados de 48 bytes (formato antigo)
	fmt.Println("\n2. Teste com 48 bytes (formato antigo):")
	data48 := "A2AB1240276010010FFFFF0165000001000000" + "0000000000"
	fmt.Printf("   Dados: '%s' (%d bytes)\n", data48, len(data48))
	
	fields48, err := creality.ParseFieldsCompat(data48)
	if err != nil {
		log.Printf("   ❌ Erro: %v", err)
	} else {
		fmt.Printf("   ✅ Parseado: Batch='%s', Color='%s', Reserve='%s'\n", 
			fields48.Batch, fields48.Color, fields48.Reserve)
	}

	// Teste 3: Dados menores que 38 bytes (padding)
	fmt.Println("\n3. Teste com dados menores (será preenchido):")
	dataShort := "A2AB124027601001"
	fmt.Printf("   Dados: '%s' (%d bytes)\n", dataShort, len(dataShort))
	
	fieldsShort, err := creality.ParseFieldsCompat(dataShort)
	if err != nil {
		log.Printf("   ❌ Erro: %v", err)
	} else {
		fmt.Printf("   ✅ Parseado: Batch='%s', Date='%s', Supplier='%s'\n", 
			fieldsShort.Batch, fieldsShort.Date, fieldsShort.Supplier)
	}

	// Teste 4: Geração de payload de 48 bytes
	fmt.Println("\n4. Teste de geração de payload 48 bytes:")
	fields := creality.NewFields()
	fields.Date = "AB124"
	fields.Supplier = "0276"
	fields.Material = "01001"
	fields.SetColor("FFFFF")
	fields.Length = "0165"
	fields.Serial = "000001"
	
	payload48, err := fields.ASCIIConcat48()
	if err != nil {
		log.Printf("   ❌ Erro: %v", err)
	} else {
		fmt.Printf("   ✅ Payload 48: '%s' (%d bytes)\n", payload48, len(payload48))
		fmt.Printf("   Dados úteis: '%s' (38 bytes)\n", payload48[:38])
		fmt.Printf("   Padding: '%s' (10 bytes)\n", payload48[38:])
	}

	fmt.Println("\n✅ Testes de compatibilidade concluídos!")
}
