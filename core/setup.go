package core

import (
	"fmt"
	"os"

	"github.com/FurqanSoftware/bullet/cfg"
	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/scope"
	"github.com/FurqanSoftware/bullet/ssh"
	"github.com/FurqanSoftware/pog"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

type setupResult struct {
	DockerInstalled bool
	EnvironPushed   bool
}

// Setup prepares servers for deployment by installing Docker and creating the application directory.
func Setup(s scope.Scope, g cfg.Configuration, environ string) error {
	results := map[string]setupResult{}

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

		r, err := setupNode(n, c, d, s)
		if err != nil {
			return err
		}

		if environ != "" {
			err = uploadEnvironFile(c, s, environ)
			if err != nil {
				return err
			}
			r.EnvironPushed = true
		}

		results[n.Name] = r
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Configure(func(cfg *tablewriter.Config) {
		cfg.Header.Formatting.AutoFormat = tw.Off
	})

	hdata := []any{"", "Docker", "Environ"}
	table.Header(hdata...)

	for _, n := range s.Nodes {
		r := results[n.Name]
		docker := "Already Installed"
		if r.DockerInstalled {
			docker = "Installed"
		}
		env := ""
		if r.EnvironPushed {
			env = "Pushed"
		}
		table.Append(n.Name, docker, env)
	}

	fmt.Println()
	return table.Render()
}

func setupNode(n scope.Node, c *ssh.Client, d distro.Distro, s scope.Scope) (setupResult, error) {
	var r setupResult

	pog.SetStatus(pogText("Installing Docker"))
	installed, err := d.InstallDocker()
	if err != nil {
		return r, err
	}
	r.DockerInstalled = installed
	pog.Info("Installed Docker")
	pog.SetStatus(nil)

	pog.SetStatus(pogText("Creating application directory"))
	err = d.MkdirAll(fmt.Sprintf("/opt/%s/releases", s.Spec.Application.Identifier))
	if err != nil {
		return r, err
	}
	err = d.Touch(fmt.Sprintf("/opt/%s/env", s.Spec.Application.Identifier))
	if err != nil {
		return r, err
	}
	pog.Info("Created application directory")
	pog.SetStatus(nil)

	return r, nil
}
