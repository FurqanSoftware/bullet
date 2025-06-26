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
	ReadFile(name string) ([]byte, error)
	WriteFile(name string, data []byte) error

	ExtractTar(name, dir string) error

	UpdateCurrent(app spec.Application, relDir string) error

	Build(app spec.Application, prog spec.Program) (bool, error)
	Restart(app spec.Application, prog spec.Program, no int) error
	Status(app spec.Application, prog spec.Program) ([]Status, error)
	Scale(app spec.Application, prog spec.Program, n int) (int, int, error)
	Log(app spec.Application, prog spec.Program, no int) error
	Signal(app spec.Application, prog spec.Program, no int, signal string) error
	Reload(app spec.Application, prog spec.Program, no int, rebuilt bool) error

	CronEnable(app spec.Application, job spec.Job) error
	CronDisable(app spec.Application, job spec.Job) error
	CronStatus(app spec.Application, job spec.Job, tw *tabwriter.Writer) error

	Run(app spec.Application, prog spec.Program) error

	Forward(app spec.Application, port string) error

	Df() error
	Top() error

	Detect() (bool, error)
}

type Status struct {
	Program spec.Program
	No      int
	Up      bool
	Healthy bool
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
