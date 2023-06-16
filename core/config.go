package core

import (
	"fmt"
	"os"

	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
	"github.com/FurqanSoftware/pog"
)

func ConfigPush(nodes []Node, spec *spec.Spec, filename string) error {
	for _, n := range nodes {
		pog.Infof("Connecting to %s", n.Label())
		c, err := ssh.Dial(n.Addr(), n.Identity)
		if err != nil {
			return err
		}

		err = uploadEnvironmentFile(c, spec, filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func uploadEnvironmentFile(c *ssh.Client, spec *spec.Spec, filename string) error {
	pog.Infof("Uploading environment file")
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	s, err := f.Stat()
	if err != nil {
		return err
	}
	return c.Push(fmt.Sprintf("/opt/%s/env", spec.Application.Identifier), s.Mode(), s.Size(), f)
}
