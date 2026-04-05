package main

import (
	"strconv"
	"strings"

	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var LogTailCmd = &cobra.Command{
	Use:               "log [program[:instance]]",
	Short:             "Tail container logs",
	Long: `Stream logs from a program's container on the selected node.
Shows the last 10 lines, then follows new output.

Specify an instance number after a colon to tail a specific
container (e.g. "web:2"). Defaults to instance 1.`,
	ValidArgsFunction: completeProgramKeysColon,
	RunE: func(cmd *cobra.Command, args []string) error {
		parts := strings.SplitN(args[0], ":", 2)
		no := 1
		if len(parts) == 2 {
			var err error
			no, err = strconv.Atoi(parts[1])
			if err != nil {
				return err
			}
		}

		s, err := NewSelector().Node(currentScope)
		if err != nil {
			return err
		}
		return core.Log(s, currentConfiguration, parts[0], no)
	},
}

func init() {
	RootCmd.AddCommand(LogTailCmd)
}
