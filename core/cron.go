package core

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/FurqanSoftware/bullet/cfg"
	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/scope"
	"github.com/FurqanSoftware/pog"
)

// CronEnable enables the specified cron jobs on all selected nodes.
func CronEnable(s scope.Scope, g cfg.Configuration, keys []string) error {
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

		pog.Info("Enabling cron job(s)")
		for _, k := range keys {
			j := s.Spec.Application.Cron.Job(k)
			if j.Command == "" {
				return fmt.Errorf("bad job key %q", k)
			}
			err = d.CronEnable(s.Spec.Application, j)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CronDisable disables the specified cron jobs on all selected nodes.
func CronDisable(s scope.Scope, g cfg.Configuration, keys []string) error {
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

		pog.Info("Disabling cron job(s)")
		for _, k := range keys {
			j := s.Spec.Application.Cron.Job(k)
			if j.Command == "" {
				return fmt.Errorf("bad job key %q", k)
			}
			err = d.CronDisable(s.Spec.Application, j)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CronStatus prints the status of all cron jobs on the selected nodes.
func CronStatus(s scope.Scope, g cfg.Configuration, keys []string) error {
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

		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

		for _, j := range s.Spec.Application.Cron.Jobs {
			err = d.CronStatus(s.Spec.Application, j, tw)
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
