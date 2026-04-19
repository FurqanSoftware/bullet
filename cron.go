package main

import (
	"github.com/FurqanSoftware/bullet/core"
	"github.com/spf13/cobra"
)

var CronEnableCmd = &cobra.Command{
	Use:               "cron:enable [job ...]",
	Short:             "Enable cron jobs",
	Long: `Create systemd timer and service units for the specified cron jobs
and enable them on the selected nodes.`,
	ValidArgsFunction: completeCronJobKeys,
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.CronEnable(currentScope, currentConfiguration, args)
	},
}

var CronDisableCmd = &cobra.Command{
	Use:               "cron:disable [job ...]",
	Short:             "Disable cron jobs",
	Long: `Stop and remove systemd timer and service units for the specified
cron jobs on the selected nodes.`,
	ValidArgsFunction: completeCronJobKeys,
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.CronDisable(currentScope, currentConfiguration, args)
	},
}

var CronStatusCmd = &cobra.Command{
	Use:   "cron:status",
	Short: "Show cron job status",
	Long: `Print the status of all cron jobs on the selected nodes, including
whether each timer is active and when it will next trigger.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.CronStatus(currentScope, currentConfiguration, args)
	},
}

func init() {
	RootCmd.AddCommand(CronEnableCmd)
	RootCmd.AddCommand(CronDisableCmd)
	RootCmd.AddCommand(CronStatusCmd)
}
