package nginx

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/baby-someday/isucon/pkg/remote"
	"github.com/baby-someday/isucon/pkg/util"
)

func RotateLogFile(host, logFilePath string, authenticationMethod remote.AuthenticationMethod) error {
	_, err := remote.Exec(
		host,
		fmt.Sprintf("echo \"\" > %s", logFilePath),
		make([]remote.Environment, 0),
		authenticationMethod,
	)
	if err != nil {
		return util.HandleError(err)
	}
	return nil
}

func Restart(host, nginxBin string, authenticationMethod remote.AuthenticationMethod) error {
	_, err := remote.Exec(
		host,
		fmt.Sprintf("sudo %s -s reopen", nginxBin),
		make([]remote.Environment, 0),
		authenticationMethod,
	)
	if err != nil {
		return util.HandleError(err)
	}
	return nil
}

func CopyLogFiles(
	outputDirPath,
	host,
	remoteAccessLogPath,
	remoteErrorLogPath string,
	authenticationMethod remote.AuthenticationMethod,
) (string, error) {
	now := time.Now()
	timestamp := now.Format("2006-01-02_15:04:05")
	accessLogFilePath := path.Join(outputDirPath, host, fmt.Sprintf("access_%s.log", timestamp))
	err := copyLogFile(
		host,
		accessLogFilePath,
		remoteAccessLogPath,
		authenticationMethod,
	)
	if err != nil {
		return "", util.HandleError(err)
	}
	err = copyLogFile(
		host,
		path.Join(outputDirPath, host, fmt.Sprintf("error_%s.log", timestamp)),
		remoteErrorLogPath,
		authenticationMethod,
	)
	if err != nil {
		return "", util.HandleError(err)
	}

	return accessLogFilePath, nil
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
