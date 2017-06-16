package ubuntu

import (
	"fmt"

	"github.com/FurqanSoftware/bullet/distro"
	"github.com/FurqanSoftware/bullet/distro/docker"
	"github.com/FurqanSoftware/bullet/distro/systemd"
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
		"echo 1 > /proc/sys/net/ipv6/conf/all/disable_ipv6",
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

func (u *Ubuntu) Install(app spec.Application, proc spec.Program) error {
	err := docker.Install(u.Client, app, proc, docker.InstallOptions{
		DockerPath: "/usr/bin/docker",
	})
	if err != nil {
		return err
	}

	return systemd.Install(u.Client, app, proc, systemd.InstallOptions{
		DockerPath: "/usr/bin/docker",
	})
}

func (u *Ubuntu) Enable(proc spec.Program) error {
	return systemd.Enable(u.Client, proc)
}

func (u *Ubuntu) Disable(proc spec.Program) error {
	return systemd.Disable(u.Client, proc)
}

func (u *Ubuntu) Start(proc spec.Program) error {
	return systemd.Start(u.Client, proc)
}

func (u *Ubuntu) Stop(proc spec.Program) error {
	return systemd.Stop(u.Client, proc)
}

func (u *Ubuntu) Restart(proc spec.Program) error {
	return systemd.Restart(u.Client, proc)
}

func (u *Ubuntu) Detect() (bool, error) {
	return true, nil
}

func init() {
	distro.DistroFuncs = append(distro.DistroFuncs, New)
}
