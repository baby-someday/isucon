package cmd

import (
	"github.com/baby-someday/isucon/internal/metricsnginx"
	"github.com/baby-someday/isucon/pkg/interaction"
	"github.com/baby-someday/isucon/pkg/nginx"
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

	metricsCmd.AddCommand(metricsNginxCmd)
}

func runMetricsNginxCommand(cmd *cobra.Command, args []string) {
	nginxConfig := nginx.Config{}
	err := util.ParseFlag(cmd, FLAG_NETWORK_PATH, &nginxConfig)
	if err != nil {
		interaction.Error(err.Error())
		return
	}

	err = metricsnginx.CopyFiles(nginxConfig.Servers)
	if err != nil {
		interaction.Error(err.Error())
		return
	}
}
