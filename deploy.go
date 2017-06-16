package bullet

import (
	"fmt"
	"log"
	"os"

	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
)

type DeployOptions struct {
	SkipBuild bool
}

func Deploy(nodes []Node, spec *spec.Spec, options DeployOptions) error {
	if options.SkipBuild {
		log.Print("Skipping build")

	} else {
		log.Print("Building")
		err := Build(spec)
		if err != nil {
			return err
		}
	}

	log.Print("Packaging")
	rel, err := Package(spec)
	if err != nil {
		return err
	}

	for _, n := range nodes {
		log.Printf("Connecting to %s", n.Addr())
		c, err := ssh.Dial(n.Addr())
		if err != nil {
			return err
		}

		log.Printf("Uploading tarball %s", n.Addr())
		tarPath := fmt.Sprintf("/tmp/%s.tar", rel.Name)
		err = uploadTarball(c, tarPath, rel.Tarball)
		if err != nil {
			return err
		}

		d, err := distro.New(c)
		if err != nil {
			return err
		}
		err = d.ExtractTar(tarPath, fmt.Sprintf("/opt/bullet/%s/releases/%s", spec.Application.Identifier, rel.Name))
		if err != nil {
			return err
		}
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
	return c.Scp(dst, s.Mode(), s.Size(), f)
}
