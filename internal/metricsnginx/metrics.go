package metricsnginx

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/baby-someday/isucon/pkg/interaction"
	"github.com/baby-someday/isucon/pkg/nginx"
	"github.com/baby-someday/isucon/pkg/output"
	"github.com/baby-someday/isucon/pkg/remote"
	"github.com/baby-someday/isucon/pkg/util"
)

func CopyFiles(servers []nginx.Server) error {
	accessLogFilePaths := make([]string, len(servers))

	for index, server := range servers {
		interaction.Message(fmt.Sprintf("%sの処理を開始します。", server.Host))
		authenticationMethod, err := remote.MakeAuthenticationMethod(server.SSH)
		if err != nil {
			return util.HandleError(err)
		}

		interaction.Message("NGINXログファイルのコピーを開始します。")
		accessLogFilePath, err := nginx.CopyLogFiles(
			output.GetNginxMetricsDirPath(),
			server.Host,
			server.Log.Access,
			server.Log.Error,
			authenticationMethod,
		)
		if err != nil {
			interaction.Error("NGINXログファイルのコピーに失敗しました。")
			return util.HandleError(err)
		}
		interaction.Message("NGINXログファイルのコピーが完了しました。")

		accessLogFilePaths[index] = accessLogFilePath

		interaction.Message("NGINXアクセスログの入れ替えを開始します。")
		err = nginx.RotateLogFile(
			server.Host,
			server.Log.Access,
			authenticationMethod,
		)
		if err != nil {
			interaction.Error("NGINXアクセスログの入れ替えに失敗しました。")
			return util.HandleError(err)
		}
		interaction.Message("NGINXアクセスログの入れ替えが完了しました。")

		defer func() {
			interaction.Message("NGINXのリスタートを開始します。")
			err := nginx.Restart(
				server.Host,
				server.Bin,
				authenticationMethod,
			)
			if err != nil {
				interaction.Error("NGINXのリスタートに失敗しました")
				util.HandleError(err)
				return
			}
			interaction.Message("NGINXのリスタートが完了しました。")
		}()

		interaction.Message("NGINXエラーログの入れ替えを開始します。")
		err = nginx.RotateLogFile(
			server.Host,
			server.Log.Error,
			authenticationMethod,
		)
		if err != nil {
			interaction.Error("NGINXエラーログの入れ替えに失敗しました。")
			return util.HandleError(err)
		}
		interaction.Message("NGINXエラーログの入れ替えが完了しました。")
	}

	interaction.Message("NGINXアクセスログの統合を開始します。")

	accessLogFilePath := path.Join(output.GetNginxMetricsDirPath(), "access.log")
	err := os.MkdirAll(path.Dir(accessLogFilePath), 0755)
	if err != nil {
		return util.HandleError(err)
	}
	accessLogFile, err := os.Create(accessLogFilePath)
	if err != nil {
		return util.HandleError(err)
	}
	if err != nil {
		interaction.Error("NGINXアクセスログの統合に失敗しました。")
		return util.HandleError(err)
	}
	defer accessLogFile.Close()
	for _, accessLogFilePath := range accessLogFilePaths {
		bytes, err := ioutil.ReadFile(accessLogFilePath)
		if err != nil {
			interaction.Error("NGINXアクセスログの統合に失敗しました。")
			return util.HandleError(err)
		}
		accessLogFile.Write(bytes)
	}
	interaction.Message("NGINXアクセスログの統合が完了しました。")

	return nil
}
