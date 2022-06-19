package remote

import (
	"golang.org/x/crypto/ssh"
)

func NewClient(host string, authenticationMethod AuthenticationMethod) (*ssh.Client, error) {
	clientConfig, err := authenticationMethod.makeConfig()
	if err != nil {
		return nil, err
	}
	return ssh.Dial("tcp", host+":22", clientConfig)
}

func NewSession(host string, environments []Environment, authenticationMethod AuthenticationMethod) (*ssh.Client, *ssh.Session, error) {
	client, err := NewClient(host, authenticationMethod)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, nil, err
	}

	for _, environment := range environments {
		err = session.Setenv(environment.Name, environment.Value)
		if err != nil {
			return nil, nil, err
		}
	}

	return client, session, nil
}
