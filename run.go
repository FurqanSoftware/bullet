package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var RunCmd = &cobra.Command{
	Use:               "run [program]",
	Short:             "Run a one-off program container",
	Long: `Start an interactive container for the given program and remove it
after it exits. Useful for running management commands or debugging.`,
	ValidArgsFunction: completeProgramKeys,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := NewSelector().Node(currentScope)
		if err != nil {
			return err
		}
		return core.Run(s, currentConfiguration, args[0])
	},
}

func init() {
	RootCmd.AddCommand(RunCmd)
}
