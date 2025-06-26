package core

import (
	"strconv"
	"strings"

	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
	"github.com/FurqanSoftware/pog"
	"github.com/antonmedv/expr"
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

func DefaultComposition(n Node, spec *spec.Spec) (*Composition, error) {
	comp := Composition{
		Sizes: make(map[string]int, len(spec.Application.Programs)),
	}
	for key, prog := range spec.Application.Programs {
		for _, scale := range prog.Scales {
			env := map[string]interface{}{
				"hasTags": func(tags ...string) bool { return n.HasTags(tags) },
				"hw": map[string]interface{}{
					"cores":  n.HW.Cores,
					"memory": n.HW.Memory,
				},
			}
			if scale.If != "" {
				prog, err := expr.Compile(scale.If, expr.Env(env))
				if err != nil {
					return nil, err
				}
				cond, err := expr.Run(prog, env)
				if err != nil {
					return nil, err
				}
				if !cond.(bool) {
					continue
				}
			}
			prog, err := expr.Compile(scale.N, expr.Env(env))
			if err != nil {
				return nil, err
			}
			n, err := expr.Run(prog, env)
			if err != nil {
				return nil, err
			}
			comp.Sizes[key] = n.(int)
		}
	}
	return &comp, nil
}

func Scale(nodes []Node, spec *spec.Spec, comp *Composition) error {
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

		err = scaleNode(n, c, d, spec, comp)
		if err != nil {
			return err
		}
	}
	return nil
}

func scaleNode(n Node, c *ssh.Client, d distro.Distro, spec *spec.Spec, comp *Composition) error {
	if len(comp.Sizes) == 0 {
		var err error
		comp, err = DefaultComposition(n, spec)
		if err != nil {
			return err
		}
	}

	pog.SetStatus(pogText("Scaling programs"))
	for k, n := range comp.Sizes {
		prog, ok := spec.Application.Programs[k]
		if !ok {
			// TODO(hjr265): This should yield an error.
			continue
		}

		pog.SetStatus(pogScalingProgram(prog))
		up, down, err := d.Scale(spec.Application, prog, n)
		if err != nil {
			return err
		}
		pog.Infof("Scaled program %s", k)
		pog.Infof("∟ Desired: %d", n)
		pog.Infof("∟ Ready: %d (%+d)", n, up-down)
		pog.SetStatus(nil)
	}
	return nil
}
