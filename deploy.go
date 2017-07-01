package main

import (
	"log"

	"github.com/FurqanSoftware/bullet/core"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/spf13/cobra"
)

var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy app to server",
	Long:  `This command packages and deploys the app to specific servers.`,
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

		rel, err := core.NewRelease(args[0])
		if err != nil {
			log.Fatal(err)
			return
		}

		err = core.Deploy(nodes, spec, rel)
		if err != nil {
			log.Fatal(err)
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(DeployCmd)
}
