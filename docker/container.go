package docker

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
)

type Container struct {
	Application spec.Application
	Program     spec.Program
	No          int

	ID     string
	Image  string
	Status string
}

type ListContainersOptions struct {
	DockerPath string
}

func ListContainers(c *ssh.Client, app spec.Application, prog spec.Program, options ListContainersOptions) ([]Container, error) {
	reName, err := regexp.Compile(fmt.Sprintf(`^%s_%s_(\d+)$`, regexp.QuoteMeta(app.Identifier), regexp.QuoteMeta(prog.Key)))
	if err != nil {
		return nil, err
	}

	conts := []Container{}
	b, err := c.Output(fmt.Sprintf("%s ps -a --format '{{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}'", options.DockerPath))
	if err != nil {
		return nil, err
	}
	s := string(b)

	lines := strings.Split(s, "\n")
	for _, line := range lines {
		parts := strings.Split(line, "\t")
		if len(parts) != 4 {
			continue
		}
		m := reName.FindStringSubmatch(parts[1])
		if m == nil {
			continue
		}
		no, err := strconv.Atoi(m[1])
		if err != nil {
			continue
		}
		conts = append(conts, Container{
			Application: app,
			Program:     prog,
			No:          no,
			ID:          parts[0],
			Image:       parts[2],
			Status:      parts[3],
		})
	}
	return conts, nil
}

type RestartContainerOptions struct {
	DockerPath string
}

func RestartContainer(c *ssh.Client, app spec.Application, prog spec.Program, no int, options RestartContainerOptions) error {
	name := fmt.Sprintf("%s_%s_%d", app.Identifier, prog.Key, no)
	image := ""
	if prog.Container.Dockerfile != "" {
		image = fmt.Sprintf("%s_%s", app.Identifier, prog.Key)
	} else {
		image = prog.Container.Image
	}

	err := deleteContainer(c, app, prog, options.DockerPath, name)
	if err != nil {
		return err
	}

	return createContainer(c, app, prog, options.DockerPath, image, name, no)
}

type ScaleContainerOptions struct {
	DockerPath string
}

func ScaleContainer(c *ssh.Client, app spec.Application, prog spec.Program, n int, options ScaleContainerOptions) (up, down int, err error) {
	conts, err := ListContainers(c, app, prog, ListContainersOptions(options))
	if err != nil {
		return
	}
	lastNo := 0
	for _, cont := range conts {
		no := cont.No
		if no > lastNo {
			lastNo = no
		}
	}

	d := len(conts) - n
	for ; d > 0; d-- {
		name := fmt.Sprintf("%s_%s_%d", app.Identifier, prog.Key, lastNo-d+1)
		err = deleteContainer(c, app, prog, options.DockerPath, name)
		if err != nil {
			return
		}
		down++
	}
	for ; d < 0; d++ {
		name := fmt.Sprintf("%s_%s_%d", app.Identifier, prog.Key, lastNo-d)
		image := fmt.Sprintf("%s_%s", app.Identifier, prog.Key)
		err = createContainer(c, app, prog, options.DockerPath, image, name, lastNo-d)
		if err != nil {
			return
		}
		up++
	}

	return
}

type LogContainerOptions struct {
	DockerPath string
}

func LogContainer(c *ssh.Client, app spec.Application, prog spec.Program, no int, options LogContainerOptions) error {
	name := fmt.Sprintf("%s_%s_%d", app.Identifier, prog.Key, no)
	return logContainer(c, app, prog, options.DockerPath, name, no)
}

type RunContainerOptions struct {
	DockerPath string
}

func RunContainer(c *ssh.Client, app spec.Application, prog spec.Program, options RunContainerOptions) error {
	name := fmt.Sprintf("%s_%s_0", app.Identifier, prog.Key)
	image := fmt.Sprintf("%s_%s", app.Identifier, prog.Key)

	conts, err := ListContainers(c, app, prog, ListContainersOptions(options))
	if err != nil {
		return err
	}
	if len(conts) > 0 {
		err := deleteContainer(c, app, prog, options.DockerPath, name)
		if err != nil {
			return err
		}
	}

	err = createAttachContainer(c, app, prog, options.DockerPath, image, name)
	if err != nil {
		return err
	}

	conts, err = ListContainers(c, app, prog, ListContainersOptions(options))
	if err != nil {
		return err
	}
	if len(conts) > 0 {
		err := deleteContainer(c, app, prog, options.DockerPath, name)
		if err != nil {
			return err
		}
	}
	return nil
}

type SignalContainerOptions struct {
	DockerPath string
}

func SignalContainer(c *ssh.Client, app spec.Application, prog spec.Program, no int, signal string, options SignalContainerOptions) error {
	name := fmt.Sprintf("%s_%s_%d", app.Identifier, prog.Key, no)
	return signalContainer(c, app, prog, options.DockerPath, name, signal)
}

type ExecuteContainerOptions struct {
	DockerPath string
}

func ExecuteContainer(c *ssh.Client, app spec.Application, prog spec.Program, no int, command string, options ExecuteContainerOptions) error {
	name := fmt.Sprintf("%s_%s_%d", app.Identifier, prog.Key, no)
	return executeContainer(c, app, prog, options.DockerPath, name, command)
}

