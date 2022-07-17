package remote

type Network struct {
	Servers []Server `yaml:"servers"`
}

type Server struct {
	Name         string        `yaml:"name"`
	Git          Git           `yaml:"git"`
	Unzip        Unzip         `yaml:"unzip"`
	Environments []Environment `yaml:"environments"`
}

type Git struct {
	Bin string `yaml:"bin"`
}

type Unzip struct {
	Bin string `yaml:"bin"`
}

type Environment struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
