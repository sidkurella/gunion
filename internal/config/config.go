package config

type OutputConfig struct {
	OutType     string
	OutFile     string
	OutPkg      string
	Command     string
	PublicValue bool
	Getters     bool
	Setters     bool
	Match       bool
	Default     bool
}

type InputConfig struct {
	Source string
	Type   string
}
