package servermaster

import (
	"errors"

	"github.com/baby-someday/isucon/pkg/remote"
)

type Config struct {
	Servers []ServerMaster `yaml:"servers"`
}

type ServerMaster struct {
	Name string     `yaml:"name"`
	Host string     `yaml:"host"`
	SSH  remote.SSH `yaml:"ssh"`
}

func FindServerMaster(name string, serverMasters []ServerMaster) (*ServerMaster, error) {
	for _, serverMaster := range serverMasters {
		if serverMaster.Name == name {
			return &serverMaster, nil
		}
	}
	return nil, errors.New("Server could not be found")
}
