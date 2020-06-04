package core

import (
	"log"

	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
)

func Restart(nodes []Node, spec *spec.Spec) error {
	for _, n := range nodes {
		log.Printf("Connecting to %s", n.Addr())
		c, err := ssh.Dial(n.Addr(), n.Identity)
		if err != nil {
			return err
		}

		d, err := distro.New(c)
		if err != nil {
			return err
		}

		log.Print("Restarting containers")
		for _, k := range spec.Application.ProgramKeys {
			p := spec.Application.Programs[k]
			err = d.RestartAll(spec.Application, p)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
