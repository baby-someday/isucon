package nginx

import (
	"context"
	"os"
	"path"

	"github.com/baby-someday/isucon/pkg/remote"
)

func CopyLogFiles(host, localAccessLogFilePath, remoteAccessLogFilePath, localErrorLogFilePath, remoteErrorLogFilePath string, authenticationMethod remote.AuthenticationMethod) error {
	err := copyLogFile(
		host,
		localAccessLogFilePath,
		remoteAccessLogFilePath,
		authenticationMethod,
	)
	if err != nil {
		return err
	}

	return copyLogFile(
		host,
		localErrorLogFilePath,
		remoteErrorLogFilePath,
		authenticationMethod,
	)
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
