package main

import (
	"log"

	"github.com/FurqanSoftware/bullet"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/spf13/cobra"
)

var RestartHosts string

var RestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart application in server",
	Long:  `This command restarts the application in the server.`,
	Run: func(cmd *cobra.Command, args []string) {
		spec, err := spec.ParseFile("Bulletspec")
		if err != nil {
			log.Fatal(err)
			return
		}

		nodes, err := bullet.ParseNodeSet(RestartHosts)
		if err != nil {
			log.Fatal(err)
			return
		}

		err = bullet.Restart(nodes, spec)
		if err != nil {
			log.Fatal(err)
			return
		}
	},
}

func init() {
	RestartCmd.Flags().StringVarP(&RestartHosts, "hosts", "H", "", "Hosts to restart in")
	RootCmd.AddCommand(RestartCmd)
}
