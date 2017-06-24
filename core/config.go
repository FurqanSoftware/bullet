package core

import (
	"fmt"
	"log"
	"os"

	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
)

func ConfigPush(nodes []Node, spec *spec.Spec, name string) error {
	for _, n := range nodes {
		log.Printf("Connecting to %s", n.Addr())
		c, err := ssh.Dial(n.Addr())
		if err != nil {
			return err
		}

		log.Print("Uploading environment file")
		f, err := os.Open(name)
		if err != nil {
			return err
		}
		defer f.Close()
		s, err := f.Stat()
		if err != nil {
			return err
		}
		err = c.Push(fmt.Sprintf("/opt/%s/env", spec.Application.Identifier), s.Mode(), s.Size(), f)
		if err != nil {
			return err
		}
	}
	return nil
}
