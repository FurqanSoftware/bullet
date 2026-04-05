package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var PruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Remove old releases",
	Long: `Delete old releases from the selected nodes, keeping the 5 most
recent. This is also run automatically at the end of each deploy.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.Prune(currentScope, currentConfiguration)
	},
}

func init() {
	RootCmd.AddCommand(PruneCmd)
}
