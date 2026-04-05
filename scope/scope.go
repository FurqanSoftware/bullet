package scope

import (
	"github.com/FurqanSoftware/bullet/spec"
)

// Scope holds the parsed spec and the set of target nodes for a command.
type Scope struct {
	Spec     *spec.Spec
	Nodes    []Node
	Identity string
}
