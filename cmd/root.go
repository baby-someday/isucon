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
	FLAG_DISTRIBUTE_PATH = "distribute"
	FLAG_GITHUB_PATH     = "github"
	FLAG_ME_PATH         = "me"
	FLAG_MYSQL_PATH      = "mysql"
	FLAG_NETWORK_PATH    = "network"
	FLAG_NGINX_PATH      = "nginx"
	FLAG_PROJECT_PATH    = "project"
	FLAG_SLACK_PATH      = "slack"

	FLAG_DISTRIBUTE_PATH_DEFAULT = "./config/distribute.yml"
	FLAG_GITHUB_PATH_DEFAULT     = "./config/github.yml"
	FLAG_ME_PATH_DEFAULT         = "./config/me.yml"
	FLAG_MYSQL_PATH_DEFAULT      = "./config/mysql.yml"
	FLAG_NETWORK_PATH_DEFAULT    = "./config/network.yml"
	FLAG_NGINX_PATH_DEFAULT      = "./config/nginx.yml"
	FLAG_PROJECT_PATH_DEFAULT    = "./config/project.yml"
	FLAG_SLACK_PATH_DEFAULT      = "./config/slack.yml"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
