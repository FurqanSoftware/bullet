package main

import (
	"log"

	"github.com/FurqanSoftware/bullet/core"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/spf13/cobra"
)

var DfCmd = &cobra.Command{
	Use:   "df",
	Short: "Print free space on node disks",
	Long:  `This command prints the available free space on disk of the node.`,
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

		nodes = core.SelectNodes(nodes)

		err = core.Df(nodes, spec)
		if err != nil {
			log.Fatal(err)
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(DfCmd)
}
