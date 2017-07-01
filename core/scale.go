package core

import (
	"log"
	"strconv"
	"strings"

	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
)

type Composition struct {
	Sizes map[string]int
}

func NewComposition(args []string) (*Composition, error) {
	comp := Composition{
		Sizes: map[string]int{},
	}
	for _, a := range args {
		parts := strings.SplitN(a, "=", 2)
		n, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		comp.Sizes[parts[0]] = n
	}
	return &comp, nil
}

func Scale(nodes []Node, spec *spec.Spec, comp *Composition) error {
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

		scaleNode(n, c, d, spec, comp)
	}
	return nil
}

func scaleNode(n Node, c *ssh.Client, d distro.Distro, spec *spec.Spec, comp *Composition) error {
	for key, n := range comp.Sizes {
		prog, ok := spec.Application.Programs[key]
		if !ok {
			// TODO(hjr265): This should yield an error.
			continue
		}

		err := d.Scale(spec.Application, prog, n)
		if err != nil {
			return err
		}
	}
	return nil
}
