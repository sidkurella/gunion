package config

type OutputConfig struct {
	OutType     string
	OutFile     string
	OutPkg      string
	PublicValue bool
	Getters     bool
	Setters     bool
	Switch      bool
	Default     bool
}

type InputConfig struct {
	Source string
	Type   string
}
