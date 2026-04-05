package main

import (
	"path/filepath"
	"strings"

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

// RootCmd is the top-level Cobra command for the bullet CLI.
var RootCmd = &cobra.Command{
	Use:          "bullet",
	Short:        "Bullet is a fast application deployment tool",
	Long:         `Bullet is a fast and flexible application deployment tool built by Furqan Software and friends. Complete documentation is available at https://bullettool.com/.`,
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		for c := cmd; c != nil; c = c.Parent() {
			if c.Name() == "completion" || c.Name() == "help" {
				return nil
			}
		}

		printBanner()

		g, err := cfg.NewLoader().
			ParseFileIfExists("Bulletcfg." + flagConfig).
			ApplyEnvironment().
			ApplyFlags(cmd.Flags()).
			Configuration()
		if err != nil {
			return err
		}

		s := scope.Scope{}
		s.Spec, err = spec.ParseFile("Bulletspec")
		if err != nil {
			return err
		}
		s.Spec.ExpandVars(g.Vars)

		s.Nodes, err = scope.ParseNodeSet(g.Hosts, g.Port, g.Identity)
		if err != nil {
			return err
		}

		currentScope = s
		currentConfiguration = g
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&flagConfig, "config", "c", "", "Name of the configuration to apply")
	RootCmd.RegisterFlagCompletionFunc("config", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		matches, _ := filepath.Glob("Bulletcfg.*")
		names := make([]string, 0, len(matches))
		for _, m := range matches {
			names = append(names, strings.TrimPrefix(m, "Bulletcfg."))
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	})
	cfg.AddFlags(RootCmd.PersistentFlags())
}
