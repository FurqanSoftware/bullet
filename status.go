package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show container status across nodes",
	Long: `Print a table showing the number of running and healthy containers
for each program on every node.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.Status(currentScope, currentConfiguration)
	},
}

func init() {
	RootCmd.AddCommand(StatusCmd)
}
