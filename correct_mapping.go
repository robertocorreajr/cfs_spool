package main

import (
	"fmt"
)

func main() {
	payload := "BB1240276A210100100000000165000001000000"
	fmt.Println("=== Mapeamento Correto dos Campos ===")
	fmt.Printf("Payload: %s (%d bytes)\n\n", payload, len(payload))
	
	// Tamanhos corretos baseados na sua tabela
	// date(5) + venderId(4) + batch(2) + filamentId(5) + color(7) + filamentLen(4) + serialNum(6) + reserve(6) = 39 bytes
	// Mas temos 38 bytes úteis, então color deve ser 6 bytes
	
	pos := 0
	date := payload[pos:pos+5]          // 5 bytes
	pos += 5
	venderId := payload[pos:pos+4]      // 4 bytes  
	pos += 4
	batch := payload[pos:pos+2]         // 2 bytes
	pos += 2
	filamentId := payload[pos:pos+5]    // 5 bytes
	pos += 5
	color := payload[pos:pos+6]         // 6 bytes (não 7)
	pos += 6
	filamentLen := payload[pos:pos+4]   // 4 bytes
	pos += 4
	serialNum := payload[pos:pos+6]     // 6 bytes
	pos += 6
	reserve := payload[pos:pos+6]       // 6 bytes
	
	fmt.Println("📋 Campos mapeados corretamente:")
	fmt.Printf("   Date: '%s' (posição 0-4)\n", date)
	fmt.Printf("   VenderId: '%s' (posição 5-8)\n", venderId)
	fmt.Printf("   Batch: '%s' (posição 9-10)\n", batch)
	fmt.Printf("   FilamentId: '%s' (posição 11-15)\n", filamentId)
	fmt.Printf("   Color: '%s' (posição 16-21)\n", color)
	fmt.Printf("   FilamentLen: '%s' (posição 22-25)\n", filamentLen)
	fmt.Printf("   SerialNum: '%s' (posição 26-31)\n", serialNum)
	fmt.Printf("   Reserve: '%s' (posição 32-37)\n", reserve)
	
	fmt.Println("\n✅ Verificações:")
	fmt.Printf("   VenderId é '0276' (Creality)? %t\n", venderId == "0276")
	fmt.Printf("   Batch é 'A2'? %t\n", batch == "A2")
	fmt.Printf("   SerialNum termina com '001'? %t\n", serialNum[3:] == "001")
	fmt.Printf("   Reserve termina com '000'? %t\n", reserve[3:] == "000")
}
