package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/FurqanSoftware/bullet/cfg"
	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/scope"
	"github.com/FurqanSoftware/bullet/ssh"
	"github.com/FurqanSoftware/pog"
	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

type deployResult struct {
	Skipped       bool
	EnvironPushed bool
	Rebuilt       map[string]bool
	Reloaded      map[string]int
	Scaled        map[string]scaleResult
}

// Deploy uploads and deploys a release to all nodes in scope.
func Deploy(s scope.Scope, g cfg.Configuration, rel *Release, environ string, scale bool) error {
	pog.Infof("Deploying %s", filepath.Base(rel.Tarball.Path))
	pog.Infof("∟ Hash: %s", rel.Hash)
	pog.Infof("∟ Size: %s", humanize.Bytes(uint64(rel.Tarball.Size)))

	results := map[string]deployResult{}

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

		var environPushed bool
		if environ != "" {
			err = uploadEnvironFile(c, s, environ)
			if err != nil {
				return err
			}
			environPushed = true
		}

		r, err := deployNode(n, c, d, s, rel)
		if err != nil {
			return err
		}
		r.EnvironPushed = environPushed

		if scale && !r.Skipped {
			comp := &Composition{Sizes: map[string]int{}}
			scaled, _, err := scaleNode(n, c, d, s, comp)
			if err != nil {
				return err
			}
			r.Scaled = scaled
		}

		results[n.Name] = r
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Configure(func(cfg *tablewriter.Config) {
		cfg.Header.Formatting.AutoFormat = tw.Off
		cfg.Footer.Alignment.Global = tw.AlignLeft
	})

	hdata := []any{"", "Environ"}
	for _, k := range s.Spec.Application.ProgramKeys {
		hdata = append(hdata, k)
	}
	table.Header(hdata...)

	rebuiltSum := map[string]int{}
	reloadedSum := map[string]int{}
	scaledSum := map[string]scaleResult{}
	for _, n := range s.Nodes {
		r := results[n.Name]
		envCell := ""
		if r.EnvironPushed {
			envCell = "Pushed"
		}
		rdata := []any{n.Name, envCell}
		if r.Skipped {
			for range s.Spec.Application.ProgramKeys {
				rdata = append(rdata, "Skipped")
			}
		} else {
			for _, k := range s.Spec.Application.ProgramKeys {
				cell := ""
				if r.Rebuilt[k] {
					rebuiltSum[k]++
					cell = "Rebuilt"
				}
				if r.Reloaded[k] > 0 {
					reloadedSum[k] += r.Reloaded[k]
					if cell != "" {
						cell += ", "
					}
					cell += fmt.Sprintf("%d Reloaded", r.Reloaded[k])
				}
				if sr, ok := r.Scaled[k]; ok {
					change := sr.Up - sr.Down
					sum := scaledSum[k]
					sum.Desired += sr.Desired
					sum.Up += sr.Up
					sum.Down += sr.Down
					scaledSum[k] = sum
					if cell != "" {
						cell += ", "
					}
					cell += fmt.Sprintf("%d Scaled (%+d)", sr.Desired, change)
				}
				if cell == "" {
					cell = "0"
				}
				rdata = append(rdata, cell)
			}
		}
		table.Append(rdata...)
	}

	fdata := []any{"", ""}
	for _, k := range s.Spec.Application.ProgramKeys {
		cell := ""
		if rebuiltSum[k] > 0 {
			cell = fmt.Sprintf("%d Rebuilt", rebuiltSum[k])
		}
		if reloadedSum[k] > 0 {
			if cell != "" {
				cell += ", "
			}
			cell += fmt.Sprintf("%d Reloaded", reloadedSum[k])
		}
		if sum, ok := scaledSum[k]; ok {
			change := sum.Up - sum.Down
			if cell != "" {
				cell += ", "
			}
			cell += fmt.Sprintf("%d Scaled (%+d)", sum.Desired, change)
		}
		if cell == "" {
			cell = "0"
		}
		fdata = append(fdata, cell)
	}
	table.Footer(fdata...)

	fmt.Println()
	return table.Render()
}

