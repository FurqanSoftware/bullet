package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var (
	flagDeployEnviron string
	flagDeployScale   bool
	flagDeploySetup   bool
)

var DeployCmd = &cobra.Command{
	Use:   "deploy [tarball]",
	Short: "Deploy a release to servers",
	Long: `Upload a tarball to the selected servers, extract it as a new release,
build Docker images, and reload running containers.

Skips nodes where the same release (by SHA256 hash) is already deployed.
Old releases are pruned automatically, keeping the 5 most recent.
Optionally push an environment file before deploying with --environ.
Use --setup to run server setup (install Docker, create directories) before deploying.
Use --scale to automatically scale programs using the Bulletspec rules after deploying.`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"tar.gz"}, cobra.ShellCompDirectiveFilterFileExt
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		rel, err := core.NewRelease(args[0])
		if err != nil {
			return err
		}

		return core.Deploy(currentScope, currentConfiguration, rel, flagDeployEnviron, flagDeploySetup, flagDeployScale)
	},
}

func init() {
	DeployCmd.Flags().StringVarP(&flagDeployEnviron, "environ", "", "", "if set, push file as environment")
	DeployCmd.Flags().BoolVarP(&flagDeploySetup, "setup", "", false, "if set, run server setup before deploying")
	DeployCmd.Flags().BoolVarP(&flagDeployScale, "scale", "", false, "if set, auto scale programs after deploying")
	RootCmd.AddCommand(DeployCmd)
}
