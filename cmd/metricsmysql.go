package cmd

import (
	"log"

	"github.com/baby-someday/isucon/internal/metricsmysql"
	"github.com/baby-someday/isucon/pkg/interaction"
	"github.com/baby-someday/isucon/pkg/mysql"
	"github.com/baby-someday/isucon/pkg/servermaster"
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
	metricsMySQLCmd.Flags().String(
		FLAG_SERVER_PATH,
		FLAG_SERVER_PATH_DEFAULT,
		"",
	)

	metricsCmd.AddCommand(metricsMySQLCmd)
}

func runMetricsMySQLCommand(cmd *cobra.Command, args []string) {
	serverMasterConfig := servermaster.Config{}
	err := util.ParseFlag(cmd, FLAG_SERVER_PATH, &serverMasterConfig)
	if err != nil {
		interaction.Error(err.Error())
		return
	}

	mysqlConfig := mysql.Config{}
	err = util.ParseFlag(cmd, FLAG_MYSQL_PATH, &mysqlConfig)
	if err != nil {
		interaction.Error(err.Error())
		return
	}

	log.Println(mysqlConfig.TTT.Name)

	err = metricsmysql.CopyLogFiles(
		serverMasterConfig.Servers,
		mysqlConfig.Servers,
	)
	if err != nil {
		interaction.Error(err.Error())
		return
	}
}
