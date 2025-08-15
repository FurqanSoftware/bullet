package core

import (
	"github.com/FurqanSoftware/bullet/cfg"
	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/scope"
	"github.com/FurqanSoftware/pog"
)

func Restart(s scope.Scope, g cfg.Configuration) error {
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

		pog.SetStatus(pogText("Restarting containers"))
		nrestart := map[string]int{}
		nrestartsum := 0
		for _, k := range s.Spec.Application.ProgramKeys {
			p := s.Spec.Application.Programs[k]
			statuses, err := d.Status(s.Spec.Application, p)
			if err != nil {
				return err
			}
			for _, status := range statuses {
				if status.No == 0 {
					continue
				}
				pog.SetStatus(pogRestartingContainer(p, status.No))
				err = d.Restart(s.Spec.Application, p, status.No)
				if err != nil {
					return err
				}
				nrestart[k]++
				nrestartsum++
			}
		}
		pog.Infof("Restarted %d container(s)", nrestartsum)
		for _, k := range s.Spec.Application.ProgramKeys {
			if nrestart[k] == 0 {
				continue
			}
			pog.Infof("∟ %s: %d", k, nrestart[k])
		}
		pog.SetStatus(nil)
	}
	return nil
}
