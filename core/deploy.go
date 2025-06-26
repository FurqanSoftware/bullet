package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
	"github.com/FurqanSoftware/pog"
	"github.com/dustin/go-humanize"
)

func Deploy(nodes []Node, spec *spec.Spec, rel *Release) error {
	pog.Infof("Deploying %s", filepath.Base(rel.Tarball.Path))
	pog.Infof("∟ Hash: %s", rel.Hash)
	pog.Infof("∟ Size: %s", humanize.Bytes(uint64(rel.Tarball.Size)))

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

		err = deployNode(n, c, d, spec, rel)
		if err != nil {
			return err
		}
	}
	return nil
}

func deployNode(n Node, c *ssh.Client, d distro.Distro, spec *spec.Spec, rel *Release) error {
	curHash, _ := d.ReadFile(fmt.Sprintf("/opt/%s/current.hash", spec.Application.Identifier))
	if rel.Hash == string(curHash) {
		pog.Info("Same as current hash. Skipping deploy.")
		return nil
	}

	tarPath := fmt.Sprintf("/tmp/%s-%s-%s.tar.gz", spec.Application.Identifier, rel.Time, rel.Hash)
	err := uploadTarball(c, tarPath, rel.Tarball)
	if err != nil {
		return err
	}
	pog.Info("Uploaded tarball")
	pog.SetStatus(nil)

	pog.SetStatus(pogText("Extracting tarball"))
	relDir := fmt.Sprintf("/opt/%s/releases/%s-%s", spec.Application.Identifier, rel.Time, rel.Hash)
	err = d.ExtractTar(tarPath, relDir)
	if err != nil {
		return err
	}
	pog.Info("Extracted tarball")
	pog.SetStatus(nil)

	pog.SetStatus(pogText("Removing tarball"))
	err = d.Remove(tarPath)
	if err != nil {
		return err
	}
	pog.Info("Removed tarball")
	pog.SetStatus(nil)

	pog.SetStatus(pogText("Updating current"))
	err = d.UpdateCurrent(spec.Application, relDir)
	if err != nil {
		return err
	}
	pog.Info("Updated current")
	pog.SetStatus(nil)

	err = d.WriteFile(fmt.Sprintf("/opt/%s/current.hash", spec.Application.Identifier), []byte(rel.Hash))
	if err != nil {
		return err
	}

	pog.SetStatus(pogText("Building images"))
	rebuilts := map[string]bool{}
	for _, p := range spec.Application.Programs {
		rebuilt, err := d.Build(spec.Application, p)
		if err != nil {
			return err
		}
		if rebuilt {
			rebuilts[p.Key] = true
		}
	}
	pog.Infof("Built %d image(s)", len(rebuilts))
	for _, k := range spec.Application.ProgramKeys {
		if !rebuilts[k] {
			continue
		}
		pog.Infof("∟ %s", k)
	}
	pog.SetStatus(nil)

	pog.SetStatus(pogText("Reloading containers"))
	nreload := map[string]int{}
	nreloadsum := 0
	for _, k := range spec.Application.ProgramKeys {
		p := spec.Application.Programs[k]
		status, err := d.Status(spec.Application, p)
		if err != nil {
			return err
		}
		for _, s := range status {
			if s.No == 0 {
				continue
			}
			pog.SetStatus(pogReloadingContainer(p, s.No))
			err = d.Reload(spec.Application, p, s.No, rebuilts[k])
			if err != nil {
				return err
			}
			nreload[k]++
			nreloadsum++
		}
	}
	pog.Infof("Reloaded %d container(s)", nreloadsum)
	for _, k := range spec.Application.ProgramKeys {
		if nreload[k] == 0 {
			continue
		}
		pog.Infof("∟ %s: %d", k, nreload[k])
	}
	pog.SetStatus(nil)

	pog.SetStatus(pogText("Removing stale releases"))
	err = d.Prune(fmt.Sprintf("/opt/%s/releases", spec.Application.Identifier), 5)
	if err != nil {
		return err
	}
	pog.Info("Removed stale releases")
	pog.SetStatus(nil)

	return nil
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
