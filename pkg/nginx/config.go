package nginx

type Config struct {
	Servers []Server `yml:"servers"`
}

type Server struct {
	Name string `yml:"name"`
	Bin  string `yml:"bin"`
	Log  Log    `yml:"log"`
}

type Log struct {
	Access string `yml:"access"`
	Error  string `yml:"error"`
}

type ALP struct {
	Bin     string      `yaml:"bin"`
	Dirs    []ALPDir    `yaml:"dirs"`
	Presets []ALPPreset `yaml:"presets"`
}

type ALPDir struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type ALPPreset struct {
	Name           string `yaml:"name"`
	File           string `yaml:"file"`
	M              string `yaml:"m"`
	O              string `yaml:"o"`
	Q              string `yaml:"q"`
	QsIgnoreValues bool   `yaml:"qs-ignore-values"`
	R              string `yaml:"r"`
	ShowFooters    bool   `yaml:"show-footers"`
	Sort           string `yaml:"sort"`
	Extra          string `yaml:"extra"`
}
