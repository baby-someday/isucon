package cmd

import (
	"github.com/baby-someday/isucon/internal/metricscpu"
	"github.com/baby-someday/isucon/pkg/interaction"
	"github.com/baby-someday/isucon/pkg/remote"
	"github.com/baby-someday/isucon/pkg/servermaster"
	"github.com/baby-someday/isucon/pkg/util"
	"github.com/spf13/cobra"
)

var metricsCPUCmd = &cobra.Command{
	Use:   "cpu",
	Short: "cpu",
	Long:  `cpu`,
	Run:   runMetricsCPUCommand,
}

const (
	FLAG_METRICS_CPU_INTERVAL = "interval"

	FLAG_METRICS_CPU_INTERVAL_DEFAULT = 1
)

func init() {
	metricsCPUCmd.Flags().Int(
		FLAG_METRICS_CPU_INTERVAL,
		FLAG_METRICS_CPU_INTERVAL_DEFAULT,
		"",
	)
	metricsCPUCmd.Flags().String(
		FLAG_NETWORK_PATH,
		FLAG_NETWORK_PATH_DEFAULT,
		"",
	)
	metricsCPUCmd.Flags().String(
		FLAG_SERVER_PATH,
		FLAG_SERVER_PATH_DEFAULT,
		"",
	)

	metricsCmd.AddCommand(metricsCPUCmd)
}

func runMetricsCPUCommand(cmd *cobra.Command, args []string) {
	interval, err := cmd.Flags().GetInt(
		FLAG_METRICS_CPU_INTERVAL,
	)
	if err != nil {
		interaction.Error(err.Error())
		return
	}

	serverMasterConfig := servermaster.Config{}
	err = util.ParseFlag(cmd, FLAG_SERVER_PATH, &serverMasterConfig)
	if err != nil {
		interaction.Error(err.Error())
		return
	}

	network := remote.Network{}
	err = util.ParseFlag(cmd, FLAG_NETWORK_PATH, &network)
	if err != nil {
		interaction.Error(err.Error())
		return
	}

	err = metricscpu.MeasureMetrics(
		interval,
		serverMasterConfig.Servers,
		network.Servers,
	)
	if err != nil {
		interaction.Error(err.Error())
		return
	}
}
