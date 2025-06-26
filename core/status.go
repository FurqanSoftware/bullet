package core

import (
	"fmt"
	"os"

	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
	"github.com/FurqanSoftware/pog"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

func Status(nodes []Node, spec *spec.Spec) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.Configure(func(cfg *tablewriter.Config) {
		cfg.Header.Formatting.AutoFormat = tw.Off
	})

	hdata := []any{""}
	for _, k := range spec.Application.ProgramKeys {
		hdata = append(hdata, k)
	}
	table.Header(hdata...)

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

		pog.SetStatus(pogText("Checking containers"))
		up := map[string]int{}
		healthy := map[string]int{}
		for _, k := range spec.Application.ProgramKeys {
			p := spec.Application.Programs[k]
			status, err := d.Status(spec.Application, p)
			if err != nil {
				return err
			}
			for _, s := range status {
				if s.Up {
					up[k]++
				}
				if s.Healthy {
					healthy[k]++
				}
			}
			if up[k] > 0 {
				pog.Infof("%s (%s)", p.Name, k)
				pog.Infof("∟ Running: %d", up[k])
				pog.Infof("∟ Healthy: %d", healthy[k])
			}
		}
		pog.SetStatus(nil)

		rdata := []any{n.Name}
		for _, k := range spec.Application.ProgramKeys {
			if up[k] == 0 {
				rdata = append(rdata, "0")
			} else {
				rdata = append(rdata, fmt.Sprintf("%d (%d)", up[k], healthy[k]))
			}
		}
		table.Append(rdata...)
	}

	return table.Render()
}
