package remote

type Network struct {
	Servers []Server `yaml:"servers"`
}

type Server struct {
	Host           string        `yaml:"host"`
	Authentication string        `yaml:"authentication"`
	SSH            SSH           `yaml:"ssh"`
	Nginx          Nginx         `yaml:"nginx"`
	Environments   []Environment `yaml:"environments"`
}

type SSH struct {
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	PublicKeyPath string `yaml:"pubkey"`
}

type Nginx struct {
	Log NginxLog `yaml:"log"`
}

type NginxLog struct {
	Access      string              `yaml:"access"`
	Error       string              `yaml:"error"`
	Persistence NginxPersistenceLog `yaml:"persistence"`
}

type NginxPersistenceLog struct {
	Access string `yaml:"access"`
	Error  string `yaml:"error"`
}

type Environment struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