func createContainer(c *ssh.Client, app spec.Application, prog spec.Program, dockerPath, image, name string, no int) error {
	appDir := fmt.Sprintf("/opt/%s", app.Identifier)

	cmd := []string{
		dockerPath,
		"run",
		"-d",
		"-e", strconv.Quote("BULLET_APPLICATION_NAME=" + app.Name),
		"-e", strconv.Quote("BULLET_APPLICATION_ID=" + app.Identifier),
		"-e", strconv.Quote("BULLET_PROGRAM_KEY=" + prog.Key),
		"-e", strconv.Quote("BULLET_PROGRAM_NAME=" + prog.Name),
		"-e", strconv.Quote("BULLET_INSTANCE_ID=" + name),
		"--env-file", appDir + "/env",
		"--name", name,
	}
	for _, p := range prog.Ports {
		m := strings.SplitN(p, ":", 2)
		if len(m) == 1 {
			m = append(m, m[0])
		}
		h, err := strconv.Atoi(m[0])
		if err != nil {
			return err
		}
		cmd = append(cmd, "-p", fmt.Sprintf("%d:%s", h+no-1, m[1]))
	}
	if prog.User != "" {
		cmd = append(cmd, "--user", prog.User)
	}
	for _, v := range prog.Volumes {
		cmd = append(cmd, "-v", v)
	}
	if prog.Healthcheck != nil {
		cmd = append(
			cmd,
			"--health-cmd", prog.Healthcheck.Command,
			"--health-interval", prog.Healthcheck.Interval.String(),
			"--health-timeout", prog.Healthcheck.Timeout.String(),
			"--health-retries", strconv.Itoa(prog.Healthcheck.Retries),
			"--health-start-period", prog.Healthcheck.StartPeriod.String(),
		)
	}
	if prog.Unsafe.NetworkHost {
		cmd = append(cmd, "--network=host")
	}

	cmd = append(
		cmd,
		"--log-driver",
		"json-file",
		"--log-opt",
		"max-size=3g",
		"--log-opt",
		`tag="{{.Name}}"`,
		"--restart", "always",
	)

	if prog.Container.ApplicationDir != nil {
		cmd = append(cmd, "-v", strconv.Quote(appDir+"/current:"+*prog.Container.ApplicationDir))
	} else {
		cmd = append(cmd, "-v", appDir+"/current:/"+app.Identifier)
	}

	if prog.Container.WorkingDir != nil {
		cmd = append(cmd, "-w", strconv.Quote(*prog.Container.WorkingDir))
	} else if prog.Container.ApplicationDir != nil {
		cmd = append(cmd, "-w", strconv.Quote(*prog.Container.ApplicationDir))
	} else {
		cmd = append(cmd, "-w", "/"+app.Identifier)
	}

	if prog.Container.Entrypoint != nil {
		cmd = append(
			cmd,
			"--entrypoint", strconv.Quote(*prog.Container.Entrypoint),
		)
	}

	cmd = append(
		cmd,
		image,
		prog.Command,
	)

	return c.Run(strings.Join(cmd, " "), false)
}

func createAttachContainer(c *ssh.Client, app spec.Application, prog spec.Program, dockerPath, image, name string) error {
	appDir := fmt.Sprintf("/opt/%s", app.Identifier)

	cmd := []string{
		dockerPath,
		"run",
		"-ti",
		"-e", strconv.Quote("BULLET_APPLICATION_NAME=" + app.Name),
		"-e", strconv.Quote("BULLET_APPLICATION_ID=" + app.Identifier),
		"-e", strconv.Quote("BULLET_PROGRAM_KEY=" + prog.Key),
		"-e", strconv.Quote("BULLET_PROGRAM_NAME=" + prog.Name),
		"-e", strconv.Quote("BULLET_INSTANCE_ID=" + name),
		"--env-file", appDir + "/env",
		"--name", name,
	}
	if prog.User != "" {
		cmd = append(cmd, "--user", prog.User)
	}
	for _, v := range prog.Volumes {
		cmd = append(cmd, "-v", v)
	}
	if prog.Unsafe.NetworkHost {
		cmd = append(cmd, "--network=host")
	}

	if prog.Container.ApplicationDir != nil {
		cmd = append(cmd, "-v", strconv.Quote(appDir+"/current:"+*prog.Container.ApplicationDir))
	} else {
		cmd = append(cmd, "-v", appDir+"/current:/"+app.Identifier)
	}

	if prog.Container.WorkingDir != nil {
		cmd = append(cmd, "-w", strconv.Quote(*prog.Container.WorkingDir))
	} else if prog.Container.ApplicationDir != nil {
		cmd = append(cmd, "-w", strconv.Quote(*prog.Container.ApplicationDir))
	} else {
		cmd = append(cmd, "-w", "/"+app.Identifier)
	}

	cmd = append(
		cmd,
		`--entrypoint=""`,
		image,
		prog.Command,
	)

	return c.RunPTY(strings.Join(cmd, " "))
}

func deleteContainer(c *ssh.Client, app spec.Application, prog spec.Program, dockerPath, name string) error {
	cmds := []string{
		fmt.Sprintf("%s stop -t 2 %s > /dev/null 2>&1 || true", dockerPath, name),
		fmt.Sprintf("%s rm %s > /dev/null 2>&1 || true", dockerPath, name),
	}
	for _, cmd := range cmds {
		err := c.Run(cmd, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func logContainer(c *ssh.Client, app spec.Application, prog spec.Program, dockerPath, name string, no int) error {
	cmd := []string{
		dockerPath,
		"logs",
		"--tail 10",
		"--follow",
		name,
	}
	return c.Run(strings.Join(cmd, " "), true)
}

func signalContainer(c *ssh.Client, app spec.Application, prog spec.Program, dockerPath, name string, signal string) error {
	cmd := []string{
		dockerPath,
		"kill",
		"--signal", signal,
		name,
	}
	return c.Run(strings.Join(cmd, " "), false)
}

func executeContainer(c *ssh.Client, app spec.Application, prog spec.Program, dockerPath, name string, command string) error {
	cmd := []string{
		dockerPath,
		"exec",
		name,
		command,
	}
	return c.Run(strings.Join(cmd, " "), false)
}
