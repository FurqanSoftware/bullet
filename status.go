package bullet

import (
	"log"

	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
)

func Status(nodes []Node, spec *spec.Spec) error {
	for _, n := range nodes {
		log.Printf("Connecting to %s", n.Addr())
		c, err := ssh.Dial(n.Addr())
		if err != nil {
			return err
		}

		d, err := distro.New(c)
		if err != nil {
			return err
		}

		for _, p := range spec.Application.Programs {
			err = d.Status(spec.Application, p)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
