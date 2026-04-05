package cfg

import (
	"os"

	"github.com/spf13/pflag"
)

// Loader builds a Configuration using a chainable builder pattern.
type Loader struct {
	c   Configuration
	err error
}

// NewLoader returns a new Loader with an empty configuration.
func NewLoader() *Loader {
	return &Loader{
		c: Configuration{},
	}
}

// Parse parses YAML configuration from b.
func (l *Loader) Parse(filename string, b []byte) *Loader {
	if l.err != nil {
		return l
	}
	l.err = l.c.Parse(filename, b)
	return l
}

// ParseFile reads and parses a YAML configuration file.
func (l *Loader) ParseFile(filename string) *Loader {
	if l.err != nil {
		return l
	}
	l.err = l.c.ParseFile(filename)
	return l
}

// ParseFileIfExists parses a YAML configuration file if it exists, skipping silently otherwise.
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

// ApplyEnvironment applies BULLET_* environment variables.
func (l *Loader) ApplyEnvironment() *Loader {
	if l.err != nil {
		return l
	}
	l.err = l.c.ApplyEnvironment()
	return l
}

// ApplyFlags applies explicitly set CLI flags.
func (l *Loader) ApplyFlags(flags *pflag.FlagSet) *Loader {
	if l.err != nil {
		return l
	}
	l.err = l.c.ApplyFlags(flags)
	return l
}

// Configuration returns the built Configuration and any error encountered during loading.
func (l *Loader) Configuration() (Configuration, error) {
	return l.c, l.err
}
