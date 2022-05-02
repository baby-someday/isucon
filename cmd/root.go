package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "isucon",
	Short: "isucon",
	Long:  `isucon`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

const (
	FLAG_CONFIG_PATH = "config"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
