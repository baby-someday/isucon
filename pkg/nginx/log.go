package nginx

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/baby-someday/isucon/pkg/remote"
)

func RotateLogFile(host, logFilePath, persistenceLogFilePath string, authenticationMethod remote.AuthenticationMethod) error {
	_, err := remote.Exec(
		host,
		fmt.Sprintf("cat %s >> %s && echo \"\" > %s", logFilePath, persistenceLogFilePath, logFilePath),
		make([]remote.Environment, 0),
		authenticationMethod,
	)
	return err
}

func CopyLogFiles(outputDirPath string, network remote.Network) error {
	for _, server := range network.Servers {
		authenticationMethod, err := remote.MakeAuthenticationMethod(server)
		if err != nil {
			return err
		}
		err = copyLogFile(
			server.Host,
			path.Join(outputDirPath, server.Host, "access.log"),
			server.Nginx.Log.Access,
			authenticationMethod,
		)
		if err != nil {
			return err
		}
		err = copyLogFile(
			server.Host,
			path.Join(outputDirPath, server.Host, "error.log"),
			server.Nginx.Log.Error,
			authenticationMethod,
		)
		if err != nil {
			return err
		}
		err = copyLogFile(
			server.Host,
			path.Join(outputDirPath, server.Host, "access.all.log"),
			server.Nginx.Log.Persistence.Access,
			authenticationMethod,
		)
		if err != nil {
			return err
		}
		err = copyLogFile(
			server.Host,
			path.Join(outputDirPath, server.Host, "error.all.log"),
			server.Nginx.Log.Persistence.Error,
			authenticationMethod,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyLogFile(host, localPath, remotePath string, authenticationMethod remote.AuthenticationMethod) error {
	err := os.MkdirAll(path.Dir(localPath), 0755)
	if err != nil {
		return err
	}
	file, err := os.Create(localPath)
	if err != nil {
		return err
	}
	return remote.CopyFromRemote(
		context.Background(),
		file,
		host,
		remotePath,
		authenticationMethod,
	)
}
