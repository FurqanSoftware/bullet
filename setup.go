package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var (
	flagSetupEnviron string
)

var SetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup server for application",
	Long:  `This command prepares the server for the application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.Setup(currentScope, currentConfiguration, flagSetupEnviron)
	},
}

func init() {
	SetupCmd.Flags().StringVarP(&flagSetupEnviron, "environ", "", "", "if set, push file as environment")
	RootCmd.AddCommand(SetupCmd)
}
