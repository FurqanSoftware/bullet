package core

import (
	"log"

	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
)

func Run(n Node, spec *spec.Spec, key string) error {
	log.Printf("Connecting to %s", n.Addr())
	c, err := ssh.Dial(n.Addr(), n.Identity)
	if err != nil {
		return err
	}

	d, err := distro.New(c)
	if err != nil {
		return err
	}

	prog, ok := spec.Application.Programs[key]
	if !ok {
		// TODO(hjr265): This should yield an error.
		return nil
	}

	return d.Run(spec.Application, prog)
}
