package docker

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
)

var dockerfileTpl = template.Must(template.New("").Parse(dockerfileTplText))

type InstallOptions struct {
	DockerPath string
}

func Install(c *ssh.Client, app spec.Application, proc spec.Process, options InstallOptions) error {
	b := bytes.Buffer{}
	err := dockerfileTpl.Execute(&b, struct {
		Process    spec.Process
		DockerPath string
	}{
		Process:    proc,
		DockerPath: options.DockerPath,
	})
	if err != nil {
		return err
	}

	appDir := fmt.Sprintf("/opt/%s", app.Identifier)

	err = c.Push(fmt.Sprintf("%s/Dockerfile.%s", appDir, proc.Name), 0644, int64(b.Len()), &b)
	if err != nil {
		return err
	}

	imageName := fmt.Sprintf("%s_%s", app.Identifier, proc.Name)
	err = c.Run(fmt.Sprintf("docker build -t %s -f %s/Dockerfile.%s %s", imageName, appDir, proc.Name, appDir))
	if err != nil {
		return err
	}

	return Run(c, app, proc, RunOptions{
		DockerPath: options.DockerPath,
	})
}

type RunOptions struct {
	DockerPath string
}

func Run(c *ssh.Client, app spec.Application, proc spec.Process, options RunOptions) error {
	name := fmt.Sprintf("%s_%s", app.Identifier, proc.Name)
	appDir := fmt.Sprintf("/opt/%s", app.Identifier)
	portArgs := ""
	for _, p := range proc.Ports {
		portArgs += fmt.Sprintf(" -p %s", p)
	}
	imageName := fmt.Sprintf("%s_%s", app.Identifier, proc.Name)

	cmds := []string{
		fmt.Sprintf("%s stop -t 2 %s", options.DockerPath, name),
		fmt.Sprintf("%s rm %s", options.DockerPath, name),
		fmt.Sprintf("%s run -d --name %s -v %s/current:/%s -w /%s --env-file %s/env %s %s %s", options.DockerPath, name, appDir, app.Identifier, app.Identifier, appDir, portArgs, imageName, proc.Command),
	}
	for _, cmd := range cmds {
		err := c.Run(cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

const dockerfileTplText = `FROM {{.Process.Image}}

{{if .Process.PreScript}}
RUN {{.Process.PreScript}}
{{end}}
`
