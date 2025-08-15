package main

import (
	"fmt"
	"log"
)

func printBanner() {
	fmt.Fprintln(log.Writer(), ` ______       _ _      _   `)
	fmt.Fprintln(log.Writer(), ` | ___ \     | | |    | |  `)
	fmt.Fprintln(log.Writer(), ` | |_/ /_   _| | | ___| |_ `)
	fmt.Fprintln(log.Writer(), ` | ___ \ | | | | |/ _ \ __|`)
	fmt.Fprintln(log.Writer(), ` | |_/ / |_| | | |  __/ |_ `)
	fmt.Fprintln(log.Writer(), ` \____/ \__,_|_|_|\___|\__|`)
	fmt.Fprintln(log.Writer())

}
