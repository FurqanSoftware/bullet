package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var ScaleCmd = &cobra.Command{
	Use:               "scale [program=count ...]",
	Short:             "Scale program instances",
	Long: `Adjust the number of container instances for one or more programs.

When called with arguments (e.g. "web=4 worker=2"), scales to the
specified counts. When called without arguments, evaluates the
scaling rules defined in the Bulletspec.`,
	ValidArgsFunction: completeProgramKeysEquals,
	RunE: func(cmd *cobra.Command, args []string) error {
		comp, err := core.NewComposition(args)
		if err != nil {
			return err
		}

		s, err := NewSelector().Nodes(currentScope)
		if err != nil {
			return err
		}
		return core.Scale(s, currentConfiguration, comp)
	},
}

func init() {
	RootCmd.AddCommand(ScaleCmd)
}
