package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/baby-someday/isucon/internal/distribute"
	"github.com/baby-someday/isucon/pkg/me"
	"github.com/baby-someday/isucon/pkg/remote"
	"github.com/baby-someday/isucon/pkg/slack"
	"github.com/baby-someday/isucon/pkg/util"
	"github.com/spf13/cobra"
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
		FLAG_CONFIG_PATH_DEFAULT,
		"",
	)
	distributeCmd.Flags().String(
		FLAG_ME_PATH,
		FLAG_ME_PATH_DEFAULT,
		"",
	)
	distributeCmd.Flags().String(
		FLAG_NETWORK_PATH,
		FLAG_NETWORK_PATH_DEFAULT,
		"",
	)
	distributeCmd.Flags().String(
		FLAG_SLACK_PATH,
		FLAG_SLACK_PATH_DEFAULT,
		"",
	)

	rootCmd.AddCommand(distributeCmd)
}

func runDistributeCommand(cmd *cobra.Command, args []string) {
	config := distribute.Config{}
	err := util.ParseFlag(cmd, FLAG_CONFIG_PATH, &config)
	if err != nil {
		log.Fatal(err)
	}

	me := me.Config{}
	err = util.ParseFlag(cmd, FLAG_ME_PATH, &me)
	if err != nil {
		log.Fatal(err)
	}

	network := remote.Network{}
	err = util.ParseFlag(cmd, FLAG_NETWORK_PATH, &network)
	if err != nil {
		log.Fatal(err)
	}

	slackConfig := slack.Config{}
	err = util.ParseFlag(cmd, FLAG_SLACK_PATH, &slackConfig)
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
