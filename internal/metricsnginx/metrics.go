package metricsnginx

import (
	"fmt"

	"github.com/baby-someday/isucon/pkg/interaction"
	"github.com/baby-someday/isucon/pkg/nginx"
	"github.com/baby-someday/isucon/pkg/output"
	"github.com/baby-someday/isucon/pkg/remote"
	"github.com/baby-someday/isucon/pkg/util"
)

func CopyFiles(servers []remote.Server) error {
	for _, server := range servers {
		interaction.Message(fmt.Sprintf("%sの処理を開始します。", server.Host))
		authenticationMethod, err := remote.MakeAuthenticationMethod(server)
		if err != nil {
			return util.HandleError(err)
		}

		interaction.Message("NGINXログファイルのコピーを開始します。")
		err = nginx.CopyLogFiles(
			output.GetNginxMetricsDirPath(),
			server.Host,
			server.Nginx.Log.Access,
			server.Nginx.Log.Error,
			authenticationMethod,
		)
		if err != nil {
			interaction.Error("NGINXログファイルのコピーに失敗しました。")
			return util.HandleError(err)
		}
		interaction.Message("NGINXログファイルのコピーが完了しました。")

		interaction.Message("NGINXアクセスログの入れ替えを開始します。")
		err = nginx.RotateLogFile(
			server.Host,
			server.Nginx.Log.Access,
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
				server.Nginx.Bin,
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
			server.Nginx.Log.Error,
			authenticationMethod,
		)
		if err != nil {
			interaction.Error("NGINXエラーログの入れ替えに失敗しました。")
			return util.HandleError(err)
		}
		interaction.Message("NGINXエラーログの入れ替えが完了しました。")
	}

	return nil
}
