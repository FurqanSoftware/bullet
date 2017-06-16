package systemd

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
)

var serviceTpl = template.Must(template.New("").Parse(serviceTplText))

type InstallOptions struct {
	DockerPath string
}

func Install(c *ssh.Client, app spec.Application, proc spec.Process, options InstallOptions) error {
	b := bytes.Buffer{}
	err := serviceTpl.Execute(&b, struct {
		Application spec.Application
		Process     spec.Process
		DockerPath  string
	}{
		Application: app,
		Process:     proc,
		DockerPath:  options.DockerPath,
	})
	if err != nil {
		return err
	}

	err = c.Push(fmt.Sprintf("/etc/systemd/system/%s.service", proc.Name), 0644, int64(b.Len()), &b)
	if err != nil {
		return err
	}

	err = c.Run("systemctl daemon-reload")
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
Description={{.Application.Name}} container for {{.Process.Name}}
After=docker.service
Requires=docker.service

[Service]
TimeoutStartSec=0
ExecStartPre=-{{.DockerPath}} stop {{.Application.Identifier}}_{{.Process.Name}}
ExecStartPre=-{{.DockerPath}} rm {{.Application.Identifier}}_{{.Process.Name}}
ExecStartPre={{.DockerPath}} pull {{.Process.Image}}
ExecStart={{.DockerPath}} run --name {{.Application.Identifier}}_{{.Process.Name}} -v /opt/{{.Application.Identifier}}/current:/{{.Application.Identifier}} -w /{{.Application.Identifier}} --env-file /opt/{{.Application.Identifier}}/env {{range $p := .Process.Ports}}-p {{$p}} {{end}} {{.Process.Image}} {{.Process.Command}}
ExecStop={{.DockerPath}} stop -t 2 {{.Application.Identifier}}_{{.Process.Name}}
Restart=always
Environment=

[Install]
WantedBy=multi-user.target
`
