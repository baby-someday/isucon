package mysql

import "github.com/baby-someday/isucon/pkg/remote"

type Config struct {
	Servers []Server `yml:"servers"`
}

type Server struct {
	Host     string     `yml:"host"`
	Defaults string     `yml:"defaults"`
	SSH      remote.SSH `yml:"ssh"`
	Bin      Bin        `yml:"bin"`
	Log      Log        `yml:"log"`
}

type Bin struct {
	MySQL      string `yml:"mysql"`
	MySQLAdmin string `yml:"mysqladmin"`
}

type Log struct {
	Slow string `yml:"slow"`
}
