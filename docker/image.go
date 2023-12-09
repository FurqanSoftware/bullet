package docker

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"

	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/bullet/ssh"
)

type Image struct {
	Application spec.Application
	Program     spec.Program
	ID          string
	Repository  string
}

type GetImageOptions struct {
	DockerPath string
}

func GetImage(c *ssh.Client, app spec.Application, prog spec.Program, options GetImageOptions) (*Image, error) {
	name := fmt.Sprintf("%s_%s", app.Identifier, prog.Key)

	img := Image{}
	b, err := c.Output(fmt.Sprintf("%s ps -a --format '{{.ID}}\t{{.Repository}}'", options.DockerPath))
	if err != nil {
		return nil, err
	}
	s := string(b)

	lines := strings.Split(s, "\n")
	for _, line := range lines {
		parts := strings.Split(line, "\t")
		if parts[1] != name {
			continue
		}
		img = Image{
			Application: app,
			Program:     prog,
			ID:          parts[0],
			Repository:  parts[1],
		}
		break
	}
	return &img, nil
}

type BuildImageOptions struct {
	DockerPath string
}

func BuildImage(c *ssh.Client, app spec.Application, prog spec.Program, options BuildImageOptions) (bool, error) {
	if prog.Container.Dockerfile != "" {
		return buildImageDockerfile(c, app, prog, options)
	} else {
		return buildImageDockerHub(c, app, prog, options)
	}
}

func buildImageDockerfile(c *ssh.Client, app spec.Application, prog spec.Program, options BuildImageOptions) (bool, error) {
	file, err := os.Open(prog.Container.Dockerfile)
	if err != nil {
		return false, err
	}
	fileBuf, err := io.ReadAll(file)
	if err != nil {
		return false, err
	}
	err = file.Close()
	if err != nil {
		return false, err
	}

	appDir := fmt.Sprintf("/opt/%s", app.Identifier)
	curDir := fmt.Sprintf("%s/current", appDir)

	shaOut, _ := c.Output(fmt.Sprintf("sha256sum %s/Dockerfile.%s", appDir, prog.Key))
	if err != nil {
		return false, err
	}
	shaOutParts := bytes.Fields(shaOut)
	if len(shaOutParts) >= 2 {
		hash := sha256.New()
		hash.Write(fileBuf)
		if hex.EncodeToString(hash.Sum(nil)[:]) == string(shaOutParts[0]) {
			return false, nil
		}
	}

	err = c.Push(fmt.Sprintf("%s/Dockerfile.%s", appDir, prog.Key), 0644, int64(len(fileBuf)), bytes.NewReader(fileBuf))
	if err != nil {
		return false, err
	}

	name := fmt.Sprintf("%s_%s", app.Identifier, prog.Key)
	return true, c.Run(fmt.Sprintf("docker build -t %s -f %s/Dockerfile.%s %s", name, appDir, prog.Key, curDir), true)
}

func buildImageDockerHub(c *ssh.Client, app spec.Application, prog spec.Program, options BuildImageOptions) (bool, error) {
	b := bytes.Buffer{}
	err := dockerfileTpl.Execute(&b, struct {
		Program spec.Program
	}{
		Program: prog,
	})
	if err != nil {
		return false, err
	}

	appDir := fmt.Sprintf("/opt/%s", app.Identifier)

	err = c.Push(fmt.Sprintf("%s/Dockerfile.%s", appDir, prog.Key), 0644, int64(b.Len()), &b)
	if err != nil {
		return false, err
	}

	name := fmt.Sprintf("%s_%s", app.Identifier, prog.Key)
	return true, c.Run(fmt.Sprintf("docker build -t %s -f %s/Dockerfile.%s %s", name, appDir, prog.Key, appDir), true)
}

var dockerfileTpl = template.Must(template.New("").Parse(dockerfileTplText))

const dockerfileTplText = `FROM {{.Program.Container.Image}}`
