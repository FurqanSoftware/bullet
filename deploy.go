package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy app to server",
	Long:  `This command packages and deploys the app to specific servers.`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"tar.gz"}, cobra.ShellCompDirectiveFilterFileExt
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		rel, err := core.NewRelease(args[0])
		if err != nil {
			return err
		}

		return core.Deploy(currentScope, currentConfiguration, rel)
	},
}

func init() {
	RootCmd.AddCommand(DeployCmd)
}
