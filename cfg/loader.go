package cfg

import (
	"os"

	"github.com/spf13/pflag"
)

type Loader struct {
	c   Configuration
	err error
}

func NewLoader() *Loader {
	return &Loader{
		c: Configuration{},
	}
}

func (l *Loader) Parse(filename string, b []byte) *Loader {
	if l.err != nil {
		return l
	}
	l.err = l.c.Parse(filename, b)
	return l
}

func (l *Loader) ParseFile(filename string) *Loader {
	if l.err != nil {
		return l
	}
	l.err = l.c.ParseFile(filename)
	return l
}

func (l *Loader) ParseFileIfExists(filename string) *Loader {
	if l.err != nil {
		return l
	}
	_, err := os.Stat(filename)
	if !os.IsNotExist(err) {
		l.err = l.c.ParseFile(filename)
	}
	return l
}

func (l *Loader) ApplyEnvironment() *Loader {
	if l.err != nil {
		return l
	}
	l.err = l.c.ApplyEnvironment()
	return l
}

func (l *Loader) ApplyFlags(flags *pflag.FlagSet) *Loader {
	if l.err != nil {
		return l
	}
	l.err = l.c.ApplyFlags(flags)
	return l
}

func (l *Loader) Configuration() (Configuration, error) {
	return l.c, l.err
}
