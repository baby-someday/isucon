package remote

import (
	"fmt"

	"github.com/baby-someday/isucon/pkg/util"
)

func Lock(path, host string, authenticationMethod AuthenticationMethod) error {
	_, err := Exec(
		host,
		fmt.Sprintf("[ ! -f %s ] && touch %s", path, path),
		make([]Environment, 0),
		authenticationMethod,
	)
	if err != nil {
		return util.HandleError(err)
	}
	return nil
}

func Unlock(path, host string, authenticationMethod AuthenticationMethod) error {
	_, err := Exec(
		host,
		fmt.Sprintf("[ -f %s ] && rm -f %s", path, path),
		make([]Environment, 0),
		authenticationMethod,
	)
	if err != nil {
		return util.HandleError(err)
	}
	return nil
}
