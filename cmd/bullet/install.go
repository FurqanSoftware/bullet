package main

import (
	"log"

	"github.com/FurqanSoftware/bullet"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/spf13/cobra"
)

var InstallHosts string

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install application in server",
	Long:  `This command installs the application as a service in the server.`,
	Run: func(cmd *cobra.Command, args []string) {
		spec, err := spec.ParseFile("Bulletspec")
		if err != nil {
			log.Fatal(err)
			return
		}

		nodes, err := bullet.ParseNodeSet(InstallHosts)
		if err != nil {
			log.Fatal(err)
			return
		}

		err = bullet.Install(nodes, spec)
		if err != nil {
			log.Fatal(err)
			return
		}
	},
}

func init() {
	InstallCmd.Flags().StringVarP(&InstallHosts, "hosts", "H", "", "Hosts to install in")
	RootCmd.AddCommand(InstallCmd)
}
