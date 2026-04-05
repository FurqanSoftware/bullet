package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var RestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart all application containers",
	Long: `Stop and recreate all running containers for the application on
the selected nodes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := NewSelector().Nodes(currentScope)
		if err != nil {
			return err
		}
		return core.Restart(s, currentConfiguration)
	},
}

func init() {
	RootCmd.AddCommand(RestartCmd)
}
