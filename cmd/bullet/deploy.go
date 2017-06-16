package main

import (
	"log"

	"github.com/FurqanSoftware/bullet"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/spf13/cobra"
)

var DeployHosts string
var DeploySkipBuild bool

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

		nodes, err := bullet.ParseNodeSet(DeployHosts)
		if err != nil {
			log.Fatal(err)
			return
		}

		err = bullet.Deploy(nodes, spec, bullet.DeployOptions{
			SkipBuild: DeploySkipBuild,
		})
		if err != nil {
			log.Fatal(err)
			return
		}
	},
}

func init() {
	DeployCmd.Flags().StringVarP(&DeployHosts, "hosts", "H", "", "Hosts to deploy to")
	DeployCmd.Flags().BoolVarP(&DeploySkipBuild, "skip-build", "", false, "Skip build")
	RootCmd.AddCommand(DeployCmd)
}
