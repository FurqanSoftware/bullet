package cfg

import (
	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	SSHRetries int           `envconfig:"SSH_RETRIES"`
	SSHTimeout time.Duration `envconfig:"SSH_TIMEOUT" default:"30s"`
}

func Load(c *Configuration) error {
	err := envconfig.Process("BULLET", c)
	return err
}

var Current Configuration

func LoadCurrent() error {
	return Load(&Current)
}
