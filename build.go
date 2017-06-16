package bullet

import (
	"os"
	"os/exec"

	"github.com/FurqanSoftware/bullet/spec"
)

func Build(spec *spec.Spec) error {
	cmd := exec.Command("/bin/bash", "-c", spec.Build.Script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
