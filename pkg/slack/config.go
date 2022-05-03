package slack

type Slack struct {
	Token   string `yaml:"token"`
	Channel string `yaml:"channel"`
}

const (
	SEPARATOR = "=================================================="
)
