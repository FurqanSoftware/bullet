package spec

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
	Key       string `yaml:"-"`
	Name      string
	Command   string
	Container Container
	Ports     []string
}

type Container struct {
	Dockerfile string
	Image      string
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
	Key      string
	Command  string
	Schedule string
}
