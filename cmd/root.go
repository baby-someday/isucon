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
	FLAG_CONFIG_PATH  = "config"
	FLAG_NETWORK_PATH = "network"
	FLAG_ME_PATH      = "me"
	FLAG_SLACK_PATH   = "slack"

	FLAG_CONFIG_PATH_DEFAULT  = "./config/distribute.yml"
	FLAG_ME_PATH_DEFAULT      = "./config/me.yml"
	FLAG_NETWORK_PATH_DEFAULT = "./config/network.yml"
	FLAG_SLACK_PATH_DEFAULT   = "./config/slack.yml"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
