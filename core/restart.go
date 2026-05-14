package core

import (
	"fmt"

	"github.com/FurqanSoftware/bullet/cfg"
	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/scope"
	"github.com/FurqanSoftware/pog"
)

// RestartTarget selects containers to restart. No == 0 means every instance
// of Key; otherwise only that specific instance.
type RestartTarget struct {
	Key string
	No  int
}

// Restart restarts application containers on the selected nodes. If targets
// is empty, every container of every program is restarted.
func Restart(s scope.Scope, g cfg.Configuration, targets []RestartTarget) error {
	for _, t := range targets {
		if _, ok := s.Spec.Application.Programs[t.Key]; !ok {
			return fmt.Errorf("unknown program key %q", t.Key)
		}
	}

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

		_, err = restartNode(d, s, targets)
		if err != nil {
			return err
		}

		if len(targets) == 0 {
			hash, err := hashRemoteEnviron(d, s)
			if err == nil {
				err = markEnvironApplied(d, s, hash)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// restartNode restarts the containers on a node that match targets and returns
// the number restarted per program key. If targets is empty, every container is
// restarted.
func restartNode(d distro.Distro, s scope.Scope, targets []RestartTarget) (map[string]int, error) {
	pog.SetStatus(pogText("Restarting containers"))
	nrestart := map[string]int{}
	nrestartsum := 0
	for _, k := range s.Spec.Application.ProgramKeys {
		if !programTargeted(targets, k) {
			continue
		}
		p := s.Spec.Application.Programs[k]
		statuses, err := d.Status(s.Spec.Application, p)
		if err != nil {
			return nrestart, err
		}
		for _, status := range statuses {
			if status.No == 0 {
				continue
			}
			if !instanceTargeted(targets, k, status.No) {
				continue
			}
			pog.SetStatus(pogRestartingContainer(p, status.No))
			err = d.Restart(s.Spec.Application, p, status.No)
			if err != nil {
				return nrestart, err
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
	return nrestart, nil
}

func programTargeted(targets []RestartTarget, key string) bool {
	if len(targets) == 0 {
		return true
	}
	for _, t := range targets {
		if t.Key == key {
			return true
		}
	}
	return false
}

func instanceTargeted(targets []RestartTarget, key string, no int) bool {
	if len(targets) == 0 {
		return true
	}
	for _, t := range targets {
		if t.Key == key && (t.No == 0 || t.No == no) {
			return true
		}
	}
	return false
}
