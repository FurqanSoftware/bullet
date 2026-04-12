package spec

import (
	"time"
)

// Application defines the application and its programs, cron jobs, and deploy settings.
type Application struct {
	Name       string
	Identifier string
	Deploy     Deploy

	Programs    map[string]Program
	ProgramKeys []string

	Cron Cron
}

// ExpandVars expands template variables in all program commands.
func (a *Application) ExpandVars(vars Vars) error {
	for k, prog := range a.Programs {
		err := prog.ExpandVars(vars)
		if err != nil {
			return err
		}
		a.Programs[k] = prog
	}
	return nil
}

// Build defines the build image and script for the application.
type Build struct {
	Image  string
	Script []string
}

// Package defines the files included in a release tarball.
type Package struct {
	Contents []string
}

// Deploy configures the deployment behavior (e.g. symlink vs replace).
type Deploy struct {
	Current string
}

// Program defines a long-running service within the application.
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

// ExpandVars expands template variables in the program's command.
func (p *Program) ExpandVars(vars Vars) error {
	var err error
	p.Command, err = vars.Expand(p.Command)
	if err != nil {
		return err
	}
	return nil
}

// Container configures the Docker image or build for a program.
type Container struct {
	Dockerfile     string
	Image          string
	Entrypoint     *string
	WorkingDir     *string
	ApplicationDir *string
}

// ProgramHealthcheck configures Docker's built-in health check for a program.
type ProgramHealthcheck struct {
	Command     string
	Interval    time.Duration
	Timeout     time.Duration
	Retries     int
	StartPeriod time.Duration
}

// Scale defines a conditional scaling rule for a program.
type Scale struct {
	If string
	N  string
}

// Reload configures how a program's containers are reloaded on deploy.
type Reload struct {
	Method     string
	Signal     string
	Command    string
	PreCommand string
}

// Unsafe holds flags for unsafe container options.
type Unsafe struct {
	NetworkHost bool
	Ulimits     []string
}

// Cron holds the scheduled jobs for the application.
type Cron struct {
	Jobs []Job
}

// Job returns the job with the given key, or an empty Job if not found.
func (c Cron) Job(k string) Job {
	for _, j := range c.Jobs {
		if j.Key == k {
			return j
		}
	}
	return Job{}
}

// Job defines a scheduled cron job.
type Job struct {
	Key         string
	Command     string
	Schedule    string
	Jitter      string
	Healthcheck JobHealthcheck
}

// JobHealthcheck configures health check pings for a cron job.
type JobHealthcheck struct {
	URL string
}
