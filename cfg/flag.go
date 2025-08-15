package cfg

import "github.com/spf13/pflag"

func AddFlags(flags *pflag.FlagSet) {
	flags.StringP("hosts", "H", "", "List of target hosts (comma separated)")
	flags.IntP("port", "p", 22, "Port to connect to")
	flags.StringP("identity", "i", "", "Path to an SSH identity file (usually a private key)")
}
