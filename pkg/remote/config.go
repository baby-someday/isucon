package remote

type Network struct {
	Servers []Server `yaml:"servers"`
}

type Server struct {
	Host         string        `yaml:"host"`
	SSH          SSH           `yaml:"ssh"`
	Git          Git           `yaml:"git"`
	Environments []Environment `yaml:"environments"`
}

type Git struct {
	Bin string `yaml:"bin"`
}

type Environment struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
