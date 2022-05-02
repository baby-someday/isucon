package distribute

type Config struct {
	Src     string   `yml:"src"`
	Dst     string   `yml:"dst"`
	Ignore  []string `yml:"ignore"`
	Command string   `yml:"command"`
}
