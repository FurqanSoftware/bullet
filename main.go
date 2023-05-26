package main

import (
	"log"
)

func main() {
	log.SetFlags(0)

	log.Println(` ______       _ _      _   `)
	log.Println(` | ___ \     | | |    | |  `)
	log.Println(` | |_/ /_   _| | | ___| |_ `)
	log.Println(` | ___ \ | | | | |/ _ \ __|`)
	log.Println(` | |_/ / |_| | | |  __/ |_ `)
	log.Println(` \____/ \__,_|_|_|\___|\__|`)
	log.Println()

	err := RootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
