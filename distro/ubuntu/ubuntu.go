package ubuntu

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strconv"
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
	version, err := u.Client.Output("docker version")
	if err == nil && len(version) > 0 {
		return nil
	}

	cmds := []string{
		"apt-get update",
		"apt-get install -y apt-transport-https ca-certificates curl software-properties-common",
		"curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -",
		`add-apt-repository -y "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"`,
		"apt-get update",
		"apt-get install -y docker-ce",
	}
	for _, cmd := range cmds {
		err := u.Client.Run(cmd, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Ubuntu) MkdirAll(name string) error {
	return u.Client.Run(fmt.Sprintf("mkdir -p %s", name), false)
}

func (u *Ubuntu) Remove(name string) error {
	return u.Client.Run(fmt.Sprintf("rm %s", name), false)
}

func (u *Ubuntu) Symlink(oldname, newname string) error {
	return u.Client.Run(fmt.Sprintf("ln -sfn %s %s", oldname, newname), false)
}

func (u *Ubuntu) Touch(name string) error {
	return u.Client.Run(fmt.Sprintf("touch %s", name), false)
}

func (u *Ubuntu) Prune(name string, n int) error {
	return u.Client.Run(fmt.Sprintf("cd %s; ls -F . | head -n -%d | xargs -r rm -r", name, n), false)
}

func (u *Ubuntu) ReadFile(name string) ([]byte, error) {
	return u.Client.Output(fmt.Sprintf("cat %s", name))
}

func (u *Ubuntu) WriteFile(name string, data []byte) error {
	return u.Client.Run(fmt.Sprintf("echo %s | base64 -d | tee %s", base64.StdEncoding.EncodeToString(data), name), false)
}

func (u *Ubuntu) ExtractTar(name, dir string) error {
	cmds := []string{
		fmt.Sprintf("mkdir %s", dir),
		fmt.Sprintf("tar -xf %s -C %s", name, dir),
	}
	for _, cmd := range cmds {
		err := u.Client.Run(cmd, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Ubuntu) UpdateCurrent(app spec.Application, relDir string) error {
	curDir := fmt.Sprintf("/opt/%s/current", app.Identifier)
	switch app.Deploy.Current {
	case "replace":
		cmds := []string{
			fmt.Sprintf("mkdir -p %s", curDir),
			fmt.Sprintf("find %s -mindepth 1 -delete", curDir),
			fmt.Sprintf("cp -a %s/. %s/", relDir, curDir),
		}
		for _, cmd := range cmds {
			err := u.Client.Run(cmd, false)
			if err != nil {
				return err
			}
		}
		return nil

	case "symlink", "":
		return u.Symlink(relDir, curDir)
	}
	panic("unreachable")
}

func (u *Ubuntu) Build(app spec.Application, prog spec.Program) (bool, error) {
	return docker.BuildImage(u.Client, app, prog, docker.BuildImageOptions{
		DockerPath: dockerPath,
	})
}

func (u *Ubuntu) Restart(app spec.Application, prog spec.Program, no int) error {
	return docker.RestartContainer(u.Client, app, prog, no, docker.RestartContainerOptions{
		DockerPath: dockerPath,
	})
}

func (u *Ubuntu) Status(app spec.Application, prog spec.Program) ([]distro.Status, error) {
	status := []distro.Status{}

	conts, err := docker.ListContainers(u.Client, app, prog, docker.ListContainersOptions{
		DockerPath: dockerPath,
	})
	if err != nil {
		return nil, err
	}

	for _, cont := range conts {
		status = append(status, distro.Status{
			Program: prog,
			No:      cont.No,
			Up:      strings.HasPrefix(cont.Status, "Up"),
			Healthy: strings.Contains(cont.Status, "(healthy)"),
		})
	}

	return status, nil
}

func (u *Ubuntu) Scale(app spec.Application, prog spec.Program, n int) (int, int, error) {
	return docker.ScaleContainer(u.Client, app, prog, n, docker.ScaleContainerOptions{
		DockerPath: dockerPath,
	})
}

func (u *Ubuntu) Log(app spec.Application, prog spec.Program, no int) error {
	return docker.LogContainer(u.Client, app, prog, no, docker.LogContainerOptions{
		DockerPath: dockerPath,
	})
}

func (u *Ubuntu) Signal(app spec.Application, prog spec.Program, no int, signal string) error {
	return docker.SignalContainer(u.Client, app, prog, no, signal, docker.SignalContainerOptions{
		DockerPath: dockerPath,
	})
}

func (u *Ubuntu) Reload(app spec.Application, prog spec.Program, no int, rebuilt bool) error {
	if prog.Reload.PreCommand != "" {
		err := docker.ExecuteContainer(u.Client, app, prog, no, prog.Reload.PreCommand, docker.ExecuteContainerOptions{
			DockerPath: dockerPath,
		})
		if err != nil {
			return err
		}
	}

	method := prog.Reload.Method
	if rebuilt {
		method = ""
	}
	switch method {
	case "signal":
		return u.Signal(app, prog, no, prog.Reload.Signal)

	case "command":
		return docker.ExecuteContainer(u.Client, app, prog, no, prog.Reload.Command, docker.ExecuteContainerOptions{
			DockerPath: dockerPath,
		})

	case "restart", "":
		return u.Restart(app, prog, no)
	}
	return nil
}

func (u *Ubuntu) CronEnable(app spec.Application, job spec.Job) error {
	appdir := fmt.Sprintf("/opt/%s", app.Identifier)
	name := app.Identifier + "_cron_" + job.Key

	servicename := "bullet_" + app.Identifier + "_" + job.Key + ".service"
	timername := "bullet_" + app.Identifier + "_" + job.Key + ".timer"

	servicepreamble := []string{}
	servicepostamble := []string{}
	if job.Healthcheck.URL != "" {
		servicepreamble = append(servicepreamble, "ExecStartPre=-curl -sS -m 10 --retry 5 "+job.Healthcheck.URL+"/start")
		servicepostamble = append(servicepostamble, "ExecStopPost=-curl -sS -m 10 --retry 5 "+job.Healthcheck.URL+"/${EXIT_STATUS}")
	}

	service := `[Unit]
Description=Bullet task ` + app.Identifier + `_` + job.Key + `
Wants=` + timername + `

[Service]
Type=oneshot
EnvironmentFile=` + fmt.Sprintf("%s/env", appdir) + `
` + strings.Join(servicepreamble, "\n") + `
ExecStart=` + fmt.Sprintf("%s run --rm --env-file %s/env --name %s -v %s/current:/%s -w /%s %s %s", dockerPath, appdir, name, appdir, app.Identifier, app.Identifier, app.Identifier+"_shell", job.Command) + `
` + strings.Join(servicepostamble, "\n") + `

[Install]
WantedBy=multi-user.target`
	err := u.Client.Push(fmt.Sprintf("/etc/systemd/system/%s", servicename), 0644, int64(len(service)), strings.NewReader(service), nil)
	if err != nil {
		return err
	}

	timerpreamble := []string{}
	timerpostamble := []string{}
	if job.Jitter != "" {
		timerpostamble = append(timerpostamble, "RandomizedDelaySec="+job.Jitter)
		timerpostamble = append(timerpostamble, "FixedRandomDelay=true")
	}

	timer := `[Unit]
Description=Bullet task ` + app.Identifier + `_` + job.Key + `
Requires=` + servicename + `

[Timer]
` + strings.Join(timerpreamble, "\n") + `
Unit=` + servicename + `
OnCalendar=` + job.Schedule + `
` + strings.Join(timerpostamble, "\n") + `

[Install]
WantedBy=timers.target`
	err = u.Client.Push(fmt.Sprintf("/etc/systemd/system/%s", timername), 0644, int64(len(timer)), strings.NewReader(timer), nil)
	if err != nil {
		return err
	}

	return u.Client.Run(fmt.Sprintf("systemctl daemon-reload && systemctl enable --now %s", timername), false)
}

func (u *Ubuntu) CronDisable(app spec.Application, job spec.Job) error {
	servicename := "bullet_" + app.Identifier + "_" + job.Key + ".service"
	timername := "bullet_" + app.Identifier + "_" + job.Key + ".timer"

	err := u.Client.Run(fmt.Sprintf("[ ! -e %s ] || systemctl disable --now %s", timername, timername), false)
	if err != nil {
		return err
	}

	for _, name := range []string{
		fmt.Sprintf("/etc/systemd/system/%s", servicename),
		fmt.Sprintf("/etc/systemd/system/%s", timername),
	} {
		err = u.Client.Run(fmt.Sprintf("[ ! -e %s ] || rm %s", name, name), false)
		if err != nil {
			return err
		}
	}

	return u.Client.Run("systemctl daemon-reload", false)
}

func (u *Ubuntu) CronStatus(app spec.Application, job spec.Job, tw *tabwriter.Writer) error {
	timername := "bullet_" + app.Identifier + "_" + job.Key + ".timer"
	status, err := u.Client.Output(fmt.Sprintf("[ ! -e /etc/systemd/system/%s ] || systemctl status %s", timername, timername))
	if err != nil {
		return err
	}

	fmt.Fprintf(tw, "%s:\t", job.Key)
	active := false
	for _, l := range bytes.Split(status, []byte("\n")) {
		l = bytes.TrimSpace(l)
		if bytes.HasPrefix(l, []byte("Active:")) {
			fmt.Fprintf(tw, "%s", bytes.TrimPrefix(l, []byte("Active: ")))
			active = true
		}
		if bytes.HasPrefix(l, []byte("Trigger:")) {
			fmt.Fprintf(tw, " (trigger: %s)", bytes.TrimPrefix(l, []byte("Trigger: ")))
		}
	}
	if !active {
		fmt.Fprintf(tw, "disabled")
	}
	fmt.Fprint(tw, "\n")

	return nil
}

func (u *Ubuntu) Run(app spec.Application, prog spec.Program) error {
	return docker.RunContainer(u.Client, app, prog, docker.RunContainerOptions{
		DockerPath: dockerPath,
	})
}

func (u *Ubuntu) Forward(app spec.Application, port string) error {
	var (
		local  int
		remote int
	)
	parts := strings.SplitN(port, ":", 2)
	var err error
	if len(parts) == 2 {
		local, err = strconv.Atoi(parts[0])
		if err != nil {
			return err
		}
		remote, err = strconv.Atoi(parts[1])
		if err != nil {
			return err
		}
	} else {
		remote, err = strconv.Atoi(parts[0])
		if err != nil {
			return err
		}
		local = remote
	}

	return u.Client.Forward(local, remote)
}

func (u *Ubuntu) Df() error {
	return u.Client.Run("df", true)
}

func (u *Ubuntu) Top() error {
	return u.Client.RunPTY("top")
}

func (u *Ubuntu) Detect() (bool, error) {
	return true, nil
}

func init() {
	distro.DistroFuncs = append(distro.DistroFuncs, New)
}
