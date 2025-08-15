package core

import (
	"fmt"
	"os"

	"github.com/FurqanSoftware/bullet/cfg"
	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/scope"
	"github.com/FurqanSoftware/pog"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

func Status(s scope.Scope, g cfg.Configuration) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.Configure(func(cfg *tablewriter.Config) {
		cfg.Header.Formatting.AutoFormat = tw.Off
		cfg.Footer.Alignment.Global = tw.AlignLeft
	})

	hdata := []any{""}
	for _, k := range s.Spec.Application.ProgramKeys {
		hdata = append(hdata, k)
	}
	table.Header(hdata...)

	upsum := map[string]int{}
	healthysum := map[string]int{}
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

		pog.SetStatus(pogText("Checking containers"))
		up := map[string]int{}
		healthy := map[string]int{}
		for _, k := range s.Spec.Application.ProgramKeys {
			p := s.Spec.Application.Programs[k]
			status, err := d.Status(s.Spec.Application, p)
			if err != nil {
				return err
			}
			for _, s := range status {
				if s.Up {
					up[k]++
					upsum[k]++
				}
				if s.Healthy {
					healthy[k]++
					healthysum[k]++
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
		for _, k := range s.Spec.Application.ProgramKeys {
			if up[k] == 0 {
				rdata = append(rdata, "0")
			} else {
				rdata = append(rdata, fmt.Sprintf("%d (%d)", up[k], healthy[k]))
			}
		}
		table.Append(rdata...)
	}

	fdata := []any{""}
	for _, k := range s.Spec.Application.ProgramKeys {
		if upsum[k] == 0 {
			fdata = append(fdata, "0")
		} else {
			fdata = append(fdata, fmt.Sprintf("%d (%d)", upsum[k], healthysum[k]))
		}
	}
	table.Footer(fdata...)

	return table.Render()
}
