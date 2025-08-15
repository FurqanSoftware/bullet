package core

import (
	"fmt"
	"os"

	"github.com/FurqanSoftware/bullet/cfg"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/scope"
	"github.com/FurqanSoftware/bullet/ssh"
	"github.com/FurqanSoftware/pog"
)

func EnvironPush(s scope.Scope, g cfg.Configuration, filename string) error {
	for _, n := range s.Nodes {
		pog.Infof("Connecting to %s", n.Label())
		c, err := sshDial(n, g)
		if err != nil {
			return err
		}

		err = uploadEnvironFile(c, s, filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func uploadEnvironFile(c *ssh.Client, s scope.Scope, filename string) error {
	pog.Infof("Uploading environment file")
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	return c.Push(fmt.Sprintf("/opt/%s/env", s.Spec.Application.Identifier), fi.Mode(), fi.Size(), f, nil)
}
