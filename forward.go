package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var ForwardCmd = &cobra.Command{
	Use:   "forward [local:]remote",
	Short: "Forward a remote port locally",
	Long: `Set up SSH port forwarding from a local port to a port on the
selected server. If only one port is given, the same port is used
locally and remotely (e.g. "8080"). Use "local:remote" to map
different ports (e.g. "3000:8080").`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := NewSelector().Node(currentScope)
		if err != nil {
			return err
		}
		return core.Forward(s, currentConfiguration, args[0])
	},
}

func init() {
	RootCmd.AddCommand(ForwardCmd)
}
