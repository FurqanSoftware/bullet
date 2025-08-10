package main

import (
	"fmt"
	"log"

	"github.com/FurqanSoftware/bullet/cfg"
	"github.com/FurqanSoftware/pog"
)

func main() {
	err := cfg.LoadCurrent()
	if err != nil {
		log.Fatal(err)
	}

	log.SetPrefix("\033[2K\r")
	log.SetFlags(0)

	pog.InitDefault()

	fmt.Fprintln(log.Writer(), ` ______       _ _      _   `)
	fmt.Fprintln(log.Writer(), ` | ___ \     | | |    | |  `)
	fmt.Fprintln(log.Writer(), ` | |_/ /_   _| | | ___| |_ `)
	fmt.Fprintln(log.Writer(), ` | ___ \ | | | | |/ _ \ __|`)
	fmt.Fprintln(log.Writer(), ` | |_/ / |_| | | |  __/ |_ `)
	fmt.Fprintln(log.Writer(), ` \____/ \__,_|_|_|\___|\__|`)
	fmt.Fprintln(log.Writer())

	err = RootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
