package core

import (
	"fmt"
	"log"

	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
)

func Prune(nodes []Node, spec *spec.Spec) error {
	for _, n := range nodes {
		log.Printf("Connecting to %s", n.Label())
		c, err := ssh.Dial(n.Addr(), n.Identity)
		if err != nil {
			return err
		}

		d, err := distro.New(c)
		if err != nil {
			return err
		}

		log.Print("Removing stale releases")
		err = d.Prune(fmt.Sprintf("/opt/%s/releases", spec.Application.Identifier), 5)
		if err != nil {
			return err
		}
	}
	return nil
}
