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

func Install(c *ssh.Client, app spec.Application, proc spec.Program, options InstallOptions) error {
	b := bytes.Buffer{}
	err := dockerfileTpl.Execute(&b, struct {
		Program    spec.Program
		DockerPath string
	}{
		Program:    proc,
		DockerPath: options.DockerPath,
	})
	if err != nil {
		return err
	}

	appDir := fmt.Sprintf("/opt/%s", app.Identifier)

	err = c.Push(fmt.Sprintf("%s/Dockerfile.%s", appDir, proc.Key), 0644, int64(b.Len()), &b)
	if err != nil {
		return err
	}

	imageName := fmt.Sprintf("%s_%s", app.Identifier, proc.Key)
	err = c.Run(fmt.Sprintf("docker build -t %s -f %s/Dockerfile.%s %s", imageName, appDir, proc.Key, appDir))
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

func Run(c *ssh.Client, app spec.Application, proc spec.Program, options RunOptions) error {
	name := fmt.Sprintf("%s_%s", app.Identifier, proc.Key)
	appDir := fmt.Sprintf("/opt/%s", app.Identifier)
	portArgs := ""
	for _, p := range proc.Ports {
		portArgs += fmt.Sprintf(" -p %s", p)
	}
	imageName := fmt.Sprintf("%s_%s", app.Identifier, proc.Key)

	cmds := []string{
		fmt.Sprintf("%s stop -t 2 %s > /dev/null 2>&1 || true", options.DockerPath, name),
		fmt.Sprintf("%s rm %s > /dev/null 2>&1 || true", options.DockerPath, name),
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

const dockerfileTplText = `FROM {{.Program.Image}}

{{range .Program.PreScript}}
	RUN {{.}}
{{end}}
`
