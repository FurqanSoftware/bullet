package core

import (
	"github.com/FurqanSoftware/bullet/cfg"
	"github.com/FurqanSoftware/bullet/distro"
	"github.com/FurqanSoftware/bullet/scope"
	"github.com/FurqanSoftware/bullet/ssh"
	"github.com/FurqanSoftware/pog"
)

func sshDial(n scope.Node, g cfg.Configuration) (*ssh.Client, error) {
	return ssh.Dial(
		n.Addr(),
		n.Identity,
		ssh.WithRetries(g.SSHRetries),
		ssh.WithTimeout(g.SSHTimeout),
	)
}

// SSH opens an interactive SSH session to the selected node.
func SSH(s scope.Scope, g cfg.Configuration) error {
	pog.SetStatus(pogConnecting(s.Nodes[0]))
	c, err := sshDial(s.Nodes[0], g)
	if err != nil {
		return err
	}
	pog.Infof("Connected to %s", s.Nodes[0].Label())
	pog.SetStatus(nil)

	d, err := distro.New(c)
	if err != nil {
		return err
	}

	return d.Shell()
}
