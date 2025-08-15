package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var DfCmd = &cobra.Command{
	Use:   "df",
	Short: "Print free space on node disks",
	Long:  `This command prints the available free space on disk of the node.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := NewSelector().Nodes(currentScope)
		return core.Df(s, currentConfiguration)
	},
}

func init() {
	RootCmd.AddCommand(DfCmd)
}
