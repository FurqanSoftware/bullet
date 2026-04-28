package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// version is set at build time via -ldflags "-X main.version=...".
var version = "dev"

type config struct {
	Addr           string        `envconfig:"ADDR" default:"127.0.0.1:9479"`
	ScrapeInterval time.Duration `envconfig:"SCRAPE_INTERVAL" default:"15s"`
	DockerPath     string        `envconfig:"DOCKER_PATH" default:"/usr/bin/docker"`
}

func main() {
	log.SetPrefix("bullet-sentinel: ")
	log.SetFlags(log.LstdFlags | log.Lmsgprefix)

	cfg := config{}
	if err := envconfig.Process("BULLET_SENTINEL", &cfg); err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&cfg.Addr, "addr", cfg.Addr, "HTTP bind address")
	flag.DurationVar(&cfg.ScrapeInterval, "scrape-interval", cfg.ScrapeInterval, "minimum interval between host introspections")
	flag.StringVar(&cfg.DockerPath, "docker-path", cfg.DockerPath, "path to the docker binary")
	flag.Parse()

	collector := &Collector{
		DockerPath:     cfg.DockerPath,
		ScrapeInterval: cfg.ScrapeInterval,
		Version:        version,
	}
	srv := newServer(cfg.Addr, collector)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("Listening on %s", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	log.Print("Shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error shutting down: %v", err)
	}

	log.Print("Bye")
}
