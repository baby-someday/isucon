package remote

import (
	"fmt"

	"github.com/baby-someday/isucon/pkg/util"
)

func Lock(path, host string, authenticationMethod AuthenticationMethod) error {
	client, err := NewClient(
		host,
		authenticationMethod,
	)
	if err != nil {
		return util.HandleError(err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return util.HandleError(err)
	}
	defer session.Close()

	commandToMakeLockFile := fmt.Sprintf("[ ! -f %s ] && touch %s", path, path)
	err = session.Run(commandToMakeLockFile)
	if err != nil {
		return util.HandleError(err)
	}

	return nil
}

func Unlock(path, host string, authenticationMethod AuthenticationMethod) error {
	client, err := NewClient(
		host,
		authenticationMethod,
	)
	if err != nil {
		return util.HandleError(err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return util.HandleError(err)
	}
	defer session.Close()

	return session.Run(fmt.Sprintf("[ -f %s ] && rm -f %s", path, path))
}
