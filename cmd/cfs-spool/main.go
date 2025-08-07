package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	flag.Parse()
	
	if len(flag.Args()) == 0 {
		usage()
		return
	}

	switch flag.Arg(0) {
	case "read-tag":
		cmdReadTag(flag.Args()[1:])
	case "write-tag":
		cmdWriteTag(flag.Args()[1:])
	default:
		usage()
	}
}

func usage() {
	fmt.Println("cfs-spool <command> [flags]")
	fmt.Println("Commands:")
	fmt.Println("  read-tag               – lê UID + conteúdo e decodifica")
	fmt.Println("  write-tag [flags]      – grava nova tag")
}

func dieIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
