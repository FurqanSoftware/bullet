package ubuntu

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

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
	} else {
		fmt.Fprintf(tw, "%s:\tdisabled\n", prog.Key)
	}
	return nil
}

func (u *Ubuntu) Scale(app spec.Application, prog spec.Program, n int) error {
	return docker.ScaleContainer(u.Client, app, prog, n, docker.ScaleContainerOptions{
		DockerPath: dockerPath,
	})
}

func (u *Ubuntu) CronEnable(app spec.Application, job spec.Job) error {
	appdir := fmt.Sprintf("/opt/%s", app.Identifier)
	name := app.Identifier + "_cron_" + job.Key

	servicename := "bullet_" + app.Identifier + "_" + job.Key + ".service"
	timername := "bullet_" + app.Identifier + "_" + job.Key + ".timer"

	service := `[Unit]
Description=Bullet task ` + app.Identifier + `_` + job.Key + `
Wants=` + timername + `

[Service]
Type=oneshot
ExecStart=` + fmt.Sprintf("%s run --rm --env-file %s/env --name %s -v %s/current:/%s -w /%s %s %s", dockerPath, appdir, name, appdir, app.Identifier, app.Identifier, app.Identifier+"_shell", job.Command) + `

[Install]
WantedBy=multi-user.target`
	err := u.Client.Push(fmt.Sprintf("/etc/systemd/system/%s", servicename), 0644, int64(len(service)), strings.NewReader(service))
	if err != nil {
		return err
	}

	timer := `[Unit]
Description=Bullet task ` + app.Identifier + `_` + job.Key + `
Requires=` + servicename + `

[Timer]
Unit=` + servicename + `
OnCalendar=` + job.Schedule + `

[Install]
WantedBy=timers.target`
	err = u.Client.Push(fmt.Sprintf("/etc/systemd/system/%s", timername), 0644, int64(len(timer)), strings.NewReader(timer))
	if err != nil {
		return err
	}

	return u.Client.Run(fmt.Sprintf("systemctl enable --now %s", timername))
}

func (u *Ubuntu) CronDisable(app spec.Application, job spec.Job) error {
	timername := "bullet_" + app.Identifier + "_" + job.Key + ".timer"
	return u.Client.Run(fmt.Sprintf("systemctl disable --now %s", timername))
}

func (u *Ubuntu) CronStatus(app spec.Application, job spec.Job, tw *tabwriter.Writer) error {
	timername := "bullet_" + app.Identifier + "_" + job.Key + ".timer"
	status, err := u.Client.Output(fmt.Sprintf("systemctl status %s", timername))
	if err != nil {
		return err
	}

	fmt.Fprintf(tw, "%s:\t", job.Key)
	for _, l := range bytes.Split(status, []byte("\n")) {
		l = bytes.TrimSpace(l)
		if bytes.HasPrefix(l, []byte("Active:")) {
			fmt.Fprintf(tw, "%s", bytes.TrimPrefix(l, []byte("Active: ")))
		}
		if bytes.HasPrefix(l, []byte("Trigger:")) {
			fmt.Fprintf(tw, " (trigger: %s)", bytes.TrimPrefix(l, []byte("Trigger: ")))
		}
	}
	fmt.Fprint(tw, "\n")

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
