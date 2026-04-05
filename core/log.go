package core

import (
	"fmt"

	"github.com/FurqanSoftware/bullet/cfg"
	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/scope"
	"github.com/FurqanSoftware/pog"
)

// Log tails the logs of a program's container instance.
func Log(s scope.Scope, g cfg.Configuration, key string, no int) error {
	pog.SetStatus(pogConnecting(s.Nodes[0]))
	c, err := sshDial(s.Nodes[0], g)
	if err != nil {
		return err
	}
	pog.Infof("Connected to %s", s.Nodes[0].Label())
	pog.SetStatus(nil)

	d, err := distro.New(c)
	if err != nil {
		return err
	}

	prog, ok := s.Spec.Application.Programs[key]
	if !ok {
		return fmt.Errorf("unknown program key %q", key)
	}

	return d.Log(s.Spec.Application, prog, no)
}
