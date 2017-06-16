package spec

type Spec struct {
	Application Application
}

type Application struct {
	Name       string
	Identifier string
	Build      Build
	Package    Package
	Deploy     Deploy
	Programs   map[string]Program
}

type Build struct {
	Image  string
	Script []string
}

type Package struct {
	Contents []string
}

type Deploy struct {
}

type Program struct {
	Key       string `yaml:"-"`
	Name      string
	Command   string
	Image     string
	Ports     []string
	PreScript []string
}
