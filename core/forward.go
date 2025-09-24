package core

import (
	"github.com/FurqanSoftware/bullet/cfg"
	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/scope"
	"github.com/FurqanSoftware/pog"
)

func Forward(s scope.Scope, g cfg.Configuration, port string) error {
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

	return d.Forward(s.Spec.Application, port)
}
