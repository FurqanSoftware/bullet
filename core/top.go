package core

import (
	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
	"github.com/FurqanSoftware/pog"
)

func Top(node Node, spec *spec.Spec) error {
	pog.SetStatus(pogConnecting(node))
	c, err := ssh.Dial(node.Addr(), node.Identity)
	if err != nil {
		return err
	}
	pog.Infof("Connected to %s", node.Label())
	pog.SetStatus(nil)

	d, err := distro.New(c)
	if err != nil {
		return err
	}

	return d.Top()
}
