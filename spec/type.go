package spec

import "time"

type Application struct {
	Name       string
	Identifier string
	Deploy     Deploy

	Programs    map[string]Program
	ProgramKeys []string

	Cron Cron
}

func (a *Application) ApplyScope(scope *Scope) error {
	for k, prog := range a.Programs {
		err := prog.ApplyScope(scope)
		if err != nil {
			return err
		}
		a.Programs[k] = prog
	}
	return nil
}

type Build struct {
	Image  string
	Script []string
}

type Package struct {
	Contents []string
}

type Deploy struct {
	Current string
}

type Program struct {
	Key         string `yaml:"-"`
	Name        string
	Command     string
	User        string
	Container   Container
	Ports       []string
	Volumes     []string
	Healthcheck *ProgramHealthcheck
	Scales      []Scale
	Reload      Reload

	Unsafe Unsafe
}

func (p *Program) ApplyScope(scope *Scope) error {
	var err error
	p.Command, err = scope.Expand(p.Command)
	if err != nil {
		return err
	}
	return nil
}

type Container struct {
	Dockerfile     string
	Image          string
	Entrypoint     *string
	WorkingDir     *string
	ApplicationDir *string
}

type ProgramHealthcheck struct {
	Command     string
	Interval    time.Duration
	Timeout     time.Duration
	Retries     int
	StartPeriod time.Duration
}

type Scale struct {
	If string
	N  string
}

type Reload struct {
	Method     string
	Signal     string
	Command    string
	PreCommand string
}

type Unsafe struct {
	NetworkHost bool
}

type Cron struct {
	Jobs []Job
}

func (c Cron) Job(k string) Job {
	for _, j := range c.Jobs {
		if j.Key == k {
			return j
		}
	}
	return Job{}
}

type Job struct {
	Key         string
	Command     string
	Schedule    string
	Jitter      string
	Healthcheck JobHealthcheck
}

type JobHealthcheck struct {
	URL string
}
