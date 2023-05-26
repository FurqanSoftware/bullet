package main

import (
	"log"

	"github.com/FurqanSoftware/bullet/core"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/spf13/cobra"
)

var CronEnableCmd = &cobra.Command{
	Use:   "cron:enable",
	Short: "Enable a cron job",
	Long:  `This command enables a cron job on the server.`,
	Run: func(cmd *cobra.Command, args []string) {
		spec, err := spec.ParseFile("Bulletspec")
		if err != nil {
			log.Fatal(err)
			return
		}

		nodes, err := core.ParseNodeSet(Hosts, Port, Identity)
		if err != nil {
			log.Fatal(err)
			return
		}

		err = core.CronEnable(nodes, spec, args)
		if err != nil {
			log.Fatal(err)
			return
		}
	},
}

var CronDisableCmd = &cobra.Command{
	Use:   "cron:disable",
	Short: "Disable a cron job",
	Long:  `This command disables a cron job on the server.`,
	Run: func(cmd *cobra.Command, args []string) {
		spec, err := spec.ParseFile("Bulletspec")
		if err != nil {
			log.Fatal(err)
			return
		}

		nodes, err := core.ParseNodeSet(Hosts, Port, Identity)
		if err != nil {
			log.Fatal(err)
			return
		}

		err = core.CronDisable(nodes, spec, args)
		if err != nil {
			log.Fatal(err)
			return
		}
	},
}

var CronStatusCmd = &cobra.Command{
	Use:   "cron:status",
	Short: "Print cron job status",
	Long:  `This command prints the status of all cron jobs.`,
	Run: func(cmd *cobra.Command, args []string) {
		spec, err := spec.ParseFile("Bulletspec")
		if err != nil {
			log.Fatal(err)
			return
		}

		nodes, err := core.ParseNodeSet(Hosts, Port, Identity)
		if err != nil {
			log.Fatal(err)
			return
		}

		err = core.CronStatus(nodes, spec, args)
		if err != nil {
			log.Fatal(err)
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(CronEnableCmd)
	RootCmd.AddCommand(CronDisableCmd)
	RootCmd.AddCommand(CronStatusCmd)
}
