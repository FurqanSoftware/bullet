package core

import (
	"fmt"

	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
	"github.com/FurqanSoftware/pog"
)

func Run(nodes []Node, spec *spec.Spec, key string) error {
	var i int = 1
	if len(nodes) > 1 {
		for i, n := range nodes {
			fmt.Printf("%d. %s\n", i+1, n.Label())
		}
		fmt.Print("? ")
		fmt.Scanf("%d", &i)
	}

	n := nodes[i-1]

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

	prog, ok := spec.Application.Programs[key]
	if !ok {
		// TODO(hjr265): This should yield an error.
		return nil
	}

	return d.Run(spec.Application, prog)
}
