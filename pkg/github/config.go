package github

type GitHub struct {
	Token      string     `yaml:"token"`
	Repository Repository `yaml:"repository"`
}

type Repository struct {
	Owner    string   `yaml:"owner"`
	Name     string   `yaml:"name"`
	URL      string   `yaml:"url"`
	Branches []string `yaml:"branches"`
}
