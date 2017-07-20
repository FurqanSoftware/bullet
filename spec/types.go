package spec

type Spec struct {
	Application Application
}

type Application struct {
	Name       string
	Identifier string
	Deploy     Deploy

	Programs    map[string]Program
	ProgramKeys []string
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
	Container Container
	Ports     []string
}

type Container struct {
	Dockerfile string
	Image      string
}
