package core

import (
	"fmt"
	"log"
	"os"

	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
)

func Deploy(nodes []Node, spec *spec.Spec, rel *Release) error {
	for _, n := range nodes {
		log.Printf("Connecting to %s", n.Label())
		c, err := ssh.Dial(n.Addr(), n.Identity)
		if err != nil {
			return err
		}

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
		log.Print("Same as current hash. Skipping deploy.")
		log.Printf(".. Hash: %s", rel.Hash)
		return nil
	}

	log.Print("Uploading tarball")
	log.Printf(".. Hash: %s", rel.Hash)
	tarPath := fmt.Sprintf("/tmp/%s-%s-%s.tar.gz", spec.Application.Identifier, rel.Time, rel.Hash)
	err := uploadTarball(c, tarPath, rel.Tarball)
	if err != nil {
		return err
	}

	log.Print("Extracting tarball")
	relDir := fmt.Sprintf("/opt/%s/releases/%s-%s", spec.Application.Identifier, rel.Time, rel.Hash)
	err = d.ExtractTar(tarPath, relDir)
	if err != nil {
		return err
	}
	log.Print("Removing tarball")
	err = d.Remove(tarPath)
	if err != nil {
		return err
	}

	log.Print("Updating current marker")
	err = d.Symlink(relDir, fmt.Sprintf("/opt/%s/current", spec.Application.Identifier))
	if err != nil {
		return err
	}
	err = d.WriteFile(fmt.Sprintf("/opt/%s/current.hash", spec.Application.Identifier), []byte(rel.Hash))
	if err != nil {
		return err
	}

	log.Print("Building images")
	rebuilts := map[string]bool{}
	for _, p := range spec.Application.Programs {
		rebuilt, err := d.Build(spec.Application, p)
		if err != nil {
			return err
		}
		rebuilts[p.Key] = rebuilt
	}

	log.Print("Reloading containers")
	for _, k := range spec.Application.ProgramKeys {
		p := spec.Application.Programs[k]
		err = d.ReloadAll(spec.Application, p, rebuilts[k])
		if err != nil {
			return err
		}
	}

	log.Print("Removing stale releases")
	err = d.Prune(fmt.Sprintf("/opt/%s/releases", spec.Application.Identifier), 5)
	if err != nil {
		return err
	}

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
	return c.Push(dst, s.Mode(), s.Size(), f)
}
