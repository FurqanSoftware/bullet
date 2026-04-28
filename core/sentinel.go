package core

import (
	"github.com/FurqanSoftware/bullet/cfg"
	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/scope"
	"github.com/FurqanSoftware/pog"
)

// SentinelBinaryFunc returns the embedded sentinel binary for the given
// architecture (e.g. "amd64", "arm64"). Implemented by the main package, which
// owns the embed directives.
type SentinelBinaryFunc func(arch string) ([]byte, error)

// SentinelInstall pushes the embedded sentinel binary and systemd unit to each
// selected node and starts the service.
func SentinelInstall(s scope.Scope, g cfg.Configuration, binaryFn SentinelBinaryFunc) error {
	for _, n := range s.Nodes {
		pog.SetStatus(pogConnecting(n))
		c, err := sshDial(n, g)
		if err != nil {
			return err
		}
		pog.Infof("Connected to %s", n.Label())
		pog.SetStatus(nil)

		d, err := distro.New(c)
		if err != nil {
			return err
		}

		pog.SetStatus(pogText("Detecting architecture"))
		arch, err := d.HostArch()
		if err != nil {
			return err
		}
		pog.Infof("Architecture: %s", arch)
		pog.SetStatus(nil)

		bin, err := binaryFn(arch)
		if err != nil {
			return err
		}

		pog.SetStatus(pogText("Installing sentinel"))
		if err := d.InstallSentinel(bin); err != nil {
			return err
		}
		pog.Info("Installed sentinel")
		pog.SetStatus(nil)
	}
	return nil
}
