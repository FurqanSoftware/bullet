package main

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

// Collector gathers host state and renders Prometheus text-format metrics.
// Snapshots are cached for the configured interval so multiple scrapers don't
// trigger redundant docker/systemctl invocations.
type Collector struct {
	DockerPath     string
	ScrapeInterval time.Duration
	Version        string

	mu       sync.Mutex
	snapshot []byte
	lastRun  time.Time
}

// Render returns the latest metrics output, refreshing it if older than
// ScrapeInterval.
func (c *Collector) Render(ctx context.Context) []byte {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.snapshot != nil && time.Since(c.lastRun) < c.ScrapeInterval {
		return c.snapshot
	}

	c.snapshot = c.gather(ctx)
	c.lastRun = time.Now()
	return c.snapshot
}

func (c *Collector) gather(ctx context.Context) []byte {
	var buf bytes.Buffer

	writeMetric(&buf, "bullet_sentinel_build_info", "Sentinel build info", "gauge",
		[]sample{{labels: map[string]string{"version": c.Version}, value: 1}})

	containers, dockerErr := listContainers(ctx, c.DockerPath)
	timers, timerErr := listTimers(ctx)

	writeMetric(&buf, "bullet_sentinel_docker_up", "Whether the docker query succeeded", "gauge",
		[]sample{{value: boolFloat(dockerErr == nil)}})
	writeMetric(&buf, "bullet_sentinel_systemd_up", "Whether the systemd query succeeded", "gauge",
		[]sample{{value: boolFloat(timerErr == nil)}})

	if dockerErr == nil {
		writeContainerMetrics(&buf, containers)
	}
	if timerErr == nil {
		writeTimerMetrics(&buf, timers)
	}

	return buf.Bytes()
}

func writeContainerMetrics(buf *bytes.Buffer, containers []Container) {
	var (
		upSamples      []sample
		healthySamples []sample
		runningCount   = map[[2]string]int{}
	)
	for _, c := range containers {
		labels := map[string]string{
			"app":      c.App,
			"program":  c.Program,
			"instance": fmt.Sprintf("%d", c.Instance),
		}
		upSamples = append(upSamples, sample{labels: labels, value: boolFloat(c.Up())})

		if healthy, has := c.Health(); has {
			healthySamples = append(healthySamples, sample{labels: labels, value: boolFloat(healthy)})
		}

		if c.Up() {
			runningCount[[2]string{c.App, c.Program}]++
		}
	}

	writeMetric(buf, "bullet_container_up", "Whether a Bullet-managed container is running", "gauge", upSamples)
	writeMetric(buf, "bullet_container_healthy", "Container healthcheck status (1 healthy, 0 unhealthy)", "gauge", healthySamples)

	var countSamples []sample
	for k, n := range runningCount {
		countSamples = append(countSamples, sample{
			labels: map[string]string{"app": k[0], "program": k[1]},
			value:  float64(n),
		})
	}
	writeMetric(buf, "bullet_program_instances_running", "Count of running instances per program", "gauge", countSamples)
}

func writeTimerMetrics(buf *bytes.Buffer, timers []Timer) {
	var activeSamples, lastTriggerSamples, lastResultSamples []sample
	for _, t := range timers {
		labels := map[string]string{"app": t.App, "job": t.Job}
		activeSamples = append(activeSamples, sample{labels: labels, value: boolFloat(t.Active)})
		lastTriggerSamples = append(lastTriggerSamples, sample{labels: labels, value: float64(t.LastTriggerUnix)})
		lastResultSamples = append(lastResultSamples, sample{labels: labels, value: boolFloat(!t.Succeeded() && t.LastResult != "")})
	}
	writeMetric(buf, "bullet_cron_timer_active", "Whether a cron timer is active", "gauge", activeSamples)
	writeMetric(buf, "bullet_cron_last_trigger_timestamp_seconds", "Unix timestamp of last trigger (0 if never)", "gauge", lastTriggerSamples)
	writeMetric(buf, "bullet_cron_last_result", "Last cron run result (0 success or unknown, 1 failure)", "gauge", lastResultSamples)
}

type sample struct {
	labels map[string]string
	value  float64
}

func writeMetric(buf *bytes.Buffer, name, help, typ string, samples []sample) {
	fmt.Fprintf(buf, "# HELP %s %s\n", name, help)
	fmt.Fprintf(buf, "# TYPE %s %s\n", name, typ)
	for _, s := range samples {
		fmt.Fprintf(buf, "%s%s %s\n", name, formatLabels(s.labels), formatFloat(s.value))
	}
}

func formatLabels(labels map[string]string) string {
	if len(labels) == 0 {
		return ""
	}
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	b.WriteByte('{')
	for i, k := range keys {
		if i > 0 {
			b.WriteByte(',')
		}
		fmt.Fprintf(&b, `%s="%s"`, k, escapeLabel(labels[k]))
	}
	b.WriteByte('}')
	return b.String()
}

func escapeLabel(v string) string {
	v = strings.ReplaceAll(v, `\`, `\\`)
	v = strings.ReplaceAll(v, `"`, `\"`)
	v = strings.ReplaceAll(v, "\n", `\n`)
	return v
}

func formatFloat(f float64) string {
	if f == float64(int64(f)) {
		return fmt.Sprintf("%d", int64(f))
	}
	return fmt.Sprintf("%g", f)
}

func boolFloat(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
