package mysql

type Config struct {
	Servers []Server `yml:"servers"`
	Test    string   `yml:"test"`
	TTT     TTT      `yml:"ttt"`
}

type TTT struct {
	Name string `yml:"name"`
}

type Server struct {
	Name     string `yml:"name"`
	Defaults string `yml:"defaults"`
	Bin      Bin    `yml:"bin"`
	Log      Log    `yml:"log"`
}

type Bin struct {
	MySQL      string `yml:"mysql"`
	MySQLAdmin string `yml:"mysqladmin"`
}

type Log struct {
	Slow string `yml:"slow"`
}

type PtQueryDigest struct {
	Bin     string                `yaml:"bin"`
	Dirs    []PtQueryDigestDir    `yaml:"dirs"`
	Presets []PtQueryDigestPreset `yaml:"presets"`
}

type PtQueryDigestDir struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type PtQueryDigestPreset struct {
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
