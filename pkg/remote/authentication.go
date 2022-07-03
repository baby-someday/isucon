package remote

import (
	"errors"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

type AuthenticationMethod interface {
	makeConfig() (*ssh.ClientConfig, error)
}

type PasswordAuthentication struct {
	User     string
	Password string
}

func (p PasswordAuthentication) makeConfig() (*ssh.ClientConfig, error) {
	return &ssh.ClientConfig{
		User: p.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(p.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, nil
}

type KeyAuthentication struct {
	User           string
	PrivateKeyPath string
}

func (p KeyAuthentication) makeConfig() (*ssh.ClientConfig, error) {
	key, err := ioutil.ReadFile(p.PrivateKeyPath)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	return &ssh.ClientConfig{
		User: p.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, nil
}

func MakeAuthenticationMethod(config SSH) (AuthenticationMethod, error) {
	var authenticationMethod AuthenticationMethod
	if config.PrivateKeyPath != "" {
		authenticationMethod = KeyAuthentication{
			User:           config.User,
			PrivateKeyPath: config.PrivateKeyPath,
		}
	} else if config.Password != "" {
		authenticationMethod = PasswordAuthentication{
			User:     config.User,
			Password: config.Password,
		}
	} else {
		return nil, errors.New("Should specify either ssh.password or ssh.privatekey")
	}

	return authenticationMethod, nil
}
