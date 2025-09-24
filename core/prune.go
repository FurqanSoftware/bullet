package core

import (
	"fmt"
	"log"

	"github.com/FurqanSoftware/bullet/cfg"
	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/scope"
	"github.com/FurqanSoftware/pog"
)

func Prune(s scope.Scope, g cfg.Configuration) error {
	for _, n := range s.Nodes {
		pog.SetStatus(pogConnecting(n))
		c, err := sshDial(n, g)
		if err != nil {
			return err
		}
		pog.Infof("Connected to %s", n.Label())
		pog.SetStatus(nil)

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
