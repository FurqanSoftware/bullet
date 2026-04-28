package main

import (
	"reflect"
	"testing"
)

func TestParseContainerName(t *testing.T) {
	cases := []struct {
		name     string
		wantApp  string
		wantProg string
		wantInst int
		wantOK   bool
	}{
		{"myapp_web_1", "myapp", "web", 1, true},
		{"myapp_worker_42", "myapp", "worker", 42, true},
		{"my_app_web_2", "my_app", "web", 2, true},
		{"a_b_c_3", "a_b", "c", 3, true},
		{"random-container", "", "", 0, false},
		{"name_only", "", "", 0, false},
		{"app_prog_notnum", "", "", 0, false},
		{"_prog_1", "", "", 0, false},
		{"app__1", "", "", 0, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			app, prog, inst, ok := parseContainerName(c.name)
			if ok != c.wantOK || app != c.wantApp || prog != c.wantProg || inst != c.wantInst {
				t.Errorf("parseContainerName(%q) = (%q, %q, %d, %v), want (%q, %q, %d, %v)",
					c.name, app, prog, inst, ok, c.wantApp, c.wantProg, c.wantInst, c.wantOK)
			}
		})
	}
}

func TestParseDockerPS(t *testing.T) {
	out := "abc123\tmyapp_web_1\tnginx:latest\tUp 2 hours (healthy)\n" +
		"def456\tmyapp_web_2\tnginx:latest\tExited (0) 5 minutes ago\n" +
		"ghi789\tunrelated\tredis\tUp 1 day\n" +
		"jkl012\tmyapp_worker_1\tworker:v3\tUp 3 hours\n"

	got := parseDockerPS(out)
	want := []Container{
		{App: "myapp", Program: "web", Instance: 1, ID: "abc123", Image: "nginx:latest", Status: "Up 2 hours (healthy)"},
		{App: "myapp", Program: "web", Instance: 2, ID: "def456", Image: "nginx:latest", Status: "Exited (0) 5 minutes ago"},
		{App: "myapp", Program: "worker", Instance: 1, ID: "jkl012", Image: "worker:v3", Status: "Up 3 hours"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseDockerPS mismatch:\ngot:  %+v\nwant: %+v", got, want)
	}
}

func TestContainerHealth(t *testing.T) {
	cases := []struct {
		status      string
		wantUp      bool
		wantHealthy bool
		wantHasHC   bool
	}{
		{"Up 2 hours (healthy)", true, true, true},
		{"Up 5 minutes (unhealthy)", true, false, true},
		{"Up 10 seconds (health: starting)", true, false, true},
		{"Up 1 day", true, false, false},
		{"Exited (0) 5 minutes ago", false, false, false},
	}
	for _, c := range cases {
		t.Run(c.status, func(t *testing.T) {
			cont := Container{Status: c.status}
			if cont.Up() != c.wantUp {
				t.Errorf("Up() = %v, want %v", cont.Up(), c.wantUp)
			}
			h, has := cont.Health()
			if h != c.wantHealthy || has != c.wantHasHC {
				t.Errorf("Health() = (%v, %v), want (%v, %v)", h, has, c.wantHealthy, c.wantHasHC)
			}
		})
	}
}
