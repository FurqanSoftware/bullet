package main

import (
	"log"

	"github.com/FurqanSoftware/bullet"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/spf13/cobra"
)

var SetupHosts string

var SetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup server for app",
	Long:  `This command prepares the server for the app.`,
	Run: func(cmd *cobra.Command, args []string) {
		spec, err := spec.ParseFile("Bulletspec")
		if err != nil {
			log.Fatal(err)
			return
		}

		nodes, err := bullet.ParseNodeSet(SetupHosts)
		if err != nil {
			log.Fatal(err)
			return
		}

		err = bullet.Setup(nodes, spec)
		if err != nil {
			log.Fatal(err)
			return
		}
	},
}

func init() {
	SetupCmd.Flags().StringVarP(&SetupHosts, "hosts", "H", "", "Hosts to configure")
	RootCmd.AddCommand(SetupCmd)
}
