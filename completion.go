package main

import (
	"strings"

	"github.com/FurqanSoftware/bullet/spec"
	"github.com/spf13/cobra"
)

func completeProgramKeys(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	s, err := spec.ParseFile("Bulletspec")
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	keys := make([]string, 0, len(s.Application.Programs))
	for k := range s.Application.Programs {
		keys = append(keys, k)
	}
	return keys, cobra.ShellCompDirectiveNoFileComp
}

func completeProgramKeysColon(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	s, err := spec.ParseFile("Bulletspec")
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	// If user already typed "web:", don't complete further (instance number is freeform).
	if strings.Contains(toComplete, ":") {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	keys := make([]string, 0, len(s.Application.Programs))
	for k := range s.Application.Programs {
		keys = append(keys, k)
	}
	return keys, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
}

func completeProgramKeysEquals(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	s, err := spec.ParseFile("Bulletspec")
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	// If user already typed "web=", don't complete further (count is freeform).
	if strings.Contains(toComplete, "=") {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	keys := make([]string, 0, len(s.Application.Programs))
	for k := range s.Application.Programs {
		keys = append(keys, k+"=")
	}
	return keys, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
}

func completeCronJobKeys(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	s, err := spec.ParseFile("Bulletspec")
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	keys := make([]string, 0, len(s.Application.Cron.Jobs))
	for _, j := range s.Application.Cron.Jobs {
		keys = append(keys, j.Key)
	}
	return keys, cobra.ShellCompDirectiveNoFileComp
}
