package benchbot

type Config struct {
	TMP        string     `yaml:"tmp"`
	Sleep      int64      `yaml:"sleep"`
	Repository Repository `yaml:"repository"`
}

type Repository struct {
	URL      string   `yaml:"url"`
	Branches []string `yaml:"branches"`
}
