package mysql

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

func FlushLogs(host, mysqlAdminBin, defaultsFilePath string, authenticationMethod remote.AuthenticationMethod) error {
	_, err := remote.Exec(
		host,
		fmt.Sprintf("sudo %s --defaults-file=%s flush-logs", mysqlAdminBin, defaultsFilePath),
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
	remoteSlowQueryLogFilePath string,
	authenticationMethod remote.AuthenticationMethod,
) (string, error) {
	now := time.Now()
	timestamp := now.Format("2006-01-02_15:04:05")
	localSlowQueryLogFilePath := path.Join(outputDirPath, host, fmt.Sprintf("slow_%s.log", timestamp))
	err := copyLogFile(
		host,
		localSlowQueryLogFilePath,
		remoteSlowQueryLogFilePath,
		authenticationMethod,
	)
	if err != nil {
		return "", util.HandleError(err)
	}

	return localSlowQueryLogFilePath, nil
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
