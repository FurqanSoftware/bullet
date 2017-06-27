package ubuntu

import (
	"fmt"
	"strings"

	"github.com/FurqanSoftware/bullet/distro"
	"github.com/FurqanSoftware/bullet/docker"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
)

type Ubuntu struct {
	Client *ssh.Client
}

func New(c *ssh.Client) distro.Distro {
	return &Ubuntu{
		Client: c,
	}
}

func (u *Ubuntu) InstallDocker() error {
	cmds := []string{
		"apt-get update",
		"apt-get install -y apt-transport-https ca-certificates curl software-properties-common",
		"curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -",
		`add-apt-repository -y "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"`,
		"apt-get update",
		"apt-get install -y docker-ce",
	}
	for _, cmd := range cmds {
		err := u.Client.Run(cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Ubuntu) MkdirAll(name string) error {
	return u.Client.Run(fmt.Sprintf("mkdir -p %s", name))
}

func (u *Ubuntu) Remove(name string) error {
	return u.Client.Run(fmt.Sprintf("rm %s", name))
}

func (u *Ubuntu) Symlink(oldname, newname string) error {
	return u.Client.Run(fmt.Sprintf("ln -sfn %s %s", oldname, newname))
}

func (u *Ubuntu) Touch(name string) error {
	return u.Client.Run(fmt.Sprintf("touch %s", name))
}

func (u *Ubuntu) ExtractTar(name, dir string) error {
	cmds := []string{
		fmt.Sprintf("mkdir %s", dir),
		fmt.Sprintf("tar -xf %s -C %s", name, dir),
	}
	for _, cmd := range cmds {
		err := u.Client.Run(cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Ubuntu) Build(app spec.Application, prog spec.Program) error {
	return docker.BuildImage(u.Client, app, prog, docker.BuildImageOptions{
		DockerPath: "/usr/bin/docker",
	})
}

func (u *Ubuntu) Restart(app spec.Application, prog spec.Program, no int) error {
	return docker.RestartContainer(u.Client, app, prog, no, docker.RestartContainerOptions{
		DockerPath: "/usr/bin/docker",
	})
}

func (u *Ubuntu) RestartAll(app spec.Application, prog spec.Program) error {
	conts, err := docker.ListContainers(u.Client, app, prog, docker.ListContainersOptions{
		DockerPath: "/usr/bin/docker",
	})
	if err != nil {
		return err
	}

	for _, cont := range conts {
		err = docker.RestartContainer(u.Client, app, prog, cont.No, docker.RestartContainerOptions{
			DockerPath: "/usr/bin/docker",
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Ubuntu) Status(app spec.Application, prog spec.Program) error {
	conts, err := docker.ListContainers(u.Client, app, prog, docker.ListContainersOptions{
		DockerPath: "/usr/bin/docker",
	})
	if err != nil {
		return err
	}

	fmt.Println(prog.Name)
	for _, cont := range conts {
		fmt.Printf("%s: %s\n", cont.ID, strings.ToLower(cont.Status))
	}
	fmt.Println()
	return nil
}

func (u *Ubuntu) Scale(app spec.Application, prog spec.Program, n int) error {
	return docker.ScaleContainer(u.Client, app, prog, n, docker.ScaleContainerOptions{
		DockerPath: "/usr/bin/docker",
	})
}

func (u *Ubuntu) Detect() (bool, error) {
	return true, nil
}

func init() {
	distro.DistroFuncs = append(distro.DistroFuncs, New)
}
