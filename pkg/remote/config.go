package remote

type Network struct {
	Servers []Server `yaml:"servers"`
}

type Server struct {
	Host           string        `yaml:"host"`
	Authentication string        `yaml:"authentication"`
	SSH            SSH           `yaml:"ssh"`
	Nginx          Nginx         `yaml:"nginx"`
	Git            Git           `yaml:"git"`
	Environments   []Environment `yaml:"environments"`
}

type SSH struct {
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	PrivateKeyPath string `yaml:"privatekey"`
}

type Nginx struct {
	Bin string   `yaml:"bin"`
	Log NginxLog `yaml:"log"`
}

type NginxLog struct {
	Access string `yaml:"access"`
	Error  string `yaml:"error"`
}

type NginxPersistenceLog struct {
	Access string `yaml:"access"`
	Error  string `yaml:"error"`
}

type Git struct {
	Bin string `yaml:"bin"`
}

type Environment struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
