package spec

import (
	"io/ioutil"
	"os"
	"sort"

	yaml "gopkg.in/yaml.v2"
)

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

func ParseFile(filename string) (*Spec, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, &Error{"ParseFile", err}
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, &Error{"ParseFile", err}
	}
	return Parse(filename, b)
}
