package remote

import "fmt"

func Lock(path, host string, authenticationMethod AuthenticationMethod) error {
	client, err := NewClient(
		host,
		authenticationMethod,
	)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Run(fmt.Sprintf("[ ! -f %s ] && touch %s", path, path))
}

func Unlock(path, host string, authenticationMethod AuthenticationMethod) error {
	client, err := NewClient(
		host,
		authenticationMethod,
	)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Run(fmt.Sprintf("[ -f %s ] && rm -f %s", path, path))
}
