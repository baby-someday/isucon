package remote

import (
	"bytes"
)

func Exec(host, command string, authenticationMethod AuthenticationMethod) ([]byte, error) {
	client, err := newClient(host, authenticationMethod)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var stdout bytes.Buffer
	session.Stdout = &stdout
	if err := session.Run(command); err != nil {
		return nil, err
	}
	return stdout.Bytes(), nil
}
