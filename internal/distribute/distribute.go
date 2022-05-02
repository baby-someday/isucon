package distribute

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/baby-someday/isucon/pkg/build"
	"github.com/baby-someday/isucon/pkg/remote"
	"golang.org/x/crypto/ssh"
)

type process struct {
	host       string
	client     *ssh.Client
	session    *ssh.Session
	stdout     io.Reader
	stderr     io.Reader
	stdoutFile *os.File
	stderrFile *os.File
}

func Distribute(ctx context.Context, network remote.Network, src, dst, lock, command string, ignore []string) error {
	err := tryToLock(
		lock,
		network,
	)
	if err != nil {
		log.Println("ロックに失敗しました、他のベンチマークが実行中です。")
		return err
	}

	zipPath := path.Join(getOutputPath(), path.Base(src)+".zip")
	err = build.Compress(src, zipPath, ignore)
	if err != nil {
		return err
	}

	processes := []process{}

	// TODO: Closeちゃんとやる
	for _, server := range network.Servers {
		authenticationMethod, err := parseAuthenticationMethod(server)
		if err != nil {
			return err
		}

		err = remote.Copy(
			ctx,
			server.Host,
			zipPath,
			dst,
			authenticationMethod,
		)
		if err != nil {
			return err
		}

		client, session, err := newSession(
			server.Host,
			server.Environments,
			authenticationMethod,
		)
		if err != nil {
			return err
		}
		stdoutPipe, err := session.StdoutPipe()
		if err != nil {
			return err
		}
		stderrPipe, err := session.StderrPipe()
		if err != nil {
			return err
		}

		stdoutFile, err := os.Create(path.Join(getOutputPath(), fmt.Sprintf("%s:stdout", server.Host)))
		if err != nil {
			return err
		}
		stderrFile, err := os.Create(path.Join(getOutputPath(), fmt.Sprintf("./%s:stderr", server.Host)))
		if err != nil {
			return err
		}

		processes = append(processes, process{
			client:     client,
			host:       server.Host,
			session:    session,
			stdout:     stdoutPipe,
			stderr:     stderrPipe,
			stdoutFile: stdoutFile,
			stderrFile: stderrFile,
		})

		go func() {
			io.Copy(stdoutFile, stdoutPipe)
		}()

		go func() {
			io.Copy(stderrFile, stderrPipe)
		}()

		go func() {
			session.Run(command)
		}()
	}

	for {
		println("🤖    終了しますか？")
		println("👉    y/n")

		var in string
		fmt.Scan(&in)

		if in == "y" {
			break
		}
	}

	for _, process := range processes {
		process.session.Signal(ssh.SIGINT)
		process.stdoutFile.Close()
		process.stderrFile.Close()
		process.session.Close()
		process.client.Close()
	}

	err = tryToUnlock(
		lock,
		network,
	)
	if err != nil {
		return err
	}

	return nil
}

func tryToLock(lock string, network remote.Network) error {
	for _, server := range network.Servers {
		authenticationMethod, err := parseAuthenticationMethod(server)
		// TODO Unlockちゃんとやる
		if err != nil {
			return err
		}
		err = remote.Lock(
			lock,
			server.Host,
			authenticationMethod,
		)
		// TODO Unlockちゃんとやる
		if err != nil {
			return err
		}
	}

	return nil
}

func tryToUnlock(lock string, network remote.Network) error {
	for _, server := range network.Servers {
		authenticationMethod, err := parseAuthenticationMethod(server)
		// TODO Unlockちゃんとやる
		if err != nil {
			return err
		}
		err = remote.Unlock(
			lock,
			server.Host,
			authenticationMethod,
		)
		// TODO Unlockちゃんとやる
		if err != nil {
			return err
		}
	}

	return nil
}

func newSession(host string, environments []remote.Environment, authenticationMethod remote.AuthenticationMethod) (*ssh.Client, *ssh.Session, error) {
	client, err := remote.NewClient(host, authenticationMethod)
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

func parseAuthenticationMethod(server remote.Server) (remote.AuthenticationMethod, error) {
	var authenticationMethod remote.AuthenticationMethod
	switch server.Authentication {
	case remote.AUTHENTICATION_METHOD_PASSWORD:
		authenticationMethod = remote.PasswordAuthentication{
			User:     server.SSH.User,
			Password: server.SSH.Password,
		}

	case remote.AUTHENTICATION_METHOD_KEY:
		// TODO
		break

	default:
		return nil, errors.New(fmt.Sprintf(
			"authentication should be followings: %s, %s",
			remote.AUTHENTICATION_METHOD_PASSWORD,
			remote.AUTHENTICATION_METHOD_KEY,
		))
	}

	return authenticationMethod, nil
}
