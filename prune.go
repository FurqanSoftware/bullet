package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var PruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Remove old release tarballs",
	Long:  `This command removes old release tarballs.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.Prune(currentScope, currentConfiguration)
	},
}

func init() {
	RootCmd.AddCommand(PruneCmd)
}