func deployNode(n scope.Node, c *ssh.Client, d distro.Distro, s scope.Scope, rel *Release) (deployResult, error) {
	r := deployResult{
		Rebuilt:  map[string]bool{},
		Reloaded: map[string]int{},
	}

	curHash, _ := d.ReadFile(fmt.Sprintf("/opt/%s/current.hash", s.Spec.Application.Identifier))
	if rel.Hash == string(curHash) {
		pog.Info("Same as current hash. Skipping deploy.")
		r.Skipped = true
		return r, nil
	}

	tarPath := fmt.Sprintf("/tmp/%s-%s-%s.tar.gz", s.Spec.Application.Identifier, rel.Time, rel.Hash)
	err := uploadTarball(c, tarPath, rel.Tarball)
	if err != nil {
		return r, err
	}
	pog.Info("Uploaded tarball")
	pog.SetStatus(nil)

	pog.SetStatus(pogText("Extracting tarball"))
	relDir := fmt.Sprintf("/opt/%s/releases/%s-%s", s.Spec.Application.Identifier, rel.Time, rel.Hash)
	err = d.ExtractTar(tarPath, relDir)
	if err != nil {
		return r, err
	}
	pog.Info("Extracted tarball")
	pog.SetStatus(nil)

	pog.SetStatus(pogText("Removing tarball"))
	err = d.Remove(tarPath)
	if err != nil {
		return r, err
	}
	pog.Info("Removed tarball")
	pog.SetStatus(nil)

	pog.SetStatus(pogText("Updating current"))
	err = d.UpdateCurrent(s.Spec.Application, relDir)
	if err != nil {
		return r, err
	}
	pog.Info("Updated current")
	pog.SetStatus(nil)

	err = d.WriteFile(fmt.Sprintf("/opt/%s/current.hash", s.Spec.Application.Identifier), []byte(rel.Hash))
	if err != nil {
		return r, err
	}

	pog.SetStatus(pogText("Building images"))
	for _, p := range s.Spec.Application.Programs {
		rebuilt, err := d.Build(s.Spec.Application, p)
		if err != nil {
			return r, err
		}
		if rebuilt {
			r.Rebuilt[p.Key] = true
		}
	}
	pog.Infof("Built %d image(s)", len(r.Rebuilt))
	for _, k := range s.Spec.Application.ProgramKeys {
		if !r.Rebuilt[k] {
			continue
		}
		pog.Infof("∟ %s", k)
	}
	pog.SetStatus(nil)

	pog.SetStatus(pogText("Reloading containers"))
	nreloadsum := 0
	for _, k := range s.Spec.Application.ProgramKeys {
		p := s.Spec.Application.Programs[k]
		statuses, err := d.Status(s.Spec.Application, p)
		if err != nil {
			return r, err
		}
		for _, status := range statuses {
			if status.No == 0 {
				continue
			}
			pog.SetStatus(pogReloadingContainer(p, status.No))
			err = d.Reload(s.Spec.Application, p, status.No, r.Rebuilt[k])
			if err != nil {
				return r, err
			}
			r.Reloaded[k]++
			nreloadsum++
		}
	}
	pog.Infof("Reloaded %d container(s)", nreloadsum)
	for _, k := range s.Spec.Application.ProgramKeys {
		if r.Reloaded[k] == 0 {
			continue
		}
		pog.Infof("∟ %s: %d", k, r.Reloaded[k])
	}
	pog.SetStatus(nil)

	pog.SetStatus(pogText("Removing stale releases"))
	err = d.Prune(fmt.Sprintf("/opt/%s/releases", s.Spec.Application.Identifier), 5)
	if err != nil {
		return r, err
	}
	pog.Info("Removed stale releases")
	pog.SetStatus(nil)

	return r, nil
}

func uploadTarball(c *ssh.Client, dst string, tar Tarball) error {
	f, err := os.Open(tar.Path)
	if err != nil {
		return err
	}
	defer f.Close()
	s, err := f.Stat()
	if err != nil {
		return err
	}
	chstatus := make(chan ssh.PushStatus)
	go func() {
		for status := range chstatus {
			pog.SetStatus(pogUploadTarball(status.N, status.Size))
		}
	}()
	return c.Push(dst, s.Mode(), s.Size(), f, chstatus)
}
