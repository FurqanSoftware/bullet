package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/FurqanSoftware/bullet/distro"
	"github.com/spf13/cobra"
)

var HostShellCmd = &cobra.Command{
	Use:   "host:shell",
	Short: "Open an interactive shell on a server",
	Long:  `Start an interactive bash session over SSH on the selected node.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := NewSelector().Node(currentScope)
		if err != nil {
			return err
		}
		return core.Shell(s, currentConfiguration)
	},
}

var (
	flagDfWatch     bool
	flagDfArguments string
)

var HostDfCmd = &cobra.Command{
	Use:   "host:df",
	Short: "Show disk usage on a server",
	Long: `Run df on the selected node to show disk space usage. Use --watch
for continuous updates or --arguments to pass additional flags to df.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := NewSelector().Node(currentScope)
		if err != nil {
			return err
		}
		return core.Df(s, currentConfiguration, distro.DfOptions{
			Watch:     flagDfWatch,
			Arguments: flagDfArguments,
		})
	},
}

var HostTopCmd = &cobra.Command{
	Use:   "host:top",
	Short: "Show running processes on a server",
	Long:  `Run top interactively on the selected node.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := NewSelector().Node(currentScope)
		if err != nil {
			return err
		}
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
