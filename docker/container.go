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

func ScaleContainer(c *ssh.Client, app spec.Application, prog spec.Program, n int, options ScaleContainerOptions) error {
	conts, err := ListContainers(c, app, prog, ListContainersOptions{
		DockerPath: options.DockerPath,
	})
	if err != nil {
		return nil
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
			return nil
		}
	}
	for ; d < 0; d++ {
		name := fmt.Sprintf("%s_%s_%d", app.Identifier, prog.Key, lastNo-d)
		image := fmt.Sprintf("%s_%s", app.Identifier, prog.Key)
		err = createContainer(c, app, prog, options.DockerPath, image, name, lastNo-d)
		if err != nil {
			return nil
		}
	}
	return nil
}

type RunContainerOptions struct {
	DockerPath string
}

func RunContainer(c *ssh.Client, app spec.Application, prog spec.Program, options RunContainerOptions) error {
	name := fmt.Sprintf("%s_%s_0", app.Identifier, prog.Key)
	image := fmt.Sprintf("%s_%s", app.Identifier, prog.Key)

	conts, err := ListContainers(c, app, prog, ListContainersOptions{
		DockerPath: options.DockerPath,
	})
	if err != nil {
		return nil
	}
	if len(conts) > 0 {
		err := deleteContainer(c, app, prog, options.DockerPath, name)
		if err != nil {
			return nil
		}
	}

	err = createAttachContainer(c, app, prog, options.DockerPath, image, name)
	if err != nil {
		return nil
	}

	conts, err = ListContainers(c, app, prog, ListContainersOptions{
		DockerPath: options.DockerPath,
	})
	if err != nil {
		return nil
	}
	if len(conts) > 0 {
		err := deleteContainer(c, app, prog, options.DockerPath, name)
		if err != nil {
			return nil
		}
	}
	return nil
}

func createContainer(c *ssh.Client, app spec.Application, prog spec.Program, dockerPath, image, name string, no int) error {
	appDir := fmt.Sprintf("/opt/%s", app.Identifier)

	cmd := []string{
		dockerPath,
		"run",
		"-d",
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

	cmd = append(
		cmd,
		"--restart", "always",
		"-v", appDir+"/current:/"+app.Identifier,
		"-w", "/"+app.Identifier,
		image,
		prog.Command,
	)

	return c.Run(strings.Join(cmd, " "))
}

func createAttachContainer(c *ssh.Client, app spec.Application, prog spec.Program, dockerPath, image, name string) error {
	appDir := fmt.Sprintf("/opt/%s", app.Identifier)

	cmds := []string{
		fmt.Sprintf("%s run -ti --env-file %s/env --name %s -v %s/current:/%s -w /%s %s %s", dockerPath, appDir, name, appDir, app.Identifier, app.Identifier, image, prog.Command),
	}
	for _, cmd := range cmds {
		err := c.RunPTY(cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteContainer(c *ssh.Client, app spec.Application, prog spec.Program, dockerPath, name string) error {
	cmds := []string{
		fmt.Sprintf("%s stop -t 2 %s > /dev/null 2>&1 || true", dockerPath, name),
		fmt.Sprintf("%s rm %s > /dev/null 2>&1 || true", dockerPath, name),
	}
	for _, cmd := range cmds {
		err := c.Run(cmd)
		if err != nil {
			return err
		}
	}
	return nil
}
