package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var (
	flagEnvironPushNoRestart bool
)

var EnvironPushCmd = &cobra.Command{
	Use:   "environ:push [file]",
	Short: "Push an environment file to servers",
	Long: `Upload a local environment file to the selected nodes. The file is
stored at /opt/<identifier>/env and is loaded by all containers and
cron jobs via Docker's --env-file flag.

If the file has changed, running containers are recreated so that they pick
up the new environment. Use --no-restart to upload without restarting.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.EnvironPush(currentScope, currentConfiguration, args[0], !flagEnvironPushNoRestart)
	},
}

func init() {
	EnvironPushCmd.Flags().BoolVarP(&flagEnvironPushNoRestart, "no-restart", "", false, "if set, do not restart containers")
	RootCmd.AddCommand(EnvironPushCmd)
}
