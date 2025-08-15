package scope

import (
	"github.com/FurqanSoftware/bullet/spec"
)

type Scope struct {
	Spec     *spec.Spec
	Nodes    []Node
	Identity string
}
