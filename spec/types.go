package spec

import "time"

type Spec struct {
	Application Application
}

type Application struct {
	Name       string
	Identifier string
	Deploy     Deploy

	Programs    map[string]Program
	ProgramKeys []string

	Cron Cron
}

type Build struct {
	Image  string
	Script []string
}

type Package struct {
	Contents []string
}

type Deploy struct {
}

type Program struct {
	Key         string `yaml:"-"`
	Name        string
	Command     string
	Container   Container
	Ports       []string
	Healthcheck *ProgramHealthcheck

	Unsafe Unsafe
}

type Container struct {
	Dockerfile string
	Image      string
}

type ProgramHealthcheck struct {
	Command     string
	Interval    time.Duration
	Timeout     time.Duration
	Retries     int
	StartPeriod time.Duration
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
