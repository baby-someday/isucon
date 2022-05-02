package remote

type Network struct {
	Servers []Server `yaml:"servers"`
}

type Server struct {
	Host           string        `yaml:"host"`
	Authentication string        `yaml:"authentication"`
	SSH            SSH           `yaml:"ssh"`
	Environments   []Environment `yaml:"environments"`
}

type SSH struct {
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	PublicKeyPath string `yaml:"pubkey"`
}

type Environment struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
