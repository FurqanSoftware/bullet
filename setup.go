package main

import (
	"log"

	"github.com/FurqanSoftware/bullet/core"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/spf13/cobra"
)

var SetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup server for application",
	Long:  `This command prepares the server for the application.`,
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

		err = core.Setup(nodes, spec)
		if err != nil {
			log.Fatal(err)
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(SetupCmd)
}
