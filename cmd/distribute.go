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

var distributeCmd = &cobra.Command{
	Use:   "distribute",
	Short: "distribute",
	Long:  `distribute`,
	Run:   runDistributeCommand,
}

func init() {
	distributeCmd.Flags().String(
		FLAG_CONFIG_PATH,
		"./config/distribute.yml",
		"",
	)
	distributeCmd.Flags().String(
		FLAG_ME_PATH,
		"./config/me.yml",
		"",
	)
	distributeCmd.Flags().String(
		FLAG_NETWORK_PATH,
		"./config/network.yml",
		"",
	)
	distributeCmd.Flags().String(
		FLAG_SLACK_PATH,
		"./config/slack.yml",
		"",
	)
	rootCmd.AddCommand(distributeCmd)
}

func runDistributeCommand(cmd *cobra.Command, args []string) {
	configFilePath, err := cmd.Flags().GetString(FLAG_CONFIG_PATH)
	if err != nil {
		log.Fatal(err.Error())
	}

	configFileBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	config := distribute.Config{}
	err = yaml.Unmarshal(configFileBytes, &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	meFilePath, err := cmd.Flags().GetString(FLAG_ME_PATH)
	if err != nil {
		log.Fatal(err.Error())
	}

	meFileBytes, err := ioutil.ReadFile(meFilePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	me := me.Config{}
	err = yaml.Unmarshal(meFileBytes, &me)
	if err != nil {
		log.Fatal(err.Error())
	}

	networkFilePath, err := cmd.Flags().GetString(FLAG_NETWORK_PATH)
	if err != nil {
		log.Fatal(err.Error())
	}

	networkFileBytes, err := ioutil.ReadFile(networkFilePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	network := remote.Network{}
	err = yaml.Unmarshal(networkFileBytes, &network)
	if err != nil {
		log.Fatal(err.Error())
	}

	slackFilePath, err := cmd.Flags().GetString(FLAG_SLACK_PATH)
	if err != nil {
		log.Fatal(err.Error())
	}

	slackFileBytes, err := ioutil.ReadFile(slackFilePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	slackConfig := slack.Config{}
	err = yaml.Unmarshal(slackFileBytes, &slackConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = slack.Post(
		slackConfig.Token,
		slackConfig.Channel,
		fmt.Sprintf("🚀 %sさんがベンチマークを開始しました", me.Name),
	)
	if err != nil {
		log.Fatal(err.Error())
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
		log.Fatal(err.Error())
	}

	err = slack.Post(
		slackConfig.Token,
		slackConfig.Channel,
		fmt.Sprintf("💨 %sさんがベンチマークを終了しました", me.Name),
	)
	if err != nil {
		log.Fatal(err.Error())
	}
}
