package bullet

import (
	"log"

	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
)

func Install(nodes []Node, spec *spec.Spec) error {
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

		log.Printf("Installing %d programs(s)", len(spec.Application.Programs))
		for _, p := range spec.Application.Programs {
			log.Printf("Installing %s", p.Name)
			err = d.Install(spec.Application, p)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
