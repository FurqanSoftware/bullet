package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var EnvironPushCmd = &cobra.Command{
	Use:   "environ:push [file]",
	Short: "Push an environment file to servers",
	Long: `Upload a local environment file to the selected nodes. The file is
stored at /opt/<identifier>/env and is loaded by all containers and
cron jobs via Docker's --env-file flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := NewSelector().Nodes(currentScope)
		if err != nil {
			return err
		}
		return core.EnvironPush(s, currentConfiguration, args[0])
	},
}

func init() {
	RootCmd.AddCommand(EnvironPushCmd)
}
