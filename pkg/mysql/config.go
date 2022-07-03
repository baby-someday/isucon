package mysql

type Config struct {
	Servers []Server `yml:"servers"`
}

type Server struct {
	Host string `yml:"host"`
	Bin  string `yml:"bin"`
}

type Log struct {
	Slow string `yml:"slow"`
}
