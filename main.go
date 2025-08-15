package main

import (
	"log"

	"github.com/FurqanSoftware/pog"
)

func main() {
	log.SetPrefix("\033[2K\r")
	log.SetFlags(0)

	pog.InitDefault()

	err := RootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
