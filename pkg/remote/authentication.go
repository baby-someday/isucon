package remote

import "golang.org/x/crypto/ssh"

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
