package core

import (
	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
	"github.com/FurqanSoftware/pog"
)

func Forward(nodes []Node, spec *spec.Spec, port string) error {
	n := SelectNode(nodes)

	pog.Infof("Connecting to %s", n.Label())
	c, err := ssh.Dial(n.Addr(), n.Identity)
	if err != nil {
		return err
	}

	d, err := distro.New(c)
	if err != nil {
		return err
	}

	return d.Forward(spec.Application, port)
}
