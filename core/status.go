package core

import (
	"os"
	"text/tabwriter"

	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
	"github.com/FurqanSoftware/pog"
)

func Status(nodes []Node, spec *spec.Spec) error {
	for _, n := range nodes {
		pog.SetStatus(pogConnecting(n))
		c, err := ssh.Dial(n.Addr(), n.Identity)
		if err != nil {
			return err
		}
		pog.Infof("Connected to %s", n.Label())
		pog.SetStatus(nil)

		d, err := distro.New(c)
		if err != nil {
			return err
		}

		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

		pog.Infof("Checking containers")
		for _, k := range spec.Application.ProgramKeys {
			p := spec.Application.Programs[k]
			err = d.Status(spec.Application, p, tw)
			if err != nil {
				return err
			}
		}
		err = tw.Flush()
		if err != nil {
			return err
		}
	}
	return nil
}
