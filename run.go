package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a specific program on the server",
	Long:  `This command runs a specific program of the app on the server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := NewSelector().Node(currentScope)
		return core.Run(s, currentConfiguration, args[0])
	},
}

func init() {
	RootCmd.AddCommand(RunCmd)
}
