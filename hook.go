package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var HookRunCmd = &cobra.Command{
	Use:   "hook:run [hook ...]",
	Short: "Run post-deploy hooks",
	Long: `Run the specified post-deploy hooks without deploying. Each hook runs
once on the first selected node, or on every selected node if its nodes
setting is "all".`,
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: completeHookKeys,
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.HookRun(currentScope, currentConfiguration, args)
	},
}

func init() {
	RootCmd.AddCommand(HookRunCmd)
}
