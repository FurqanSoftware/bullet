package main

import (
	"log"

	"github.com/FurqanSoftware/bullet"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/spf13/cobra"
)

var ConfigHosts string

var ConfigPushCmd = &cobra.Command{
	Use:   "config:push",
	Short: "Push application configuration from file to server",
	Long:  `This command configures the application in the server based on specific environment file.`,
	Run: func(cmd *cobra.Command, args []string) {
		spec, err := spec.ParseFile("Bulletspec")
		if err != nil {
			log.Fatal(err)
			return
		}

		nodes, err := bullet.ParseNodeSet(ConfigHosts)
		if err != nil {
			log.Fatal(err)
			return
		}

		err = bullet.ConfigPush(nodes, spec, args[0])
		if err != nil {
			log.Fatal(err)
			return
		}
	},
}

func init() {
	ConfigPushCmd.Flags().StringVarP(&ConfigHosts, "hosts", "H", "", "Hosts to configure")
	RootCmd.AddCommand(ConfigPushCmd)
}
