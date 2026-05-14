package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var RestartCmd = &cobra.Command{
	Use:   "restart [program[:instance] ...]",
	Short: "Restart application containers",
	Long: `Stop and recreate running containers for the application on the
selected nodes.

Without arguments, restarts every container of every program.

Each argument filters which containers to restart. A bare program key
restarts every instance of that program; "key:N" restarts only instance
N. Filters are additive — "web:2 worker worker:3" restarts web instance
2 plus every worker instance.`,
	ValidArgsFunction: completeProgramKeysColon,
	RunE: func(cmd *cobra.Command, args []string) error {
		var targets []core.RestartTarget
		for _, a := range args {
			parts := strings.SplitN(a, ":", 2)
			t := core.RestartTarget{Key: parts[0]}
			if len(parts) == 2 {
				n, err := strconv.Atoi(parts[1])
				if err != nil {
					return fmt.Errorf("bad argument %q: %w", a, err)
				}
				if n < 1 {
					return fmt.Errorf("bad argument %q: instance must be >= 1", a)
				}
				t.No = n
			}
			targets = append(targets, t)
		}
		return core.Restart(currentScope, currentConfiguration, targets)
	},
}

func init() {
	RootCmd.AddCommand(RestartCmd)
}
