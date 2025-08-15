package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print application status",
	Long:  `This command prints the application status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.Status(currentScope, currentConfiguration)
	},
}

func init() {
	RootCmd.AddCommand(StatusCmd)
}
