package remote

import "bytes"

func Exec(host, command string, environments []Environment, authenticationMethod AuthenticationMethod) ([]uint8, error) {
	client, session, err := NewSession(host, environments, authenticationMethod)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	defer session.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr
	err = session.Run(command)
	if err != nil {
		return nil, err
	}

	return stdout.Bytes(), nil
}
