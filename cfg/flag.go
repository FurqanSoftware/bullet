package cfg

import "github.com/spf13/pflag"

// AddFlags registers the global CLI flags (hosts, port, identity) on the given flag set.
func AddFlags(flags *pflag.FlagSet) {
	flags.StringP("hosts", "H", "", "List of target hosts (comma separated)")
	flags.IntP("port", "p", 22, "Port to connect to")
	flags.StringP("identity", "i", "", "Path to an SSH identity file (usually a private key)")
}
