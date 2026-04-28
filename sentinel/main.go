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
)

// version is set at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	log.SetPrefix("bullet-sentinel: ")
	log.SetFlags(log.LstdFlags | log.Lmsgprefix)

	configPath := flag.String("config", "/etc/bullet/sentinel.yaml", "path to YAML config file")
	addrFlag := flag.String("addr", "", "HTTP bind address (overrides config)")
	scrapeFlag := flag.Duration("scrape-interval", 0, "minimum interval between host introspections (overrides config)")
	dockerPathFlag := flag.String("docker-path", "", "path to the docker binary (overrides config)")
	flag.Parse()

	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "addr":
			cfg.Addr = *addrFlag
		case "scrape-interval":
			cfg.ScrapeInterval = *scrapeFlag
		case "docker-path":
			cfg.DockerPath = *dockerPathFlag
		}
	})

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
