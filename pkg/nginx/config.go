package nginx

import "github.com/baby-someday/isucon/pkg/remote"

type Config struct {
	Servers []Server `yml:"servers"`
}

type Server struct {
	Host string     `yml:"host"`
	SSH  remote.SSH `yml:"ssh"`
	Bin  string     `yml:"bin"`
	Log  Log        `yml:"log"`
}

type Log struct {
	Access string `yml:"access"`
	Error  string `yml:"error"`
}
