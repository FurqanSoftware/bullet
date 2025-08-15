package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/FurqanSoftware/bullet/distro"
	"github.com/spf13/cobra"
)

var HostShellCmd = &cobra.Command{
	Use:   "host:shell",
	Short: "Connects to node over SSH",
	Long:  `This command starts an SSH session with the selected node.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := NewSelector().Node(currentScope)
		return core.Shell(s, currentConfiguration)
	},
}

var (
	flagDfWatch     bool
	flagDfArguments string
)

var HostDfCmd = &cobra.Command{
	Use:   "host:df",
	Short: "Print free space on node disks",
	Long:  `This command prints the available free space on disk of the node.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := NewSelector().Node(currentScope)
		return core.Df(s, currentConfiguration, distro.DfOptions{
			Watch:     flagDfWatch,
			Arguments: flagDfArguments,
		})
	},
}

var HostTopCmd = &cobra.Command{
	Use:   "host:top",
	Short: "Display Linux processes",
	Long:  `This command displays Linux processes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := NewSelector().Node(currentScope)
		return core.Top(s, currentConfiguration)
	},
}

func init() {
	RootCmd.AddCommand(HostShellCmd)

	HostDfCmd.Flags().BoolVarP(&flagDfWatch, "watch", "w", false, "runs df with watch")
	HostDfCmd.Flags().StringVarP(&flagDfArguments, "arguments", "a", "", "additional arguments for df")
	RootCmd.AddCommand(HostDfCmd)

	RootCmd.AddCommand(HostTopCmd)
}
