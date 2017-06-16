package systemd

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
)

var serviceTpl = template.Must(template.New("").Parse(serviceTplText))

func Install(c *ssh.Client, app spec.Application, proc spec.Process) error {
	b := bytes.Buffer{}
	err := serviceTpl.Execute(&b, struct {
		Application spec.Application
		Process     spec.Process
	}{
		Application: app,
		Process:     proc,
	})
	if err != nil {
		return err
	}

	err = c.Scp(fmt.Sprintf("/etc/systemd/system/%s.service", proc.Name), 0644, int64(b.Len()), &b)
	if err != nil {
		return err
	}

	return Enable(c, proc)
}

func Enable(c *ssh.Client, proc spec.Process) error {
	return c.Run(fmt.Sprintf("systemctl enable %s.service", proc.Name))
}

func Disable(c *ssh.Client, proc spec.Process) error {
	return c.Run(fmt.Sprintf("systemctl disable %s.service", proc.Name))
}

func Start(c *ssh.Client, proc spec.Process) error {
	return c.Run(fmt.Sprintf("systemctl start %s.service", proc.Name))
}

func Stop(c *ssh.Client, proc spec.Process) error {
	return c.Run(fmt.Sprintf("systemctl stop %s.service", proc.Name))
}

func Restart(c *ssh.Client, proc spec.Process) error {
	return c.Run(fmt.Sprintf("systemctl restart %s.service", proc.Name))
}

const serviceTplText = `[Unit]
Description=

[Service]
ExecStart=docker run -v /opt/bullet/{{.Application.Identifier}}/current:/{{.Application.Identifier}} -w /{{.Application.Identifier}} alpine:3.5 {{.Process.Command}}
WorkingDirectory=/opt/bullet/{{.Application.Identifier}}/current
Restart=always
Environment=

[Install]
WantedBy=multi-user.target
`
