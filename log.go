package main

import (
	"strconv"
	"strings"

	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var LogTailCmd = &cobra.Command{
	Use:   "log",
	Short: "Tail application log",
	Long:  `This command tails the application log.`,
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

		s := NewSelector().Node(currentScope)
		return core.Log(s, currentConfiguration, parts[0], no)
	},
}

func init() {
	RootCmd.AddCommand(LogTailCmd)
}
