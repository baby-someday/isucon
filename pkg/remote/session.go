package remote

import (
	"golang.org/x/crypto/ssh"
)

func newClient(host string, authenticationMethod AuthenticationMethod) (*ssh.Client, error) {
	return ssh.Dial("tcp", host+":22", authenticationMethod.makeConfig())
}
