package spec

import (
	"io"
	"os"
	"sort"

	yaml "gopkg.in/yaml.v2"
)

// Spec represents a parsed Bulletspec file.
type Spec struct {
	Application Application
}

// Parse parses a Bulletspec from YAML bytes.
func Parse(filename string, b []byte) (*Spec, error) {
	spec := Spec{}
	err := yaml.Unmarshal(b, &spec)
	if err != nil {
		return nil, err
	}
	keys := []string{}
	for k, v := range spec.Application.Programs {
		v.Key = k
		spec.Application.Programs[k] = v
		keys = append(keys, k)
	}
	sort.Strings(keys)
	spec.Application.ProgramKeys = keys
	return &spec, nil
}

// ParseFile reads and parses a Bulletspec YAML file.
func ParseFile(filename string) (*Spec, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, &Error{"ParseFile", err}
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, &Error{"ParseFile", err}
	}
	return Parse(filename, b)
}

// ExpandVars expands template variables in the spec's program commands.
func (s *Spec) ExpandVars(vars Vars) error {
	return s.Application.ExpandVars(vars)
}
