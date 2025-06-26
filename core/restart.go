package core

import (
	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
	"github.com/FurqanSoftware/pog"
)

func Restart(nodes []Node, spec *spec.Spec) error {
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

		pog.SetStatus(pogText("Restarting containers"))
		nrestart := map[string]int{}
		nrestartsum := 0
		for _, k := range spec.Application.ProgramKeys {
			p := spec.Application.Programs[k]
			status, err := d.Status(spec.Application, p)
			if err != nil {
				return err
			}
			for _, s := range status {
				if s.No == 0 {
					continue
				}
				pog.SetStatus(pogRestartingContainer(p, s.No))
				err = d.Restart(spec.Application, p, s.No)
				if err != nil {
					return err
				}
				nrestart[k]++
				nrestartsum++
			}
		}
		pog.Infof("Restarted %d container(s)", nrestartsum)
		for _, k := range spec.Application.ProgramKeys {
			if nrestart[k] == 0 {
				continue
			}
			pog.Infof("∟ %s: %d", k, nrestart[k])
		}
		pog.SetStatus(nil)
	}
	return nil
}
