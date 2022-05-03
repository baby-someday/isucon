package cmd

import (
	"github.com/spf13/cobra"
)

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "metrics",
	Long:  `metrics`,
}

func init() {
	rootCmd.AddCommand(metricsCmd)
}
