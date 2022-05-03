package remote

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/ssh"
)

const (
	AUTHENTICATION_METHOD_PASSWORD = "password"
	AUTHENTICATION_METHOD_KEY      = "key"
)

type AuthenticationMethod interface {
	makeConfig() *ssh.ClientConfig
}

type PasswordAuthentication struct {
	User     string
	Password string
}

func (p PasswordAuthentication) makeConfig() *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: p.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(p.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
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
		// TODO
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
