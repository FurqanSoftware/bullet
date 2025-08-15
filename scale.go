package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var ScaleCmd = &cobra.Command{
	Use:   "scale",
	Short: "Scale a specific service on the server",
	Long:  `This command scales a specific service of the app on the server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		comp, err := core.NewComposition(args)
		if err != nil {
			return err
		}

		s := NewSelector().Nodes(currentScope)
		return core.Scale(s, currentConfiguration, comp)
	},
}

func init() {
	RootCmd.AddCommand(ScaleCmd)
}
