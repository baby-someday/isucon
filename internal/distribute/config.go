package distribute

type Config struct {
	Src     string   `yml:"src"`
	Dst     string   `yml:"dst"`
	Lock    string   `yml:"lock"`
	Ignore  []string `yml:"ignore"`
	Command string   `yml:"command"`
}
