package main

import "testing"

func TestParseTimerUnit(t *testing.T) {
	cases := []struct {
		unit     string
		wantApp  string
		wantJob  string
		wantOK   bool
	}{
		{"bullet_myapp_cleanup.timer", "myapp", "cleanup", true},
		{"bullet_my_app_cleanup.timer", "my_app", "cleanup", true},
		{"bullet_a_b_c.timer", "a_b", "c", true},
		{"unrelated.timer", "", "", false},
		{"bullet_only.timer", "", "", false},
		{"bullet_myapp_cleanup.service", "", "", false},
	}
	for _, c := range cases {
		t.Run(c.unit, func(t *testing.T) {
			app, job, _, ok := parseTimerUnit(c.unit)
			if ok != c.wantOK || app != c.wantApp || job != c.wantJob {
				t.Errorf("parseTimerUnit(%q) = (%q, %q, %v), want (%q, %q, %v)",
					c.unit, app, job, ok, c.wantApp, c.wantJob, c.wantOK)
			}
		})
	}
}

func TestParseShowOutput(t *testing.T) {
	out := "ActiveState=active\nLastTriggerUSec=1714000000000000\nUnused=ignored=line\n\n"
	got := parseShowOutput(out)
	if got["ActiveState"] != "active" {
		t.Errorf("ActiveState = %q, want active", got["ActiveState"])
	}
	if got["LastTriggerUSec"] != "1714000000000000" {
		t.Errorf("LastTriggerUSec = %q", got["LastTriggerUSec"])
	}
	if got["Unused"] != "ignored=line" {
		t.Errorf("Unused = %q, want ignored=line (Cut splits on first =)", got["Unused"])
	}
}

func TestParseUSecTimestamp(t *testing.T) {
	cases := []struct {
		in   string
		want int64
	}{
		{"", 0},
		{"0", 0},
		{"1714000000000000", 1714000000},
		{"garbage", 0},
	}
	for _, c := range cases {
		got := parseUSecTimestamp(c.in)
		if got != c.want {
			t.Errorf("parseUSecTimestamp(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}
