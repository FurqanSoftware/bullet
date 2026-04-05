package spec

import (
	"bytes"
	"text/template"
)

// Vars is a map of template variables used for expanding values in the Bulletspec.
type Vars map[string]any

// Expand expands Go template expressions in text using the variables.
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
