package core

import (
	"github.com/FurqanSoftware/bullet/cfg"
	"github.com/FurqanSoftware/bullet/scope"
	"github.com/FurqanSoftware/bullet/ssh"
)

func sshDial(n scope.Node, g cfg.Configuration) (*ssh.Client, error) {
	return ssh.Dial(
		n.Addr(),
		n.Identity,
		ssh.WithRetries(g.SSHRetries),
		ssh.WithTimeout(g.SSHTimeout),
	)
}
