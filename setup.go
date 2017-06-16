package bullet

import (
	"fmt"
	"log"

	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
)

func Setup(nodes []Node, spec *spec.Spec) error {
	for _, n := range nodes {
		log.Print("Connecting to %s", n.Addr())
		c, err := ssh.Dial(n.Addr())
		if err != nil {
			return err
		}

		d, err := distro.New(c)
		if err != nil {
			return err
		}
		log.Print("Installing Docker")
		err = d.InstallDocker()
		if err != nil {
			return err
		}

		log.Print("Creating directories")
		err = d.MkdirAll(fmt.Sprintf("/opt/%s/releases", spec.Application.Identifier))
		if err != nil {
			return err
		}
	}
	return nil
}
