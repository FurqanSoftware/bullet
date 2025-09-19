package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var EnvironPushCmd = &cobra.Command{
	Use:   "environ:push",
	Short: "Push application environment from file to server",
	Long:  `This command configures the application in the server based on specific environment file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := NewSelector().Nodes(currentScope)
		return core.EnvironPush(s, currentConfiguration, args[0])
	},
}

func init() {
	RootCmd.AddCommand(EnvironPushCmd)
}
