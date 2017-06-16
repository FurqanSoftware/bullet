package spec

type Spec struct {
	Application Application
	Build       Build
	Package     Package
	Deploy      Deploy
	Processes   []Process `toml:"process"`
}

type Application struct {
	Name       string
	Identifier string
}

type Build struct {
	Script string
	Image  string
}

type Package struct {
	Contents []string
}

type Deploy struct {
}

type Process struct {
	Name      string
	Command   string
	Image     string
	Ports     []string
	PreScript string
}
