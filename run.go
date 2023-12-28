package main

import (
	"log"

	"github.com/FurqanSoftware/bullet/core"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/spf13/cobra"
)

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a specific program on the server",
	Long:  `This command runs a specific program of the app on the server.`,
	Run: func(cmd *cobra.Command, args []string) {
		spec, err := spec.ParseFile("Bulletspec")
		if err != nil {
			log.Fatal(err)
			return
		}
		if Scope != "" {
			err = spec.ApplyScopeFile(Scope)
			if err != nil {
				log.Fatal(err)
				return
			}
		}

		nodes, err := core.ParseNodeSet(Hosts, Port, Identity)
		if err != nil {
			log.Fatal(err)
			return
		}

		err = core.Run(nodes, spec, args[0])
		if err != nil {
			log.Fatal(err)
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(RunCmd)
}
