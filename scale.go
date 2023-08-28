package main

import (
	"log"

	"github.com/FurqanSoftware/bullet/core"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/spf13/cobra"
)

var ScaleCmd = &cobra.Command{
	Use:   "scale",
	Short: "Scale a specific service on the server",
	Long:  `This command scales a specific service of the app on the server.`,
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

		comp, err := core.NewComposition(args)
		if err != nil {
			log.Fatal(err)
			return
		}

		err = core.Scale(nodes, spec, comp)
		if err != nil {
			log.Fatal(err)
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(ScaleCmd)
}
