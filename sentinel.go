package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var SentinelInstallCmd = &cobra.Command{
	Use:   "sentinel:install",
	Short: "Install the sentinel on the selected nodes",
	Long: `Push the embedded bullet-sentinel binary and systemd unit, and enable
bullet-sentinel.service on each selected node. The sentinel config file at
/etc/bullet/sentinel.yaml is expected to be managed out-of-band.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.SentinelInstall(currentScope, currentConfiguration, sentinelBinary)
	},
}

func init() {
	RootCmd.AddCommand(SentinelInstallCmd)
}
