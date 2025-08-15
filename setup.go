package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var SetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup server for application",
	Long:  `This command prepares the server for the application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		environ, _ := cmd.Flags().GetString("environ")
		return core.Setup(currentScope, currentConfiguration, environ)
	},
}

func init() {
	SetupCmd.Flags().String("environ", "", "if set, push file as environment")
	RootCmd.AddCommand(SetupCmd)
}
