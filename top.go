package main

import (
	"log"

	"github.com/FurqanSoftware/bullet/core"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/spf13/cobra"
)

var TopCmd = &cobra.Command{
	Use:   "top",
	Short: "Display Linux processes",
	Long:  `This command displays Linux processes.`,
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

		node := core.SelectNode(nodes)

		err = core.Top(node, spec)
		if err != nil {
			log.Fatal(err)
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(TopCmd)
}
