package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var ForwardCmd = &cobra.Command{
	Use:   "forward",
	Short: "Forward a specific port from the server",
	Long:  `This command forwards a specific program from the server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := NewSelector().Node(currentScope)
		return core.Forward(s, currentConfiguration, args[0])
	},
}

func init() {
	RootCmd.AddCommand(ForwardCmd)
}
