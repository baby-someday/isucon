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
