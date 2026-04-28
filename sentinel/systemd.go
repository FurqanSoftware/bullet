package main

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Timer is a Bullet-managed cron timer discovered on the host.
type Timer struct {
	App string
	Job string

	UnitName    string
	ServiceName string

	Active            bool
	LastTriggerUnix   int64
	LastResult        string
}

// Succeeded returns whether the last invocation reported success. Returns
// false if the job has never run or systemd reports a non-success Result.
func (t Timer) Succeeded() bool {
	return t.LastResult == "success"
}

var reBulletTimer = regexp.MustCompile(`^bullet_(.+)\.timer$`)

func listTimers(ctx context.Context) ([]Timer, error) {
	out, err := exec.CommandContext(ctx, "systemctl", "list-units", "--all", "--no-legend", "--plain", "--type=timer").Output()
	if err != nil {
		return nil, fmt.Errorf("systemctl list-units: %w", err)
	}

	var timers []Timer
	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		app, job, unit, ok := parseTimerUnit(fields[0])
		if !ok {
			continue
		}
		t := Timer{
			App:         app,
			Job:         job,
			UnitName:    unit,
			ServiceName: strings.TrimSuffix(unit, ".timer") + ".service",
		}
		if err := loadTimerDetails(ctx, &t); err != nil {
			return nil, err
		}
		timers = append(timers, t)
	}
	return timers, nil
}

func parseTimerUnit(unit string) (app, job, name string, ok bool) {
	m := reBulletTimer.FindStringSubmatch(unit)
	if m == nil {
		return "", "", "", false
	}
	inner := m[1]
	idx := strings.LastIndex(inner, "_")
	if idx < 1 || idx >= len(inner)-1 {
		return "", "", "", false
	}
	return inner[:idx], inner[idx+1:], unit, true
}

func loadTimerDetails(ctx context.Context, t *Timer) error {
	timerProps, err := showProperties(ctx, t.UnitName, "ActiveState", "LastTriggerUSec")
	if err != nil {
		return err
	}
	t.Active = timerProps["ActiveState"] == "active"
	t.LastTriggerUnix = parseUSecTimestamp(timerProps["LastTriggerUSec"])

	svcProps, err := showProperties(ctx, t.ServiceName, "Result")
	if err != nil {
		return err
	}
	t.LastResult = svcProps["Result"]
	return nil
}

func showProperties(ctx context.Context, unit string, props ...string) (map[string]string, error) {
	args := []string{"show", unit, "--property=" + strings.Join(props, ",")}
	out, err := exec.CommandContext(ctx, "systemctl", args...).Output()
	if err != nil {
		return nil, fmt.Errorf("systemctl show %s: %w", unit, err)
	}
	return parseShowOutput(string(out)), nil
}

func parseShowOutput(out string) map[string]string {
	m := map[string]string{}
	for _, line := range strings.Split(out, "\n") {
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		m[k] = v
	}
	return m
}

// parseUSecTimestamp parses a microseconds-since-epoch string as returned by
// systemctl show for *USec properties. Returns 0 for empty, "0", or invalid
// values.
func parseUSecTimestamp(s string) int64 {
	if s == "" || s == "0" {
		return 0
	}
	usec, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return usec / 1_000_000
}
