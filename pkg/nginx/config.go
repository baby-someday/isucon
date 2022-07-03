package nginx

type Config struct {
	Servers []Server `yml:"servers"`
}

type Server struct {
	Name string `yml:"name"`
	Bin  string `yml:"bin"`
	Log  Log    `yml:"log"`
}

type Log struct {
	Access string `yml:"access"`
	Error  string `yml:"error"`
}
