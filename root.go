package main

import (
	"log"

	"github.com/FurqanSoftware/bullet/cfg"
	"github.com/FurqanSoftware/bullet/scope"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/spf13/cobra"
)

var (
	flagConfig string

	currentScope         scope.Scope
	currentConfiguration cfg.Configuration
)

var RootCmd = &cobra.Command{
	Use:   "bullet",
	Short: "Bullet is a fast application deployment tool",
	Long:  `Bullet is a fast and flexible application deployment tool built by Furqan Software and friends. Complete documentation is available at https://bullettool.com/.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		printBanner()

		g, err := cfg.NewLoader().
			ParseFileIfExists("Bulletcfg." + flagConfig).
			ApplyEnvironment().
			ApplyFlags(cmd.Flags()).
			Configuration()
		if err != nil {
			log.Fatal(err)
		}

		s := scope.Scope{}
		s.Spec, err = spec.ParseFile("Bulletspec")
		if err != nil {
			log.Fatal(err)
			return
		}
		s.Spec.ExpandVars(g.Vars)

		s.Nodes, err = scope.ParseNodeSet(g.Hosts, g.Port, g.Identity)
		if err != nil {
			log.Fatal(err)
			return
		}

		currentScope = s
		currentConfiguration = g
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&flagConfig, "config", "c", "", "Name of the configuration to apply")
	cfg.AddFlags(RootCmd.PersistentFlags())
}
