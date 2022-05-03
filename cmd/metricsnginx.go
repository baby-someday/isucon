package cmd

import (
	"log"

	"github.com/baby-someday/isucon/internal/metricsnginx"
	"github.com/baby-someday/isucon/pkg/remote"
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
		FLAG_NETWORK_PATH,
		FLAG_NETWORK_PATH_DEFAULT,
		"",
	)

	metricsCmd.AddCommand(metricsNginxCmd)
}

func runMetricsNginxCommand(cmd *cobra.Command, args []string) {
	network := remote.Network{}
	err := util.ParseFlag(cmd, FLAG_NETWORK_PATH, &network)
	if err != nil {
		log.Fatal(err)
	}

	err = metricsnginx.CopyFiles(network)
	if err != nil {
		log.Fatal(err)
	}
}
