package spec

import (
	"fmt"
	"strings"
)

// Validate checks the spec for missing or invalid fields and returns an error
// describing all problems found.
func (s *Spec) Validate() error {
	var errs []string

	if s.Application.Identifier == "" {
		errs = append(errs, "application.identifier is required")
	}

	if s.Application.Deploy.Current != "" {
		switch s.Application.Deploy.Current {
		case "symlink", "replace":
		default:
			errs = append(errs, fmt.Sprintf("application.deploy.current: unknown method %q (must be \"symlink\" or \"replace\")", s.Application.Deploy.Current))
		}
	}

	for k, prog := range s.Application.Programs {
		prefix := fmt.Sprintf("application.programs.%s", k)

		if prog.Container.Image == "" && prog.Container.Dockerfile == "" {
			errs = append(errs, fmt.Sprintf("%s.container: image or dockerfile is required", prefix))
		}

		for i, p := range prog.Ports {
			parts := strings.SplitN(p, ":", 2)
			if len(parts) != 2 {
				errs = append(errs, fmt.Sprintf("%s.ports[%d]: must be in HOST:CONTAINER format", prefix, i))
			}
		}

		if prog.Healthcheck != nil {
			if prog.Healthcheck.Command == "" {
				errs = append(errs, fmt.Sprintf("%s.healthcheck.command is required", prefix))
			}
		}

		if prog.Reload.Method != "" {
			switch prog.Reload.Method {
			case "signal":
				if prog.Reload.Signal == "" {
					errs = append(errs, fmt.Sprintf("%s.reload.signal is required when method is \"signal\"", prefix))
				}
			case "command":
				if prog.Reload.Command == "" {
					errs = append(errs, fmt.Sprintf("%s.reload.command is required when method is \"command\"", prefix))
				}
			case "restart":
			default:
				errs = append(errs, fmt.Sprintf("%s.reload.method: unknown method %q (must be \"signal\", \"command\", or \"restart\")", prefix, prog.Reload.Method))
			}
		}

		for i, scale := range prog.Scales {
			if scale.N == "" {
				errs = append(errs, fmt.Sprintf("%s.scales[%d].n is required", prefix, i))
			}
		}
	}

	for i, job := range s.Application.Cron.Jobs {
		prefix := fmt.Sprintf("application.cron.jobs[%d]", i)

		if job.Key == "" {
			errs = append(errs, fmt.Sprintf("%s.key is required", prefix))
		}
		if job.Command == "" {
			errs = append(errs, fmt.Sprintf("%s.command is required", prefix))
		}
		if job.Schedule == "" {
			errs = append(errs, fmt.Sprintf("%s.schedule is required", prefix))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("bulletspec validation failed:\n  %s", strings.Join(errs, "\n  "))
	}
	return nil
}
