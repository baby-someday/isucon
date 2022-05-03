package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/baby-someday/isucon/internal/distribute"
	"github.com/baby-someday/isucon/pkg/github"
	"github.com/baby-someday/isucon/pkg/me"
	"github.com/baby-someday/isucon/pkg/project"
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

const (
	FLAG_DISTRIBUTE_FROM = "from"

	FLAG_DISTRIBUTE_FROM_DEFAULT = distribute.FROM_LOCAL
)

func init() {
	distributeCmd.Flags().String(
		FLAG_CONFIG_PATH,
		FLAG_CONFIG_PATH_DEFAULT,
		"",
	)
	distributeCmd.Flags().String(
		FLAG_GITHUB_PATH,
		FLAG_GITHUB_PATH_DEFAULT,
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
		FLAG_PROJECT_PATH,
		FLAG_PROJECT_PATH_DEFAULT,
		"",
	)
	distributeCmd.Flags().String(
		FLAG_SLACK_PATH,
		FLAG_SLACK_PATH_DEFAULT,
		"",
	)
	distributeCmd.Flags().String(
		FLAG_DISTRIBUTE_FROM,
		FLAG_DISTRIBUTE_FROM_DEFAULT,
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

	from, err := cmd.Flags().GetString(FLAG_DISTRIBUTE_FROM)
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

	switch from {
	case distribute.FROM_LOCAL:
		distributeFromLocal(
			cmd,
			context.Background(),
			config,
			network,
		)

	case distribute.FROM_GIT_HUB:
		err = distributeFromGitHub(
			cmd,
			context.Background(),
			config,
			network,
		)
	}

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

func distributeFromLocal(
	cmd *cobra.Command,
	ctx context.Context,
	config distribute.Config,
	network remote.Network,
) error {
	project := project.Project{}
	err := util.ParseFlag(
		cmd,
		FLAG_PROJECT_PATH,
		&project,
	)
	if err != nil {
		return err
	}
	return distribute.DistributeFromLocal(
		context.Background(),
		network,
		project.Src,
		config.Dst,
		config.Lock,
		config.Command,
		config.Ignore,
	)
}

func distributeFromGitHub(
	cmd *cobra.Command,
	ctx context.Context,
	config distribute.Config,
	network remote.Network,
) error {
	github := github.GitHub{}
	err := util.ParseFlag(
		cmd,
		FLAG_GITHUB_PATH,
		&github,
	)
	if err != nil {
		return err
	}

	println("🤖    どのブランチをデプロイしますか？")
	print("👉    ")
	for index, branch := range github.Repository.Branches {
		print(fmt.Sprintf("%d:%s    ", index, branch))
	}
	println()

	var index int
	fmt.Scan(&index)

	if len(github.Repository.Branches) <= index {
		return errors.New("bad index")
	}

	return distribute.DistributeFromGitHub(
		ctx,
		network,
		github.Token,
		github.Repository.Owner,
		github.Repository.Name,
		github.Repository.URL,
		github.Repository.Branches[index],
		config.Dst,
		config.Lock,
		config.Command,
		config.Ignore,
	)
}
