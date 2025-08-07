package main

import (
	"fmt"
	"log"

	"github.com/robertocorreajr/cfs_spool/internal/rfid"
)

func main() {
	fmt.Println("=== Teste Básico de Leitura ===")
	
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
	fmt.Printf("UID: %s\n", uid)

	// Tentar ler bloco 0 (sempre legível)
	block0, err := rdr.ReadBlockDirect(0)
	if err != nil {
		fmt.Printf("Erro ao ler bloco 0: %v\n", err)
	} else {
		fmt.Printf("Bloco 0: %s\n", block0)
	}

	// Tentar ler outros blocos básicos
	for i := byte(1); i <= 3; i++ {
		block, err := rdr.ReadBlockDirect(i)
		if err != nil {
			fmt.Printf("Bloco %d: ERRO - %v\n", i, err)
		} else {
			fmt.Printf("Bloco %d: %s\n", i, block)
		}
	}
}
