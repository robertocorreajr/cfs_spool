package main

import (
	"fmt"
)

func main() {
	// Payload da tag original
	payload := "BB1240276A210100100000000165000001000000"
	
	fmt.Println("=== Análise do Payload da Tag Original ===")
	fmt.Printf("Payload: %s (%d bytes)\n\n", payload, len(payload))
	
	// Procurar por padrões conhecidos
	fmt.Println("🔍 Procurando padrões conhecidos:")
	fmt.Printf("- Creality (0276): posição %d\n", findSubstring(payload, "0276"))
	fmt.Printf("- A2: posição %d\n", findSubstring(payload, "A2"))
	fmt.Printf("- 000001 (serial): posição %d\n", findSubstring(payload, "000001"))
	fmt.Printf("- 000000 (reserve): posição %d\n", findSubstring(payload, "000000"))
	
	fmt.Println("\n📋 Tentativas de mapeamento:")
	
	// Tentativa 1: Formato original que estava sendo usado
	fmt.Println("\n1. Mapeamento atual (pode estar errado):")
	if len(payload) >= 38 {
		fmt.Printf("   Batch: '%s' (pos 0-1)\n", payload[0:2])
		fmt.Printf("   Date: '%s' (pos 2-6)\n", payload[2:7])
		fmt.Printf("   Supplier: '%s' (pos 7-10)\n", payload[7:11])
		fmt.Printf("   Material: '%s' (pos 11-15)\n", payload[11:16])
		fmt.Printf("   Color: '%s' (pos 16-21)\n", payload[16:22])
		fmt.Printf("   Length: '%s' (pos 22-25)\n", payload[22:26])
		fmt.Printf("   Serial: '%s' (pos 26-31)\n", payload[26:32])
		fmt.Printf("   Reserve: '%s' (pos 32-37)\n", payload[32:38])
	}
	
	// Tentativa 2: Procurar onde está realmente o 0276 e A2
	fmt.Println("\n2. Analisando posições reais:")
	for i := 0; i < len(payload)-3; i++ {
		substr := payload[i:i+4]
		if substr == "0276" {
			fmt.Printf("   0276 (Creality) encontrado na posição %d\n", i)
		}
		if i < len(payload)-1 && payload[i:i+2] == "A2" {
			fmt.Printf("   A2 encontrado na posição %d\n", i)
		}
	}
	
	// Tentativa 3: Mapeamento baseado nas posições encontradas
	fmt.Println("\n3. Proposta de novo mapeamento:")
	fmt.Println("   Se 0276 está na posição correta para Supplier...")
	
	// Vou tentar várias combinações
	testMappings := []struct{
		name string
		batch, date, supplier, material, color, length, serial, reserve string
	}{
		{
			"Tentativa A",
			payload[0:2], payload[2:7], payload[7:11], payload[11:16], 
			payload[16:22], payload[22:26], payload[26:32], payload[32:38],
		},
		{
			"Tentativa B (shift +1)",
			payload[1:3], payload[3:8], payload[8:12], payload[12:17],
			payload[17:23], payload[23:27], payload[27:33], payload[33:38],
		},
	}
	
	for _, tm := range testMappings {
		fmt.Printf("\n   %s:\n", tm.name)
		fmt.Printf("     Batch: '%s'\n", tm.batch)
		fmt.Printf("     Date: '%s'\n", tm.date)
		fmt.Printf("     Supplier: '%s'\n", tm.supplier)
		fmt.Printf("     Material: '%s'\n", tm.material)
		fmt.Printf("     Color: '%s'\n", tm.color)
		fmt.Printf("     Length: '%s'\n", tm.length)
		fmt.Printf("     Serial: '%s'\n", tm.serial)
		fmt.Printf("     Reserve: '%s'\n", tm.reserve)
	}
}

func findSubstring(str, substr string) int {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
