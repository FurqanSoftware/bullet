package core

import (
	"github.com/FurqanSoftware/bullet/cfg"
	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/scope"
	"github.com/FurqanSoftware/pog"
)

func Forward(s scope.Scope, g cfg.Configuration, port string) error {
	pog.Infof("Connecting to %s", s.Nodes[0].Label())
	c, err := sshDial(s.Nodes[0], g)
	if err != nil {
		return err
	}

	d, err := distro.New(c)
	if err != nil {
		return err
	}

	return d.Forward(s.Spec.Application, port)
}
