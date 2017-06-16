package systemd

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
)

var serviceTpl = template.Must(template.New("").Parse(serviceTplText))

type InstallOptions struct {
	DockerPath string
}

func Install(c *ssh.Client, app spec.Application, proc spec.Program, options InstallOptions) error {
	b := bytes.Buffer{}
	err := serviceTpl.Execute(&b, struct {
		Application spec.Application
		Program     spec.Program
		DockerPath  string
	}{
		Application: app,
		Program:     proc,
		DockerPath:  options.DockerPath,
	})
	if err != nil {
		return err
	}

	err = c.Push(fmt.Sprintf("/etc/systemd/system/%s.service", proc.Key), 0644, int64(b.Len()), &b)
	if err != nil {
		return err
	}

	err = c.Run("systemctl daemon-reload")
	if err != nil {
		return err
	}

	return Enable(c, proc)
}

func Enable(c *ssh.Client, proc spec.Program) error {
	return c.Run(fmt.Sprintf("systemctl enable %s.service", proc.Key))
}

func Disable(c *ssh.Client, proc spec.Program) error {
	return c.Run(fmt.Sprintf("systemctl disable %s.service", proc.Key))
}

func Start(c *ssh.Client, proc spec.Program) error {
	return c.Run(fmt.Sprintf("systemctl start %s.service", proc.Key))
}

func Stop(c *ssh.Client, proc spec.Program) error {
	return c.Run(fmt.Sprintf("systemctl stop %s.service", proc.Key))
}

func Restart(c *ssh.Client, proc spec.Program) error {
	return c.Run(fmt.Sprintf("systemctl restart %s.service", proc.Key))
}

const serviceTplText = `[Unit]
Description={{.Program.Name}}
After=docker.service
Requires=docker.service

[Service]
ExecStart={{.DockerPath}} start {{.Application.Identifier}}_{{.Program.Key}}
ExecStop={{.DockerPath}} stop -t 2 {{.Application.Identifier}}_{{.Program.Key}}
Restart=always

[Install]
WantedBy=multi-user.target
`
