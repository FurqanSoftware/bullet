package spec

import (
	"strings"
	"testing"
	"time"
)

const hookSpecBase = `
application:
  identifier: app
  programs:
    web:
      command: ./web
      container:
        image: alpine
    shell:
      command: /bin/sh
      container:
        image: alpine
`

func TestParseHooks(t *testing.T) {
	s, err := Parse("Bulletspec", []byte(hookSpecBase+`
  hooks:
    post_deploy:
      - key: makeindexes
        command: ./clt db:makeindexes
      - key: purgecache
        command: ./clt cache:purge
        program: web
        nodes: all
        timeout: 10m
`))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	hooks := s.Application.Hooks.PostDeploy
	if len(hooks) != 2 {
		t.Fatalf("want 2 hooks, got %d", len(hooks))
	}

	h := hooks[0]
	if h.Key != "makeindexes" || h.Command != "./clt db:makeindexes" {
		t.Errorf("unexpected hook: %+v", h)
	}
	if h.ProgramKey() != "shell" {
		t.Errorf("want default program shell, got %q", h.ProgramKey())
	}
	if h.AllNodes() {
		t.Error("want default nodes first")
	}
	if h.Timeout != 0 {
		t.Errorf("want no timeout, got %s", h.Timeout)
	}

	h = hooks[1]
	if h.ProgramKey() != "web" {
		t.Errorf("want program web, got %q", h.ProgramKey())
	}
	if !h.AllNodes() {
		t.Error("want nodes all")
	}
	if h.Timeout != 10*time.Minute {
		t.Errorf("want timeout 10m, got %s", h.Timeout)
	}

	if _, ok := s.Application.Hooks.PostDeployHook("purgecache"); !ok {
		t.Error("PostDeployHook: want purgecache found")
	}
	if _, ok := s.Application.Hooks.PostDeployHook("missing"); ok {
		t.Error("PostDeployHook: want missing not found")
	}
}

func TestValidateHooks(t *testing.T) {
	cases := []struct {
		name  string
		hooks string
		want  string
	}{
		{
			name: "missing key",
			hooks: `
      - command: ./x`,
			want: "application.hooks.post_deploy[0].key is required",
		},
		{
			name: "missing command",
			hooks: `
      - key: x`,
			want: "application.hooks.post_deploy[0].command is required",
		},
		{
			name: "duplicate key",
			hooks: `
      - key: x
        command: ./x
      - key: x
        command: ./y`,
			want: `application.hooks.post_deploy[1].key: duplicate key "x"`,
		},
		{
			name: "unknown program",
			hooks: `
      - key: x
        command: ./x
        program: worker`,
			want: `application.hooks.post_deploy[0].program: unknown program "worker"`,
		},
		{
			name: "bad nodes",
			hooks: `
      - key: x
        command: ./x
        nodes: some`,
			want: `application.hooks.post_deploy[0].nodes: unknown value "some"`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := Parse("Bulletspec", []byte(hookSpecBase+`
  hooks:
    post_deploy:`+c.hooks+"\n"))
			if err == nil {
				t.Fatal("want error, got nil")
			}
			if !strings.Contains(err.Error(), c.want) {
				t.Errorf("want error containing %q, got %q", c.want, err)
			}
		})
	}
}

func TestValidateHooksDefaultProgramMissing(t *testing.T) {
	_, err := Parse("Bulletspec", []byte(`
application:
  identifier: app
  programs:
    web:
      command: ./web
      container:
        image: alpine
  hooks:
    post_deploy:
      - key: x
        command: ./x
`))
	if err == nil || !strings.Contains(err.Error(), `program: unknown program "shell"`) {
		t.Errorf("want unknown program shell error, got %v", err)
	}
}
