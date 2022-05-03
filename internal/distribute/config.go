package distribute

type Config struct {
	Dst     string   `yml:"dst"`
	Lock    string   `yml:"lock"`
	Ignore  []string `yml:"ignore"`
	Command string   `yml:"command"`
}
