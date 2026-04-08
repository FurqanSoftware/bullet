package core

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/FurqanSoftware/bullet/cfg"
	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/scope"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
	"github.com/FurqanSoftware/pog"
	"github.com/antonmedv/expr"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

// Composition maps program keys to their desired instance counts.
type Composition struct {
	Sizes map[string]int
}

// NewComposition parses "key=count" arguments into a Composition.
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

// DefaultComposition evaluates the scaling rules from the spec against a node's properties.
func DefaultComposition(n scope.Node, spec *spec.Spec) (*Composition, error) {
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

type scaleResult struct {
	Desired int
	Up      int
	Down    int
}

// Scale adjusts the number of container instances on each node to match the composition.
func Scale(s scope.Scope, g cfg.Configuration, comp *Composition) error {
	type nodeResult struct {
		Name    string
		Results map[string]scaleResult
		Keys    []string
	}
	var nodeResults []nodeResult

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

		results, keys, err := scaleNode(n, c, d, s, comp)
		if err != nil {
			return err
		}

		nodeResults = append(nodeResults, nodeResult{
			Name:    n.Name,
			Results: results,
			Keys:    keys,
		})
	}

	// Collect all program keys in order.
	seen := map[string]bool{}
	var keys []string
	for _, nr := range nodeResults {
		for _, k := range nr.Keys {
			if !seen[k] {
				seen[k] = true
				keys = append(keys, k)
			}
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Configure(func(cfg *tablewriter.Config) {
		cfg.Header.Formatting.AutoFormat = tw.Off
		cfg.Footer.Alignment.Global = tw.AlignLeft
	})

	hdata := []any{""}
	for _, k := range keys {
		hdata = append(hdata, k)
	}
	table.Header(hdata...)

	desiredSum := map[string]int{}
	changeSum := map[string]int{}
	for _, nr := range nodeResults {
		rdata := []any{nr.Name}
		for _, k := range keys {
			r, ok := nr.Results[k]
			if !ok {
				rdata = append(rdata, "")
			} else {
				change := r.Up - r.Down
				desiredSum[k] += r.Desired
				changeSum[k] += change
				rdata = append(rdata, fmt.Sprintf("%d (%+d)", r.Desired, change))
			}
		}
		table.Append(rdata...)
	}

	fdata := []any{""}
	for _, k := range keys {
		fdata = append(fdata, fmt.Sprintf("%d (%+d)", desiredSum[k], changeSum[k]))
	}
	table.Footer(fdata...)

	fmt.Println()
	return table.Render()
}

func scaleNode(n scope.Node, c *ssh.Client, d distro.Distro, s scope.Scope, comp *Composition) (map[string]scaleResult, []string, error) {
	if len(comp.Sizes) == 0 {
		var err error
		comp, err = DefaultComposition(n, s.Spec)
		if err != nil {
			return nil, nil, err
		}
	}

	results := map[string]scaleResult{}
	var keys []string

	pog.SetStatus(pogText("Scaling programs"))
	for k, n := range comp.Sizes {
		prog, ok := s.Spec.Application.Programs[k]
		if !ok {
			return nil, nil, fmt.Errorf("unknown program key %q", k)
		}

		pog.SetStatus(pogScalingProgram(prog))
		up, down, err := d.Scale(s.Spec.Application, prog, n)
		if err != nil {
			return nil, nil, err
		}
		pog.Infof("Scaled program %s", k)
		pog.Infof("∟ Desired: %d", n)
		pog.Infof("∟ Ready: %d (%+d)", n, up-down)
		pog.SetStatus(nil)

		results[k] = scaleResult{
			Desired: n,
			Up:      up,
			Down:    down,
		}
		keys = append(keys, k)
	}
	return results, keys, nil
}
