package main

import "github.com/spf13/cobra"

var (
	Hosts    string
	Port     int
	Identity string
)

var RootCmd = &cobra.Command{
	Use:   "bullet",
	Short: "Bullet is a fast application deploy tool",
	Long:  `Bullet is a fast and flexible application deploy tool built by Furqan Software and friends. Complete documentation is available at https://bullettool.com/.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&Hosts, "hosts", "H", "", "List of target hosts (comma separated)")
	RootCmd.PersistentFlags().IntVarP(&Port, "port", "p", 22, "Port to connect to")
	RootCmd.PersistentFlags().StringVarP(&Identity, "identitiy", "i", "", "Path to an SSH identity file (usually a private key)")
}
