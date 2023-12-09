package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/FurqanSoftware/bullet/core"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/spf13/cobra"
)

var LogTailCmd = &cobra.Command{
	Use:   "log",
	Short: "Tail application log",
	Long:  `This command tails the application log.`,
	Run: func(cmd *cobra.Command, args []string) {
		spec, err := spec.ParseFile("Bulletspec")
		if err != nil {
			log.Fatal(err)
			return
		}

		nodes, err := core.ParseNodeSet(Hosts, Port, Identity)
		if err != nil {
			log.Fatal(err)
			return
		}

		parts := strings.SplitN(args[0], ":", 2)
		no := 1
		if len(parts) == 2 {
			no, err = strconv.Atoi(parts[1])
			if err != nil {
				log.Fatal(err)
				return
			}
		}

		err = core.Log(nodes, spec, parts[0], no)
		if err != nil {
			log.Fatal(err)
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(LogTailCmd)
}
