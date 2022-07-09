package cmd

import (
	"github.com/spf13/cobra"
)

var analysisCmd = &cobra.Command{
	Use:   "analysis",
	Short: "analysis",
	Long:  `analysis`,
}

func init() {
	rootCmd.AddCommand(analysisCmd)
}
