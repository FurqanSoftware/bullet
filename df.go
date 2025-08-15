package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/FurqanSoftware/bullet/distro"
	"github.com/spf13/cobra"
)

var (
	flagDfWatch     bool
	flagDfArguments string
)

var DfCmd = &cobra.Command{
	Use:   "df",
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

func init() {
	DfCmd.Flags().BoolVarP(&flagDfWatch, "watch", "w", false, "runs df with watch")
	DfCmd.Flags().StringVarP(&flagDfArguments, "arguments", "a", "", "additional arguments for df")
	RootCmd.AddCommand(DfCmd)
}
