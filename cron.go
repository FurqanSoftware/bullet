package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var CronEnableCmd = &cobra.Command{
	Use:   "cron:enable",
	Short: "Enable a cron job",
	Long:  `This command enables a cron job on the server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := selectNodes(currentScope)
		return core.CronEnable(s, currentConfiguration, args)
	},
}

var CronDisableCmd = &cobra.Command{
	Use:   "cron:disable",
	Short: "Disable a cron job",
	Long:  `This command disables a cron job on the server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := selectNodes(currentScope)
		return core.CronDisable(s, currentConfiguration, args)
	},
}

var CronStatusCmd = &cobra.Command{
	Use:   "cron:status",
	Short: "Print cron job status",
	Long:  `This command prints the status of all cron jobs.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := selectNodes(currentScope)
		return core.CronStatus(s, currentConfiguration, args)
	},
}

func init() {
	RootCmd.AddCommand(CronEnableCmd)
	RootCmd.AddCommand(CronDisableCmd)
	RootCmd.AddCommand(CronStatusCmd)
}
