package core

import (
	"log"

	"github.com/FurqanSoftware/bullet/cfg"
	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/scope"
)

func Log(s scope.Scope, g cfg.Configuration, key string, no int) error {
	log.Printf("Connecting to %s", s.Nodes[0].Label())
	c, err := sshDial(s.Nodes[0], g)
	if err != nil {
		return err
	}

	d, err := distro.New(c)
	if err != nil {
		return err
	}

	prog, ok := s.Spec.Application.Programs[key]
	if !ok {
		// TODO(hjr265): This should yield an error.
		return nil
	}

	return d.Log(s.Spec.Application, prog, no)
}
