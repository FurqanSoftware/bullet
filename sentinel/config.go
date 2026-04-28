package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// Config holds the sentinel runtime configuration. Defaults apply if a field is
// absent from the config file and no env override is set.
type Config struct {
	Addr           string        `yaml:"addr"`
	ScrapeInterval time.Duration `yaml:"scrape_interval"`
	DockerPath     string        `yaml:"docker_path"`
}

func defaultConfig() Config {
	return Config{
		Addr:           "127.0.0.1:9479",
		ScrapeInterval: 15 * time.Second,
		DockerPath:     "/usr/bin/docker",
	}
}

// loadConfig reads YAML config from path (if it exists) and applies env
// overrides. Caller is expected to apply explicit flag overrides afterwards.
func loadConfig(path string) (Config, error) {
	cfg := defaultConfig()

	data, err := os.ReadFile(path)
	switch {
	case err == nil:
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return cfg, fmt.Errorf("parse %s: %w", path, err)
		}
	case errors.Is(err, fs.ErrNotExist):
		// fall through to defaults + env
	default:
		return cfg, fmt.Errorf("read %s: %w", path, err)
	}

	if v := os.Getenv("BULLET_SENTINEL_ADDR"); v != "" {
		cfg.Addr = v
	}
	if v := os.Getenv("BULLET_SENTINEL_SCRAPE_INTERVAL"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return cfg, fmt.Errorf("BULLET_SENTINEL_SCRAPE_INTERVAL: %w", err)
		}
		cfg.ScrapeInterval = d
	}
	if v := os.Getenv("BULLET_SENTINEL_DOCKER_PATH"); v != "" {
		cfg.DockerPath = v
	}

	return cfg, nil
}
