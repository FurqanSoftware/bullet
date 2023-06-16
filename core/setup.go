package core

import (
	"fmt"

	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
	"github.com/FurqanSoftware/pog"
)

func Setup(nodes []Node, spec *spec.Spec, config string) error {
	for _, n := range nodes {
		pog.Infof("Connecting to %s", n.Addr())
		c, err := ssh.Dial(n.Addr(), n.Identity)
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

		if config != "" {
			err = uploadEnvironmentFile(c, spec, config)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func setupNode(n Node, c *ssh.Client, d distro.Distro, spec *spec.Spec) error {
	pog.Info("Installing Docker")
	err := d.InstallDocker()
	if err != nil {
		return err
	}

	pog.Info("Creating application directory")
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
