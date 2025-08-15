package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var TopCmd = &cobra.Command{
	Use:   "top",
	Short: "Display Linux processes",
	Long:  `This command displays Linux processes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := NewSelector().Node(currentScope)
		return core.Top(s, currentConfiguration)
	},
}

func init() {
	RootCmd.AddCommand(TopCmd)
}
