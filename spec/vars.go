package spec

import (
	"bytes"
	"text/template"
)

type Vars map[string]any

func (v Vars) Expand(text string) (string, error) {
	t, err := template.New("").Parse(text)
	if err != nil {
		return "", err
	}
	b := bytes.Buffer{}
	err = t.Execute(&b, struct {
		Vars Vars
	}{
		Vars: v,
	})
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
