package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/FurqanSoftware/bullet/core"
	"github.com/FurqanSoftware/bullet/distro"
	"github.com/spf13/cobra"
)

var (
	flagLogLines      string
	flagLogSince      string
	flagLogNoFollow   bool
	flagLogTimestamps bool
)

var LogTailCmd = &cobra.Command{
	Use:   "log [program[:instance]]",
	Short: "Tail container logs",
	Long: `Stream logs from a program's container on the selected node.
Shows the last 10 lines, then follows new output.

Specify an instance number after a colon to tail a specific
container (e.g. "web:2"). Defaults to instance 1.`,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeProgramKeysColon,
	RunE: func(cmd *cobra.Command, args []string) error {
		key, no, err := parseProgramInstance(args[0])
		if err != nil {
			return err
		}
		options, err := logOptions()
		if err != nil {
			return err
		}

		s, err := NewSelector().Node(currentScope)
		if err != nil {
			return err
		}
		return core.Log(s, currentConfiguration, key, no, options)
	},
}

// parseProgramInstance splits a "program[:instance]" argument into the
// program key and instance number, defaulting the instance to 1.
func parseProgramInstance(arg string) (string, int, error) {
	parts := strings.SplitN(arg, ":", 2)
	if len(parts) == 1 {
		return parts[0], 1, nil
	}
	no, err := strconv.Atoi(parts[1])
	if err != nil || no < 1 {
		return "", 0, fmt.Errorf("invalid instance %q, must be a positive integer", parts[1])
	}
	return parts[0], no, nil
}

// logOptions builds distro.LogOptions from the shared log flags.
func logOptions() (distro.LogOptions, error) {
	if flagLogLines != "all" {
		n, err := strconv.Atoi(flagLogLines)
		if err != nil || n < 0 {
			return distro.LogOptions{}, fmt.Errorf("invalid --lines %q, must be a non-negative integer or \"all\"", flagLogLines)
		}
	}
	return distro.LogOptions{
		Lines:      flagLogLines,
		Since:      flagLogSince,
		Follow:     !flagLogNoFollow,
		Timestamps: flagLogTimestamps,
	}, nil
}

// addLogFlags registers the shared log flags on cmd.
func addLogFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&flagLogLines, "lines", "n", "10", `number of lines to show from the end, or "all"`)
	cmd.Flags().StringVarP(&flagLogSince, "since", "S", "", `show entries newer than a duration (e.g. "2h") or timestamp`)
	cmd.Flags().BoolVarP(&flagLogNoFollow, "no-follow", "N", false, "print logs and exit instead of following")
	cmd.Flags().BoolVarP(&flagLogTimestamps, "timestamps", "t", false, "show timestamps")
}

func init() {
	addLogFlags(LogTailCmd)
	RootCmd.AddCommand(LogTailCmd)
}
