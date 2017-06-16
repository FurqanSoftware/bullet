package main

import (
	"log"
)

func main() {
	err := RootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
