package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/baby-someday/isucon/internal/distribute"
	"github.com/baby-someday/isucon/pkg/me"
	"github.com/baby-someday/isucon/pkg/remote"
	"github.com/baby-someday/isucon/pkg/slack"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "metrics",
	Long:  `metrics`,
}

func init() {
	rootCmd.AddCommand(metricsCmd)
}

func runMetricsCommand(cmd *cobra.Command, args []string) {
	configFilePath, err := cmd.Flags().GetString(FLAG_CONFIG_PATH)
	if err != nil {
		log.Fatal(err)
	}

	configFileBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	config := distribute.Config{}
	err = yaml.Unmarshal(configFileBytes, &config)
	if err != nil {
		log.Fatal(err)
	}

	meFilePath, err := cmd.Flags().GetString(FLAG_ME_PATH)
	if err != nil {
		log.Fatal(err)
	}

	meFileBytes, err := ioutil.ReadFile(meFilePath)
	if err != nil {
		log.Fatal(err)
	}

	me := me.Config{}
	err = yaml.Unmarshal(meFileBytes, &me)
	if err != nil {
		log.Fatal(err)
	}

	networkFilePath, err := cmd.Flags().GetString(FLAG_NETWORK_PATH)
	if err != nil {
		log.Fatal(err)
	}

	networkFileBytes, err := ioutil.ReadFile(networkFilePath)
	if err != nil {
		log.Fatal(err)
	}

	network := remote.Network{}
	err = yaml.Unmarshal(networkFileBytes, &network)
	if err != nil {
		log.Fatal(err)
	}

	slackFilePath, err := cmd.Flags().GetString(FLAG_SLACK_PATH)
	if err != nil {
		log.Fatal(err)
	}

	slackFileBytes, err := ioutil.ReadFile(slackFilePath)
	if err != nil {
		log.Fatal(err)
	}

	slackConfig := slack.Config{}
	err = yaml.Unmarshal(slackFileBytes, &slackConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = slack.PostMessage(
		slackConfig.Token,
		slackConfig.Channel,
		fmt.Sprintf("🚀 %sさんがベンチマークを開始しました", me.Name),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = distribute.Distribute(
		context.Background(),
		network,
		config.Src,
		config.Dst,
		config.Lock,
		config.Command,
		config.Ignore,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = slack.PostMessage(
		slackConfig.Token,
		slackConfig.Channel,
		fmt.Sprintf("💨 %sさんがベンチマークを終了しました", me.Name),
	)
	if err != nil {
		log.Fatal(err)
	}
}
