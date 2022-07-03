package metricsmysql

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/baby-someday/isucon/pkg/interaction"
	"github.com/baby-someday/isucon/pkg/mysql"
	"github.com/baby-someday/isucon/pkg/output"
	"github.com/baby-someday/isucon/pkg/remote"
	"github.com/baby-someday/isucon/pkg/servermaster"
	"github.com/baby-someday/isucon/pkg/util"
)

func CopyLogFiles(
	serverMasters []servermaster.ServerMaster,
	servers []mysql.Server,
) error {
	sloqQueryLogFilePaths := make([]string, len(servers))

	for index, server := range servers {
		serverMaster, err := servermaster.FindServerMaster(
			server.Name,
			serverMasters,
		)

		interaction.Message(fmt.Sprintf("%sの処理を開始します。", serverMaster.Host))
		authenticationMethod, err := remote.MakeAuthenticationMethod(serverMaster.SSH)
		if err != nil {
			return util.HandleError(err)
		}

		interaction.Message("MySQLログファイルのコピーを開始します。")
		slowQueryLogFilePath, err := mysql.CopyLogFiles(
			output.GetMySQLMetricsDirPath(),
			serverMaster.Host,
			server.Log.Slow,
			authenticationMethod,
		)
		if err != nil {
			interaction.Error("MySQLログファイルのコピーに失敗しました。")
			return util.HandleError(err)
		}
		interaction.Message("MySQLログファイルのコピーが完了しました。")

		sloqQueryLogFilePaths[index] = slowQueryLogFilePath

		interaction.Message("MySQLアクセスログの入れ替えを開始します。")
		err = mysql.RotateLogFile(
			serverMaster.Host,
			server.Log.Slow,
			authenticationMethod,
		)
		if err != nil {
			interaction.Error("MySQLアクセスログの入れ替えに失敗しました。")
			return util.HandleError(err)
		}
		interaction.Message("MySQLアクセスログの入れ替えが完了しました。")

		defer func() {
			interaction.Message("MySQLのflush-logsを開始します。")
			err := mysql.FlushLogs(
				serverMaster.Host,
				server.Bin.MySQLAdmin,
				server.Defaults,
				authenticationMethod,
			)
			if err != nil {
				interaction.Error("MySQLのflush-logsに失敗しました")
				util.HandleError(err)
				interaction.Error(err.Error())
				return
			}
			interaction.Message("MySQLのflush-logsが完了しました。")
		}()

		interaction.Message("MySQLログの入れ替えを開始します。")
		err = mysql.RotateLogFile(
			serverMaster.Host,
			server.Log.Slow,
			authenticationMethod,
		)
		if err != nil {
			interaction.Error("MySQLログの入れ替えに失敗しました。")
			return util.HandleError(err)
		}
		interaction.Message("MySQLログの入れ替えが完了しました。")
	}

	interaction.Message("MySQLログの統合を開始します。")

	slowQueryLogFilePath := path.Join(output.GetMySQLMetricsDirPath(), "slow.log")
	err := os.MkdirAll(path.Dir(slowQueryLogFilePath), 0755)
	if err != nil {
		return util.HandleError(err)
	}
	slowQueryLogFile, err := os.Create(slowQueryLogFilePath)
	if err != nil {
		return util.HandleError(err)
	}
	if err != nil {
		interaction.Error("MySQLログの統合に失敗しました。")
		return util.HandleError(err)
	}
	defer slowQueryLogFile.Close()
	for _, slowQueryLogFilePath := range sloqQueryLogFilePaths {
		bytes, err := ioutil.ReadFile(slowQueryLogFilePath)
		if err != nil {
			interaction.Error("MySQLログの統合に失敗しました。")
			return util.HandleError(err)
		}
		slowQueryLogFile.Write(bytes)
	}
	interaction.Message("MySQLログの統合が完了しました。")

	return nil
}
