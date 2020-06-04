package core

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
)

func Status(nodes []Node, spec *spec.Spec) error {
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

		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

		fmt.Println("[Programs]")
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
		fmt.Println()

		fmt.Println("[Cron]")
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
		fmt.Println()
	}
	return nil
}
