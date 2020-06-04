package ubuntu

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/tabwriter"

	"github.com/FurqanSoftware/bullet/distro"
	"github.com/FurqanSoftware/bullet/docker"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
	cryptossh "golang.org/x/crypto/ssh"
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

func (u *Ubuntu) Prune(name string, n int) error {
	return u.Client.Run(fmt.Sprintf("cd %s; ls -F . | head -n -%d | xargs -r rm -r", name, n))
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
		DockerPath: dockerPath,
	})
}

func (u *Ubuntu) Restart(app spec.Application, prog spec.Program, no int) error {
	return docker.RestartContainer(u.Client, app, prog, no, docker.RestartContainerOptions{
		DockerPath: dockerPath,
	})
}

func (u *Ubuntu) RestartAll(app spec.Application, prog spec.Program) error {
	conts, err := docker.ListContainers(u.Client, app, prog, docker.ListContainersOptions{
		DockerPath: dockerPath,
	})
	if err != nil {
		return err
	}

	for _, cont := range conts {
		err = docker.RestartContainer(u.Client, app, prog, cont.No, docker.RestartContainerOptions{
			DockerPath: dockerPath,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Ubuntu) Status(app spec.Application, prog spec.Program, tw *tabwriter.Writer) error {
	conts, err := docker.ListContainers(u.Client, app, prog, docker.ListContainersOptions{
		DockerPath: dockerPath,
	})
	if err != nil {
		return err
	}

	if len(conts) > 0 {
		for _, cont := range conts {
			fmt.Fprintf(tw, "%s:\t%s\t(%s)\n", prog.Key, strings.ToLower(cont.Status), cont.ID)
		}
	}
	return nil
}

func (u *Ubuntu) Scale(app spec.Application, prog spec.Program, n int) error {
	return docker.ScaleContainer(u.Client, app, prog, n, docker.ScaleContainerOptions{
		DockerPath: dockerPath,
	})
}

func (u *Ubuntu) CronEnable(app spec.Application, job spec.Job) error {
	crontab, err := u.Client.Output("crontab -l")
	_, ok := err.(*cryptossh.ExitError)
	if ok {
		err = nil
	}
	if err != nil {
		return err
	}
	lines := []string{}
	re := regexp.MustCompile("# Bullet " + regexp.QuoteMeta(app.Identifier) + "_" + regexp.QuoteMeta(job.Key))
	for _, l := range bytes.Split(bytes.TrimSpace(crontab), []byte("\n")) {
		if !re.Match(l) {
			lines = append(lines, string(l))
		}
	}
	if len(lines) == 1 && lines[0] == "" {
		lines = lines[:0]
	}
	appdir := fmt.Sprintf("/opt/%s", app.Identifier)
	name := app.Identifier + "_cron_" + job.Key
	lines = append(lines, job.Schedule+" "+fmt.Sprintf("%s run --rm --env-file %s/env --name %s -v %s/current:/%s -w /%s %s %s", dockerPath, appdir, name, appdir, app.Identifier, app.Identifier, app.Identifier+"_shell", job.Command)+" && (mkdir -p "+appdir+"/logs && touch "+appdir+"/logs/cron.log && echo `date +\\%Y-\\%m-\\%d \\%H:\\%M:\\%S` 'Job "+job.Key+" succeeded' >> "+appdir+"/logs/cron.log) || (mkdir -p "+appdir+"/logs && touch "+appdir+"/logs/cron.log && echo `date +\\%Y-\\%m-\\%d \\%H:\\%M:\\%S` 'Job "+job.Key+" failed' >> "+appdir+"/logs/cron.log) # Bullet "+app.Identifier+"_"+job.Key)
	crontab = []byte(strings.Join(lines, "\n") + "\n")

	err = u.Client.Push("/tmp/crontab", 0600, int64(len(crontab)), bytes.NewReader(crontab))
	if err != nil {
		return err
	}

	return u.Client.Run("crontab /tmp/crontab")
}

func (u *Ubuntu) CronDisable(app spec.Application, job spec.Job) error {
	crontab, err := u.Client.Output("crontab -l")
	_, ok := err.(*cryptossh.ExitError)
	if ok {
		err = nil
	}
	if err != nil {
		return err
	}
	lines := []string{}
	match := false
	re := regexp.MustCompile("# Bullet " + regexp.QuoteMeta(app.Identifier) + "_" + regexp.QuoteMeta(job.Key))
	for _, l := range bytes.Split(bytes.TrimSpace(crontab), []byte("\n")) {
		if !re.Match(l) {
			lines = append(lines, string(l))
		} else {
			match = true
		}
	}
	if !match {
		return nil
	}
	crontab = []byte(strings.Join(lines, "\n") + "\n")

	err = u.Client.Push("/tmp/crontab", 0600, int64(len(crontab)), bytes.NewReader(crontab))
	if err != nil {
		return err
	}

	return u.Client.Run("crontab /tmp/crontab")
}

func (u *Ubuntu) CronStatus(app spec.Application, job spec.Job, tw *tabwriter.Writer) error {
	crontab, err := u.Client.Output("crontab -l")
	_, ok := err.(*cryptossh.ExitError)
	if ok {
		err = nil
	}
	if err != nil {
		return err
	}
	match := false
	re := regexp.MustCompile("# Bullet " + regexp.QuoteMeta(app.Identifier) + "_" + regexp.QuoteMeta(job.Key))
	for _, l := range bytes.Split(bytes.TrimSpace(crontab), []byte("\n")) {
		if re.Match(l) {
			match = true
		}
	}

	fmt.Fprintf(tw, "%s:\t", job.Key)
	if match {
		fmt.Fprintln(tw, "enabled")
	} else {
		fmt.Fprintln(tw, "disabled")
	}
	return nil
}

func (u *Ubuntu) Run(app spec.Application, prog spec.Program) error {
	return docker.RunContainer(u.Client, app, prog, docker.RunContainerOptions{
		DockerPath: dockerPath,
	})
}

func (u *Ubuntu) Detect() (bool, error) {
	return true, nil
}

func init() {
	distro.DistroFuncs = append(distro.DistroFuncs, New)
}
