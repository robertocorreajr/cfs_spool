package main

import (
	"fmt"

	"github.com/robertocorreajr/cfs_spool/internal/creality"
)

func main() {
	fmt.Println("=== Teste dos Novos Nomes de Materiais ===")

	// Teste com materiais genéricos (sem prefixo "Generic")
	testMaterials := []string{"00001", "00002", "00003", "00004", "00005"}
	
	for _, code := range testMaterials {
		fields := creality.Fields{Material: code}
		fmt.Printf("Código %s -> %s\n", code, fields.GetMaterialName())
	}
	
	fmt.Println("\n=== Teste com materiais Creality ===")
	crealityMaterials := []string{"04001", "05001", "06001", "07001", "01001"}
	
	for _, code := range crealityMaterials {
		fields := creality.Fields{Material: code}
		fmt.Printf("Código %s -> %s\n", code, fields.GetMaterialName())
	}
	
	fmt.Println("\n✅ Teste de nomes de materiais concluído!")
}
