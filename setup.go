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
		log.Printf("Connecting to %s", n.Addr())
		c, err := ssh.Dial(n.Addr())
		if err != nil {
			return err
		}

		d, err := distro.New(c)
		if err != nil {
			return err
		}

		err = setupNode(n, c, d, spec)
		if err != nil {
			return err
		}
	}
	return nil
}

func setupNode(n Node, c *ssh.Client, d distro.Distro, spec *spec.Spec) error {
	log.Print("Installing Docker")
	err := d.InstallDocker()
	if err != nil {
		return err
	}

	log.Print("Creating application directory")
	err = d.MkdirAll(fmt.Sprintf("/opt/%s/releases", spec.Application.Identifier))
	if err != nil {
		return err
	}
	err = d.Touch(fmt.Sprintf("/opt/%s/env", spec.Application.Identifier))
	if err != nil {
		return err
	}

	return nil
}
