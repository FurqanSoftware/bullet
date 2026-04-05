package cfg

import (
	"io"
	"os"
	"time"

	"github.com/FurqanSoftware/bullet/spec"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/pflag"
	yaml "gopkg.in/yaml.v2"
)

// Configuration holds the runtime configuration for a Bullet session.
type Configuration struct {
	Hosts    string `envconfig:"HOSTS"`
	Port     int    `envconfig:"PORT"`
	Identity string `envconfig:"IDENTITY"`

	Vars spec.Vars `envconfig:"-"`

	SSHRetries int           `envconfig:"SSH_RETRIES"`
	SSHTimeout time.Duration `envconfig:"SSH_TIMEOUT" default:"30s"`
}

// Parse unmarshals YAML configuration from b into c.
func (c *Configuration) Parse(filename string, b []byte) error {
	return yaml.Unmarshal(b, c)
}

// ParseFile reads and parses a YAML configuration file.
func (c *Configuration) ParseFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return &Error{"ParseFile", err}
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return &Error{"ParseFile", err}
	}
	return c.Parse(filename, b)
}

// ApplyEnvironment applies BULLET_* environment variables to the configuration.
func (c *Configuration) ApplyEnvironment() error {
	return envconfig.Process("BULLET", c)
}

// ApplyFlags overrides configuration values with any explicitly set CLI flags.
func (c *Configuration) ApplyFlags(flags *pflag.FlagSet) error {
	flagHosts := flags.Lookup("hosts")
	if flagHosts.Changed {
		var err error
		c.Hosts, err = flags.GetString("hosts")
		if err != nil {
			return err
		}
	}
	flagPort := flags.Lookup("port")
	if flagPort.Changed {
		var err error
		c.Port, err = flags.GetInt("port")
		if err != nil {
			return err
		}
	}
	flagIdentity := flags.Lookup("identity")
	if flagIdentity.Changed {
		var err error
		c.Identity, err = flags.GetString("identity")
		if err != nil {
			return err
		}
	}
	return nil
}
