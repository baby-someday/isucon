package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/baby-someday/isucon/internal/distribute"
	"github.com/baby-someday/isucon/pkg/github"
	"github.com/baby-someday/isucon/pkg/interaction"
	"github.com/baby-someday/isucon/pkg/me"
	"github.com/baby-someday/isucon/pkg/nginx"
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
		FLAG_DISTRIBUTE_PATH,
		FLAG_DISTRIBUTE_PATH_DEFAULT,
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
		FLAG_NGINX_PATH,
		FLAG_NGINX_PATH_DEFAULT,
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
	err := util.ParseFlag(cmd, FLAG_DISTRIBUTE_PATH, &config)
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

	slackConfig := slack.Slack{}
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
		fmt.Sprintf("%s\n🚀  %sさんがベンチマークを開始しました  🚀", slack.SEPARATOR, me.Name),
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
		fmt.Sprintf("💨  %sさんがベンチマークを終了しました  💨\n%s", me.Name, slack.SEPARATOR),
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
	interaction.Message("ローカルからデプロイを開始します。")

	project := project.Project{}
	err := util.ParseFlag(
		cmd,
		FLAG_PROJECT_PATH,
		&project,
	)
	if err != nil {
		return util.HandleError(err)
	}

	nginxConfig := nginx.Config{}
	err = util.ParseFlag(
		cmd,
		FLAG_NGINX_PATH,
		&nginxConfig,
	)
	if err != nil {
		return util.HandleError(err)
	}

	err = distribute.DistributeFromLocal(
		context.Background(),
		network,
		nginxConfig,
		project.Src,
		config.Dst,
		config.Lock,
		config.Command,
		config.Ignore,
	)
	if err != nil {
		return util.HandleError(err)
	}

	interaction.Message("ローカルのデプロイが完了しました。")
	return nil
}

func distributeFromGitHub(
	cmd *cobra.Command,
	ctx context.Context,
	config distribute.Config,
	network remote.Network,
) error {
	interaction.Message("GitHubからデプロイを開始します。")

	nginxConfig := nginx.Config{}
	err := util.ParseFlag(
		cmd,
		FLAG_NGINX_PATH,
		&nginxConfig,
	)
	if err != nil {
		return util.HandleError(err)
	}

	github := github.GitHub{}
	err = util.ParseFlag(
		cmd,
		FLAG_GITHUB_PATH,
		&github,
	)
	if err != nil {
		return util.HandleError(err)
	}

	slack := slack.Slack{}
	err = util.ParseFlag(
		cmd,
		FLAG_SLACK_PATH,
		&slack,
	)
	if err != nil {
		return util.HandleError(err)
	}

	indexString := interaction.Choose(
		"どのブランチをデプロイしますか？",
		len(github.Repository.Branches),
		func(index int) (string, string) {
			return strconv.Itoa(index), github.Repository.Branches[index]
		},
	)
	index, err := strconv.Atoi(indexString)
	if err != nil {
		return util.HandleError(err)
	}

	if len(github.Repository.Branches) <= index {
		return errors.New("bad index")
	}

	err = distribute.DistributeFromGitHub(
		ctx,
		network,
		nginxConfig,
		github.Token,
		github.Repository.Owner,
		github.Repository.Name,
		github.Repository.URL,
		github.Repository.Branches[index],
		slack.Token,
		slack.Channel,
		config.Dst,
		config.Lock,
		config.Command,
		config.Ignore,
	)
	if err != nil {
		return util.HandleError(err)
	}

	interaction.Message("GitHubのデプロイが完了しました。")
	return nil
}
