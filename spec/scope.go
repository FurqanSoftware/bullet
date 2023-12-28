package spec

import (
	"bytes"
	"io"
	"os"
	"text/template"

	"gopkg.in/yaml.v2"
)

type Scope struct {
	Vars map[string]string
}

func (s *Scope) Expand(text string) (string, error) {
	t, err := template.New("").Parse(text)
	if err != nil {
		return "", err
	}
	b := bytes.Buffer{}
	err = t.Execute(&b, s)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func (s *Scope) Parse(filename string, b []byte) error {
	return yaml.Unmarshal(b, &s.Vars)
}

func (s *Scope) ParseFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	return s.Parse(filename, b)
}
