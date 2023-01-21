package config

type Config struct {
	Type        string
	OutFile     string
	OutPkg      string
	PublicValue bool
	Getters     bool
	Setters     bool
	Switch      bool
}
