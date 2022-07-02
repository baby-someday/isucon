package nginx

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/baby-someday/isucon/pkg/remote"
	"github.com/baby-someday/isucon/pkg/util"
)

func RotateLogFile(host, logFilePath, persistenceLogFilePath string, authenticationMethod remote.AuthenticationMethod) error {
	_, err := remote.Exec(
		host,
		fmt.Sprintf("cat %s >> %s && echo \"\" > %s", logFilePath, persistenceLogFilePath, logFilePath),
		make([]remote.Environment, 0),
		authenticationMethod,
	)
	return util.HandleError(err)
}

func CopyLogFiles(
	outputDirPath,
	host,
	remoteAccessLogPath,
	remoteErrorLogPath,
	remotePersistenceAccessLogPath,
	remotePersistenceErrorLogPath string,
	authenticationMethod remote.AuthenticationMethod,
) error {
	err := copyLogFile(
		host,
		path.Join(outputDirPath, host, "access.log"),
		remoteAccessLogPath,
		authenticationMethod,
	)
	if err != nil {
		return util.HandleError(err)
	}
	err = copyLogFile(
		host,
		path.Join(outputDirPath, host, "error.log"),
		remoteErrorLogPath,
		authenticationMethod,
	)
	if err != nil {
		return util.HandleError(err)
	}
	err = copyLogFile(
		host,
		path.Join(outputDirPath, host, "access.all.log"),
		remotePersistenceAccessLogPath,
		authenticationMethod,
	)
	if err != nil {
		return util.HandleError(err)
	}
	err = copyLogFile(
		host,
		path.Join(outputDirPath, host, "error.all.log"),
		remotePersistenceErrorLogPath,
		authenticationMethod,
	)
	if err != nil {
		return util.HandleError(err)
	}

	return nil
}

func copyLogFile(host, localPath, remotePath string, authenticationMethod remote.AuthenticationMethod) error {
	err := os.MkdirAll(path.Dir(localPath), 0755)
	if err != nil {
		return util.HandleError(err)
	}
	file, err := os.Create(localPath)
	if err != nil {
		return util.HandleError(err)
	}
	return remote.CopyFromRemote(
		context.Background(),
		file,
		host,
		remotePath,
		authenticationMethod,
	)
}
