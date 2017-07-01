package main

import (
	"log"

	"github.com/FurqanSoftware/bullet/core"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/spf13/cobra"
)

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print application status",
	Long:  `This command prints the application status.`,
	Run: func(cmd *cobra.Command, args []string) {
		spec, err := spec.ParseFile("Bulletspec")
		if err != nil {
			log.Fatal(err)
			return
		}

		nodes, err := core.ParseNodeSet(Hosts, Identity)
		if err != nil {
			log.Fatal(err)
			return
		}

		err = core.Status(nodes, spec)
		if err != nil {
			log.Fatal(err)
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(StatusCmd)
}
