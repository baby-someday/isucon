package cmd

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/baby-someday/isucon/internal/distribute"
	"github.com/baby-someday/isucon/pkg/github"
	"github.com/baby-someday/isucon/pkg/interaction"
	"github.com/baby-someday/isucon/pkg/me"
	"github.com/baby-someday/isucon/pkg/mysql"
	"github.com/baby-someday/isucon/pkg/nginx"
	"github.com/baby-someday/isucon/pkg/project"
	"github.com/baby-someday/isucon/pkg/remote"
	"github.com/baby-someday/isucon/pkg/servermaster"
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
		FLAG_MYSQL_PATH,
		FLAG_MYSQL_PATH_DEFAULT,
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
		FLAG_SERVER_PATH,
		FLAG_SERVER_PATH_DEFAULT,
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
		interaction.Error(err.Error())
		return
	}

	me := me.Config{}
	err = util.ParseFlag(cmd, FLAG_ME_PATH, &me)
	if err != nil {
		interaction.Error(err.Error())
		return
	}

	serverMasterConfig := servermaster.Config{}
	err = util.ParseFlag(cmd, FLAG_SERVER_PATH, &serverMasterConfig)
	if err != nil {
		interaction.Error(err.Error())
		return
	}

	network := remote.Network{}
	err = util.ParseFlag(cmd, FLAG_NETWORK_PATH, &network)
	if err != nil {
		interaction.Error(err.Error())
		return
	}

	nginxConfig := nginx.Config{}
	err = util.ParseFlag(
		cmd,
		FLAG_NGINX_PATH,
		&nginxConfig,
	)
	if err != nil {
		interaction.Error(err.Error())
		return
	}

	mysqlConfig := mysql.Config{}
	err = util.ParseFlag(
		cmd,
		FLAG_MYSQL_PATH,
		&mysqlConfig,
	)
	if err != nil {
		interaction.Error(err.Error())
		return
	}

	slackConfig := slack.Slack{}
	err = util.ParseFlag(cmd, FLAG_SLACK_PATH, &slackConfig)
	if err != nil {
		interaction.Error(err.Error())
		return
	}

	from, err := cmd.Flags().GetString(FLAG_DISTRIBUTE_FROM)
	if err != nil {
		interaction.Error(err.Error())
		return
	}

	err = slack.PostMessage(
		slackConfig.Token,
		slackConfig.Channel,
		fmt.Sprintf("%s\n🚀  %sさんがベンチマークを開始しました  🚀", slack.SEPARATOR, me.Name),
	)
	if err != nil {
		interaction.Error(err.Error())
		return
	}

	switch from {
	case distribute.FROM_LOCAL:
		distributeFromLocal(
			cmd,
			context.Background(),
			config,
			serverMasterConfig,
			network,
			nginxConfig,
			mysqlConfig,
		)

	case distribute.FROM_GIT_HUB:
		err = distributeFromGitHub(
			cmd,
			context.Background(),
			config,
			serverMasterConfig,
			network,
			nginxConfig,
			mysqlConfig,
		)
	}

	if err != nil {
		interaction.Error(err.Error())
		return
	}

	err = slack.PostMessage(
		slackConfig.Token,
		slackConfig.Channel,
		fmt.Sprintf("💨  %sさんがベンチマークを終了しました  💨\n%s", me.Name, slack.SEPARATOR),
	)
	if err != nil {
		interaction.Error(err.Error())
		return
	}
}

func distributeFromLocal(
	cmd *cobra.Command,
	ctx context.Context,
	config distribute.Config,
	serverMasterConfig servermaster.Config,
	network remote.Network,
	nginxConfig nginx.Config,
	mysqlConfig mysql.Config,
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

	err = distribute.DistributeFromLocal(
		context.Background(),
		serverMasterConfig.Servers,
		network.Servers,
		nginxConfig,
		mysqlConfig,
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
	serverMasterConfig servermaster.Config,
	network remote.Network,
	nginxConfig nginx.Config,
	mysqlConfig mysql.Config,
) error {
	interaction.Message("GitHubからデプロイを開始します。")

	github := github.GitHub{}
	err := util.ParseFlag(
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
		serverMasterConfig.Servers,
		network.Servers,
		nginxConfig,
		mysqlConfig,
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
