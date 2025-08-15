package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var RestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart application in server",
	Long:  `This command restarts the application in the server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := NewSelector().Nodes(currentScope)
		return core.Restart(s, currentConfiguration)
	},
}

func init() {
	RootCmd.AddCommand(RestartCmd)
}
