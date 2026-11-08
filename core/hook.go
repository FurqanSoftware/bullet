package core

import (
	"fmt"

	"github.com/FurqanSoftware/bullet/cfg"
	"github.com/FurqanSoftware/bullet/distro"
	_ "github.com/FurqanSoftware/bullet/distro/ubuntu"
	"github.com/FurqanSoftware/bullet/scope"
	"github.com/FurqanSoftware/bullet/spec"
	"github.com/FurqanSoftware/pog"
)

// HookRun runs the specified post-deploy hooks, honouring each hook's nodes setting.
func HookRun(s scope.Scope, g cfg.Configuration, keys []string) error {
	hooks := []spec.Hook{}
	for _, k := range keys {
		hook, ok := s.Spec.Application.Hooks.PostDeployHook(k)
		if !ok {
			return fmt.Errorf("bad hook key %q", k)
		}
		hooks = append(hooks, hook)
	}
	return runHooks(s, g, hooks)
}

func runHooks(s scope.Scope, g cfg.Configuration, hooks []spec.Hook) error {
	if len(hooks) == 0 || len(s.Nodes) == 0 {
		return nil
	}

	for _, hook := range hooks {
		prog, ok := s.Spec.Application.Programs[hook.ProgramKey()]
		if !ok {
			return fmt.Errorf("hook %s: unknown program %q", hook.Key, hook.ProgramKey())
		}

		nodes := s.Nodes[:1]
		if hook.AllNodes() {
			nodes = s.Nodes
		}

		for _, n := range nodes {
			pog.SetStatus(pogConnecting(n))
			c, err := sshDial(n, g)
			if err != nil {
				return err
			}
			pog.SetStatus(nil)

			d, err := distro.New(c)
			if err != nil {
				return err
			}

			pog.Infof("Running hook %s on %s", hook.Key, n.Label())
			err = d.RunHook(s.Spec.Application, prog, hook)
			if err != nil {
				return fmt.Errorf("hook %s failed on %s: %w", hook.Key, n.Label(), err)
			}
			pog.Infof("Finished hook %s on %s", hook.Key, n.Label())
		}
	}
	return nil
}
