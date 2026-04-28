package main

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Container is a Bullet-managed Docker container discovered on the host.
type Container struct {
	App      string
	Program  string
	Instance int

	ID     string
	Image  string
	Status string
}

// Up reports whether the container is currently running.
func (c Container) Up() bool {
	return strings.HasPrefix(c.Status, "Up")
}

// Health returns whether the container reports healthy. The second return
// value indicates whether the container has a healthcheck configured at all.
func (c Container) Health() (healthy, hasHealthcheck bool) {
	switch {
	case strings.Contains(c.Status, "(healthy)"):
		return true, true
	case strings.Contains(c.Status, "(unhealthy)"), strings.Contains(c.Status, "(health: starting)"):
		return false, true
	}
	return false, false
}

// reBulletContainer matches Bullet's <app>_<program>_<instance> naming. The
// app and program portions are split by the last underscore in the prefix —
// this means program keys with underscores can't be parsed by the sentinel
// without out-of-band info, but matches Bullet's typical usage.
var reBulletContainer = regexp.MustCompile(`^(.+)_(\d+)$`)

func listContainers(ctx context.Context, dockerPath string) ([]Container, error) {
	cmd := exec.CommandContext(ctx, dockerPath, "ps", "-a", "--format", "{{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("docker ps: %w", err)
	}
	return parseDockerPS(string(out)), nil
}

func parseDockerPS(out string) []Container {
	var conts []Container
	for _, line := range strings.Split(out, "\n") {
		parts := strings.Split(line, "\t")
		if len(parts) != 4 {
			continue
		}
		app, prog, inst, ok := parseContainerName(parts[1])
		if !ok {
			continue
		}
		conts = append(conts, Container{
			App:      app,
			Program:  prog,
			Instance: inst,
			ID:       parts[0],
			Image:    parts[2],
			Status:   parts[3],
		})
	}
	return conts
}

func parseContainerName(name string) (app, program string, instance int, ok bool) {
	m := reBulletContainer.FindStringSubmatch(name)
	if m == nil {
		return "", "", 0, false
	}
	n, err := strconv.Atoi(m[2])
	if err != nil {
		return "", "", 0, false
	}
	prefix := m[1]
	idx := strings.LastIndex(prefix, "_")
	if idx < 1 || idx >= len(prefix)-1 {
		return "", "", 0, false
	}
	return prefix[:idx], prefix[idx+1:], n, true
}
