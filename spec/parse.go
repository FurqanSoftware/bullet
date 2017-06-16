package spec

import (
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

func Parse(filename string, r io.Reader) (*Spec, error) {
	spec := Spec{}
	_, err := toml.DecodeReader(r, &spec)
	if err != nil {
		return nil, err
	}
	return &spec, nil
}

func ParseFile(filename string) (*Spec, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, &Error{"ParseFile", err}
	}
	defer f.Close()
	return Parse(filename, f)
}
