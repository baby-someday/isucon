package slack

type Config struct {
	Token   string `yaml:"token"`
	Channel string `yaml:"channel"`
}
