package bullet

import (
	"os"
	"os/exec"

	"github.com/FurqanSoftware/bullet/spec"
)

func Build(spec *spec.Spec) error {
	for _, l := range spec.Application.Build.Script {
		cmd := exec.Command("/bin/bash", "-c", l)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
