package distro

import (
	"errors"
	"text/tabwriter"

	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
)

var ErrBadDistribution = errors.New("distro: unsupported distribution")

type Distro interface {
	InstallDocker() error

	MkdirAll(name string) error
	Remove(name string) error
	Symlink(oldname, newname string) error
	Touch(name string) error
	Prune(name string, n int) error

	ExtractTar(name, dir string) error

	Build(app spec.Application, prog spec.Program) error
	Restart(app spec.Application, prog spec.Program, no int) error
	RestartAll(app spec.Application, prog spec.Program) error
	Status(app spec.Application, prog spec.Program, tw *tabwriter.Writer) error
	Scale(app spec.Application, prog spec.Program, n int) error

	CronEnable(app spec.Application, job spec.Job) error
	CronDisable(app spec.Application, job spec.Job) error
	CronStatus(app spec.Application, job spec.Job, tw *tabwriter.Writer) error

	Run(app spec.Application, prog spec.Program) error

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
