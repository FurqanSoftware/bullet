package core

import (
	"fmt"
	"log"

	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
)

func Log(nodes []Node, spec *spec.Spec, key string, no int) error {
	var i int = 1
	if len(nodes) > 1 {
		for i, n := range nodes {
			fmt.Printf("%d. %s\n", i+1, n.Label())
		}
		fmt.Print("? ")
		fmt.Scanf("%d", &i)
	}

	n := nodes[i-1]

	log.Printf("Connecting to %s", n.Label())
	c, err := ssh.Dial(n.Addr(), n.Identity)
	if err != nil {
		return err
	}

	d, err := distro.New(c)
	if err != nil {
		return err
	}

	prog, ok := spec.Application.Programs[key]
	if !ok {
		// TODO(hjr265): This should yield an error.
		return nil
	}

	return d.Log(spec.Application, prog, no)
}
