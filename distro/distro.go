package distro

import (
	"errors"

	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
)

var ErrBadDistribution = errors.New("distro: unsupported distribution")

type Distro interface {
	InstallDocker() error

	MkdirAll(name string) error
	Remove(name string) error
	Symlink(oldname, newname string) error

	ExtractTar(name, dir string) error

	Install(app spec.Application, proc spec.Program) error
	Enable(proc spec.Program) error
	Disable(proc spec.Program) error
	Start(proc spec.Program) error
	Stop(proc spec.Program) error
	Restart(proc spec.Program) error

	Detect() (bool, error)
}

func New(c *ssh.Client) (Distro, error) {
	for _, fn := range DistroFuncs {
		d := fn(c)
		ok, err := d.Detect()
		if err != nil {
			return nil, err
		}

		if ok {
			return d, nil
		}
	}
	return nil, ErrBadDistribution
}

type DistroFunc func(c *ssh.Client) Distro

var DistroFuncs = []DistroFunc{}
