package cmd

import (
	"github.com/baby-someday/isucon/internal/metricsmysql"
	"github.com/baby-someday/isucon/pkg/interaction"
	"github.com/baby-someday/isucon/pkg/mysql"
	"github.com/baby-someday/isucon/pkg/util"
	"github.com/spf13/cobra"
)

var metricsMySQLCmd = &cobra.Command{
	Use:   "mysql",
	Short: "mysql",
	Long:  `mysql`,
	Run:   runMetricsMySQLCommand,
}

func init() {
	metricsMySQLCmd.Flags().String(
		FLAG_MYSQL_PATH,
		FLAG_MYSQL_PATH_DEFAULT,
		"",
	)

	metricsCmd.AddCommand(metricsMySQLCmd)
}

func runMetricsMySQLCommand(cmd *cobra.Command, args []string) {
	mysqlConfig := mysql.Config{}
	err := util.ParseFlag(cmd, FLAG_MYSQL_PATH, &mysqlConfig)
	if err != nil {
		interaction.Error(err.Error())
		return
	}

	err = metricsmysql.CopyLogFiles(mysqlConfig.Servers)
	if err != nil {
		interaction.Error(err.Error())
		return
	}
}
