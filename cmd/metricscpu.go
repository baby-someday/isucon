package cmd

import (
	"log"

	"github.com/baby-someday/isucon/internal/metricscpu"
	"github.com/baby-someday/isucon/pkg/remote"
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
	metricsCPUCmd.Flags().String(
		FLAG_NETWORK_PATH,
		FLAG_NETWORK_PATH_DEFAULT,
		"",
	)
	metricsCPUCmd.Flags().Int(
		FLAG_METRICS_CPU_INTERVAL,
		FLAG_METRICS_CPU_INTERVAL_DEFAULT,
		"",
	)

	metricsCmd.AddCommand(metricsCPUCmd)
}

func runMetricsCPUCommand(cmd *cobra.Command, args []string) {
	network := remote.Network{}
	err := util.ParseFlag(cmd, FLAG_NETWORK_PATH, &network)
	if err != nil {
		log.Fatal(err)
	}

	interval, err := cmd.Flags().GetInt(
		FLAG_METRICS_CPU_INTERVAL,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = metricscpu.MeasureMetrics(interval, network.Servers)
	if err != nil {
		log.Fatal(err)
	}
}
