package core

import (
	"fmt"
	"log"

	"github.com/FurqanSoftware/bullet/cfg"
	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/scope"
)

func Prune(s scope.Scope, g cfg.Configuration) error {
	for _, n := range s.Nodes {
		log.Printf("Connecting to %s", n.Label())
		c, err := sshDial(n, g)
		if err != nil {
			return err
		}

		d, err := distro.New(c)
		if err != nil {
			return err
		}

		log.Print("Removing stale releases")
		err = d.Prune(fmt.Sprintf("/opt/%s/releases", s.Spec.Application.Identifier), 5)
		if err != nil {
			return err
		}
	}
	return nil
}
