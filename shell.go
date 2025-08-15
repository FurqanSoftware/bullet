package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var ShellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Connects to node over SSH",
	Long:  `This command starts an SSH session with the selected node.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := NewSelector().Node(currentScope)
		return core.Shell(s, currentConfiguration)
	},
}

func init() {
	RootCmd.AddCommand(ShellCmd)
}
