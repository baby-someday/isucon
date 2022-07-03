package remote

import (
	"bytes"
	"errors"
	"fmt"
)

type SSH struct {
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	PrivateKeyPath string `yaml:"privatekey"`
}

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
		return nil, errors.New(fmt.Sprintf("%s: %s", err.Error(), stderr.String()))
	}

	return stdout.Bytes(), nil
}
