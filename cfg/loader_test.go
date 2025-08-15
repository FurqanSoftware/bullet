package cfg

import (
	"testing"

	"github.com/FurqanSoftware/bullet/spec"
	. "github.com/onsi/gomega"
	"github.com/spf13/pflag"
)

func TestLoaderParse(t *testing.T) {
	g := NewWithT(t)

	l := NewLoader()
	cfg, err := l.Parse("Bulletcfg.production", []byte(bulletcfgProduction)).Configuration()
	g.Expect(err).To(BeNil())
	g.Expect(cfg).To(Equal(Configuration{
		Hosts:    "@./manifest.yml",
		Identity: "./id_rsa",
		Vars: spec.Vars{
			"BASE": "https://bullet.fstools.co",
			"ENV":  "production",
		},
	}))
}

func TestLoaderApplyFlags(t *testing.T) {
	g := NewWithT(t)

	flags := pflag.FlagSet{}
	AddFlags(&flags)
	err := flags.Parse([]string{"-H", "@./manifest-staging-alt.yml", "--identity=./id_rsa_staging-alt"})
	g.Expect(err).To(BeNil())

	l := NewLoader()
	cfg, err := l.Parse("Bulletcfg.staging", []byte(bulletcfgStaging)).
		ApplyFlags(&flags).
		Configuration()
	g.Expect(err).To(BeNil())
	g.Expect(cfg).To(Equal(Configuration{
		Hosts:    "@./manifest-staging-alt.yml",
		Identity: "./id_rsa_staging-alt",
		Vars: spec.Vars{
			"BASE": "https://bullet-staging.fstools.co",
			"ENV":  "staging",
		},
	}))
}

const (
	bulletcfgProduction = `
hosts: '@./manifest.yml'
identity: ./id_rsa

vars:
  BASE: https://bullet.fstools.co
  ENV: production
`

	bulletcfgStaging = `
hosts: '@./manifest-staging.yml'
identity: ./id_rsa_staging

vars:
  BASE: https://bullet-staging.fstools.co
  ENV: staging
`
)
