package cmd

import (
	"github.com/baby-someday/isucon/internal/metricsnginx"
	"github.com/baby-someday/isucon/pkg/interaction"
	"github.com/baby-someday/isucon/pkg/nginx"
	"github.com/baby-someday/isucon/pkg/servermaster"
	"github.com/baby-someday/isucon/pkg/util"
	"github.com/spf13/cobra"
)

var metricsNginxCmd = &cobra.Command{
	Use:   "nginx",
	Short: "nginx",
	Long:  `nginx`,
	Run:   runMetricsNginxCommand,
}

func init() {
	metricsNginxCmd.Flags().String(
		FLAG_NGINX_PATH,
		FLAG_NGINX_PATH_DEFAULT,
		"",
	)
	metricsNginxCmd.Flags().String(
		FLAG_SERVER_PATH,
		FLAG_SERVER_PATH,
		"",
	)

	metricsCmd.AddCommand(metricsNginxCmd)
}

func runMetricsNginxCommand(cmd *cobra.Command, args []string) {
	serverMasterConfig := servermaster.Config{}
	err := util.ParseFlag(cmd, FLAG_SERVER_PATH, &serverMasterConfig)
	if err != nil {
		interaction.Error(err.Error())
		return
	}

	nginxConfig := nginx.Config{}
	err = util.ParseFlag(cmd, FLAG_NGINX_PATH, &nginxConfig)
	if err != nil {
		interaction.Error(err.Error())
		return
	}

	err = metricsnginx.CopyLogFiles(
		serverMasterConfig.Servers,
		nginxConfig.Servers,
	)
	if err != nil {
		interaction.Error(err.Error())
		return
	}
}
