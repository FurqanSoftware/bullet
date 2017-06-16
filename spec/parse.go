package spec

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

func Parse(filename string, b []byte) (*Spec, error) {
	spec := Spec{}
	err := yaml.Unmarshal(b, &spec)
	if err != nil {
		return nil, err
	}
	for k, v := range spec.Application.Programs {
		v.Key = k
		spec.Application.Programs[k] = v
	}
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
