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
	Short: "Prepare servers for deployment",
	Long: `Install Docker and create the application directory structure on the
selected servers. Optionally push an initial environment file with --environ.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.Setup(currentScope, currentConfiguration, flagSetupEnviron)
	},
}

func init() {
	SetupCmd.Flags().StringVarP(&flagSetupEnviron, "environ", "", "", "if set, push file as environment")
	RootCmd.AddCommand(SetupCmd)
}
