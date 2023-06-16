package core

import (
	"log"
	"os"
	"text/tabwriter"

	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
	"github.com/FurqanSoftware/pog"
)

func CronEnable(nodes []Node, spec *spec.Spec, keys []string) error {
	for _, n := range nodes {
		log.Printf("Connecting to %s", n.Addr())
		c, err := ssh.Dial(n.Addr(), n.Identity)
		if err != nil {
			return err
		}

		d, err := distro.New(c)
		if err != nil {
			return err
		}

		log.Print("Enabling cron job(s)")
		for _, k := range keys {
			j := spec.Application.Cron.Job(k)
			if j.Command == "" {
				log.Fatalf("Bad job key %q", k)
			}
			err = d.CronEnable(spec.Application, j)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CronDisable(nodes []Node, spec *spec.Spec, keys []string) error {
	for _, n := range nodes {
		log.Printf("Connecting to %s", n.Addr())
		c, err := ssh.Dial(n.Addr(), n.Identity)
		if err != nil {
			return err
		}

		d, err := distro.New(c)
		if err != nil {
			return err
		}

		log.Print("Disabling cron job(s)")
		for _, k := range keys {
			j := spec.Application.Cron.Job(k)
			if j.Command == "" {
				log.Fatalf("Bad job key %q", k)
			}
			err = d.CronDisable(spec.Application, j)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CronStatus(nodes []Node, spec *spec.Spec, keys []string) error {
	for _, n := range nodes {
		pog.SetStatus(pogConnecting(n.Addr()))
		c, err := ssh.Dial(n.Addr(), n.Identity)
		if err != nil {
			return err
		}
		pog.Infof("Connected to %s", n.Addr())
		pog.SetStatus(nil)

		d, err := distro.New(c)
		if err != nil {
			return err
		}

		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

		for _, j := range spec.Application.Cron.Jobs {
			err = d.CronStatus(spec.Application, j, tw)
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
