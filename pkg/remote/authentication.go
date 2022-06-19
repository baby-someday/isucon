package remote

import (
	"errors"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

const (
	AUTHENTICATION_METHOD_PASSWORD = "password"
	AUTHENTICATION_METHOD_KEY      = "key"
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

func MakeAuthenticationMethod(server Server) (AuthenticationMethod, error) {
	var authenticationMethod AuthenticationMethod
	switch server.Authentication {
	case AUTHENTICATION_METHOD_PASSWORD:
		authenticationMethod = PasswordAuthentication{
			User:     server.SSH.User,
			Password: server.SSH.Password,
		}

	case AUTHENTICATION_METHOD_KEY:
		authenticationMethod = KeyAuthentication{
			User:           server.SSH.User,
			PrivateKeyPath: server.SSH.PrivateKeyPath,
		}
		break

	default:
		return nil, errors.New(fmt.Sprintf(
			"authentication should be followings: %s, %s",
			AUTHENTICATION_METHOD_PASSWORD,
			AUTHENTICATION_METHOD_KEY,
		))
	}

	return authenticationMethod, nil
}
