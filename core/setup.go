package core

import (
	"fmt"

	"github.com/FurqanSoftware/bullet/cfg"
	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/scope"
	"github.com/FurqanSoftware/bullet/ssh"
	"github.com/FurqanSoftware/pog"
)

func Setup(s scope.Scope, g cfg.Configuration, environ string) error {
	for _, n := range s.Nodes {
		pog.Infof("Connecting to %s", n.Label())
		c, err := sshDial(n, g)
		if err != nil {
			return err
		}

		d, err := distro.New(c)
		if err != nil {
			return err
		}

		err = setupNode(n, c, d, s)
		if err != nil {
			return err
		}

		if environ != "" {
			err = uploadEnvironFile(c, s, environ)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func setupNode(n scope.Node, c *ssh.Client, d distro.Distro, s scope.Scope) error {
	pog.Info("Installing Docker")
	err := d.InstallDocker()
	if err != nil {
		return err
	}

	pog.Info("Creating application directory")
	err = d.MkdirAll(fmt.Sprintf("/opt/%s/releases", s.Spec.Application.Identifier))
	if err != nil {
		return err
	}
	err = d.Touch(fmt.Sprintf("/opt/%s/env", s.Spec.Application.Identifier))
	if err != nil {
		return err
	}

	return nil
}
