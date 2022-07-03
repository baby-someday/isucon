package distribute

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"

	"github.com/baby-someday/isucon/internal/metricscpu"
	"github.com/baby-someday/isucon/internal/metricsmysql"
	"github.com/baby-someday/isucon/internal/metricsnginx"
	"github.com/baby-someday/isucon/pkg/build"
	"github.com/baby-someday/isucon/pkg/github"
	"github.com/baby-someday/isucon/pkg/interaction"
	"github.com/baby-someday/isucon/pkg/mysql"
	"github.com/baby-someday/isucon/pkg/nginx"
	"github.com/baby-someday/isucon/pkg/output"
	"github.com/baby-someday/isucon/pkg/remote"
	"github.com/baby-someday/isucon/pkg/slack"
	"github.com/baby-someday/isucon/pkg/util"
	"golang.org/x/crypto/ssh"
)

type process struct {
	host       string
	client     *ssh.Client
	session    *ssh.Session
	stdout     io.Reader
	stderr     io.Reader
	stdoutFile *os.File
	stderrFile *os.File
}

type action struct {
	name     string
	callback func() error
}

const (
	FROM_LOCAL   = "local"
	FROM_GIT_HUB = "github"
)

func DistributeFromLocal(
	ctx context.Context,
	network remote.Network,
	nginxConfig nginx.Config,
	mysqlConfig mysql.Config,
	src,
	dst,
	lock,
	command string,
	ignore []string,
) error {
	return distribute(
		ctx,
		network,
		dst,
		lock,
		command,
		ignore,
		[]action{
			makeCPUMetricsAction(network.Servers),
			makeNginxMetricsAction(nginxConfig.Servers),
			makeMySQLMetricsAction(mysqlConfig.Servers),
		},
		deloyFromLocal(
			ctx,
			network,
			src,
			dst,
			ignore,
		),
	)
}

func DistributeFromGitHub(
	ctx context.Context,
	network remote.Network,
	nginxConfig nginx.Config,
	mysqlConfig mysql.Config,
	githubToken,
	repositoryOwner,
	repositoryName,
	repositoryURL,
	repositoryBranch,
	slackToken,
	slcakChannel,
	dst,
	lock,
	command string,
	ignore []string,
) error {
	err := distribute(
		ctx,
		network,
		dst,
		lock,
		command,
		ignore,
		[]action{
			makeCPUMetricsAction(network.Servers),
			makeNginxMetricsAction(nginxConfig.Servers),
			makeMySQLMetricsAction(mysqlConfig.Servers),
			makeSaveScoreAction(
				githubToken,
				repositoryOwner,
				repositoryName,
				repositoryBranch,
				slackToken,
				slcakChannel,
			),
		},
		deloyFromGitHub(
			ctx,
			network,
			repositoryOwner,
			repositoryName,
			repositoryBranch,
			dst,
		),
	)
	if err != nil {
		return util.HandleError(err)
	}

	return nil
}

func distribute(
	ctx context.Context,
	network remote.Network,
	dst,
	lock,
	command string,
	ignore []string,
	actions []action,
	deploy func() error,
) error {
	interaction.Message("ロックの取得を開始します。")
	err := tryToLock(
		lock,
		network,
	)
	if err != nil {
		interaction.Error("ロックの取得に失敗しました。")
		return util.HandleError(err)
	}
	interaction.Message("ロックを取得しました。")

	defer func() {
		interaction.Message("ロックの解除を開始します。")
		err := tryToUnlock(
			lock,
			network,
		)
		if err != nil {
			interaction.Error("ロックの解除に失敗しました。")
			util.HandleError(err)
			return
		}
		interaction.Message("ロックを解除しました。")
	}()

	interaction.Message("デプロイを開始します。")
	err = deploy()
	if err != nil {
		return util.HandleError(err)
	}
	interaction.Message("デプロイが完了しました。")

	processes := []process{}

	for _, server := range network.Servers {
		interaction.Message(fmt.Sprintf("%sへのSSH接続を開始します。", server.Host))
		authenticationMethod, err := remote.MakeAuthenticationMethod(server.SSH)
		if err != nil {
			return util.HandleError(err)
		}

		client, session, err := remote.NewSession(
			server.Host,
			server.Environments,
			authenticationMethod,
		)
		if err != nil {
			return util.HandleError(err)
		}

		stdoutPipe, err := session.StdoutPipe()
		if err != nil {
			return util.HandleError(err)
		}
		stderrPipe, err := session.StderrPipe()
		if err != nil {
			return util.HandleError(err)
		}

		stdoutFilePath := path.Join(output.GetDistributeOutputDirPath(), server.Host, "stdout")
		err = os.MkdirAll(path.Dir(stdoutFilePath), 0755)
		if err != nil {
			return util.HandleError(err)
		}
		stdoutFile, err := os.Create(stdoutFilePath)
		if err != nil {
			return util.HandleError(err)
		}

		stderrFilePath := path.Join(output.GetDistributeOutputDirPath(), server.Host, "stderr")
		err = os.MkdirAll(path.Dir(stderrFilePath), 0755)
		if err != nil {
			return util.HandleError(err)
		}
		stderrFile, err := os.Create(stderrFilePath)
		if err != nil {
			return util.HandleError(err)
		}

		processes = append(processes, process{
			client:     client,
			host:       server.Host,
			session:    session,
			stdout:     stdoutPipe,
			stderr:     stderrPipe,
			stdoutFile: stdoutFile,
			stderrFile: stderrFile,
		})

		go io.Copy(stdoutFile, stdoutPipe)

		go io.Copy(stderrFile, stderrPipe)

		go session.Run(command)

		interaction.Message(fmt.Sprintf("%sへのSSH接続が完了しました。", server.Host))
	}

	for {
		in := interaction.Choose(
			"操作を選んでください",
			len(actions)+1,
			func(index int) (string, string) {
				if len(actions) <= index {
					return "q", "quit"
				}
				return strconv.Itoa(index), actions[index].name
			},
		)

		if in == "q" {
			break
		}

		index, err := strconv.ParseInt(in, 10, 64)
		if err != nil || int64(len(actions)) <= index {
			continue
		}

		err = actions[index].callback()
		if err != nil {
			interaction.Error(err.Error())
			continue
		}
	}

	for _, process := range processes {
		process.session.Signal(ssh.SIGINT)
		process.stdoutFile.Close()
		process.stderrFile.Close()
		process.session.Close()
		process.client.Close()
	}

	return nil
}

func deloyFromLocal(
	ctx context.Context,
	network remote.Network,
	src,
	dst string,
	ignore []string,
) func() error {
	return func() error {
		interaction.Message("zipファイルの作成を開始します。")
		zipPath := path.Join(output.GetDistributeOutputDirPath(), path.Base(src)+".zip")
		err := build.Compress(src, zipPath, ignore)
		if err != nil {
			interaction.Message("zipファイルの作成に失敗しました。")
			return util.HandleError(err)
		}
		interaction.Message("zipファイルの作成に成功しました。")

		interaction.Message("プロジェクトのコピーを開始します。")
		for _, server := range network.Servers {
			interaction.Message(fmt.Sprintf("%sの処理を開始します。", server.Host))
			authenticationMethod, err := remote.MakeAuthenticationMethod(server.SSH)
			if err != nil {
				return util.HandleError(err)
			}

			err = remote.CopyFromLocal(
				ctx,
				server.Host,
				zipPath,
				dst,
				authenticationMethod,
			)
			if err != nil {
				return util.HandleError(err)
			}
			interaction.Message(fmt.Sprintf("%sの処理が完了しました。", server.Host))
		}
		interaction.Message("プロジェクトのコピーが完了しました。")

		return nil
	}
}

func deloyFromGitHub(
	ctx context.Context,
	network remote.Network,
	repositoryOwner,
	repositoryName,
	branch,
	dst string,
) func() error {
	return func() error {
		for _, server := range network.Servers {
			interaction.Message(fmt.Sprintf("%sの処理を開始します。", server.Host))
			authenticationMethod, err := remote.MakeAuthenticationMethod(server.SSH)
			if err != nil {
				return util.HandleError(err)
			}

			command := fmt.Sprintf(
				"rm -rf %s && mkdir -p %s && %s clone -b %s git@github.com:%s/%s.git %s",
				dst,
				dst,
				server.Git.Bin,
				branch,
				repositoryOwner,
				repositoryName,
				dst,
			)
			_, err = remote.Exec(
				server.Host,
				command,
				server.Environments,
				authenticationMethod,
			)
			if err != nil {
				return util.HandleError(err)
			}
			interaction.Message(fmt.Sprintf("%sの処理が完了しました。", server.Host))
		}
		return nil
	}
}

func tryToLock(lock string, network remote.Network) error {
	for _, server := range network.Servers {
		interaction.Message(fmt.Sprintf("%sのロック取得を開始します。", server.Host))
		authenticationMethod, err := remote.MakeAuthenticationMethod(server.SSH)
		if err != nil {
			return util.HandleError(err)
		}
		err = remote.Lock(
			lock,
			server.Host,
			authenticationMethod,
		)
		if err != nil {
			interaction.Error(fmt.Sprintf("%sのロック取得に失敗しました。", server.Host))
			return util.HandleError(err)
		}
		interaction.Message(fmt.Sprintf("%sのロック取得が完了しました。", server.Host))
	}

	return nil
}

func tryToUnlock(
	lock string,
	network remote.Network,
) error {
	for _, server := range network.Servers {
		interaction.Message(fmt.Sprintf("%sのロック解除を開始します。", server.Host))
		authenticationMethod, err := remote.MakeAuthenticationMethod(server.SSH)
		if err != nil {
			interaction.Error(fmt.Sprintf("%sのロック解除に失敗しました。", server.Host))
			return util.HandleError(err)
		}
		err = remote.Unlock(
			lock,
			server.Host,
			authenticationMethod,
		)
		if err != nil {
			interaction.Error(fmt.Sprintf("%sのロック解除に失敗しました。", server.Host))
			return util.HandleError(err)
		}
		interaction.Message(fmt.Sprintf("%sのロック解除が完了しました。", server.Host))
	}

	return nil
}

func makeCPUMetricsAction(servers []remote.Server) action {
	return action{
		name: "metrics-cpu",
		callback: func() error {
			var interval int64
			for {
				println("🤖    何秒間隔で取得しますか？")
				var in string
				fmt.Scan(&in)
				var err error
				interval, err = strconv.ParseInt(in, 10, 64)
				if err != nil {
					continue
				}
				break
			}

			err := metricscpu.MeasureMetrics(int(interval), servers)
			if err != nil {
				interaction.Error("CPUのメトリクス取得に失敗しました。")
				return util.HandleError(err)
			}
			return nil
		},
	}
}

func makeNginxMetricsAction(servers []nginx.Server) action {
	return action{
		name: "metrics-nginx",
		callback: func() error {
			err := metricsnginx.CopyLogFiles(servers)
			if err != nil {
				interaction.Error("NGINXのメトリクス取得に失敗しました。")
				return util.HandleError(err)
			}
			return nil
		},
	}
}

func makeMySQLMetricsAction(servers []mysql.Server) action {
	return action{
		name: "metrics-mysql",
		callback: func() error {
			err := metricsmysql.CopyLogFiles(servers)
			if err != nil {
				interaction.Error("MySQLのメトリクス取得に失敗しました。")
				return util.HandleError(err)
			}
			return nil
		},
	}
}

func makeSaveScoreAction(
	githubToken,
	repositoryOwner,
	repositoryName,
	repositoryBranch,
	slackToken,
	slackChannel string,
) action {
	return action{
		name: "save-score",
		callback: func() error {
			commit, err := github.GetCommit(
				githubToken,
				repositoryOwner,
				repositoryName,
				repositoryBranch,
			)
			if err != nil {
				return util.HandleError(err)
			}

			println("🤖    スコアを入力してください")
			var score int
			fmt.Scan(&score)

			terminate := "baby-someday:terminate"
			println("🤖    ベンチマーク結果を入力してください")
			println(fmt.Sprintf("      ※終了する場合は %s を入力してください", terminate))
			var body = fmt.Sprintf(`
### スコア
%d
		
### ブランチ
%s
		
### コミット
%s
		
### 結果
`, score, repositoryBranch, commit.Sha1)
			var githubIssueBody = body + "```\n"
			for {
				scanner := bufio.NewScanner(os.Stdin)
				if !scanner.Scan() {
					break
				}
				line := scanner.Text()
				if line == terminate {
					break
				}
				githubIssueBody += line + "\n"
			}
			githubIssueBody += "\n```\n"

			postIssueResponse, err := github.PostIssue(
				githubToken,
				repositoryOwner,
				repositoryName,
				fmt.Sprintf("ベンチマーク: Score@%d Branch@%s Commit@%s", score, repositoryBranch, commit.GetShortSha1()),
				githubIssueBody,
				[]string{github.TAG_BENCHMARK, fmt.Sprintf("branch/%s", repositoryBranch), fmt.Sprintf("commit/%s", commit.GetShortSha1())},
			)
			if err != nil {
				return util.HandleError(err)
			}

			issueID, err := postIssueResponse.GetID()
			if err != nil {
				return util.HandleError(err)
			}
			slack.PostMessage(
				slackToken,
				slackChannel,
				fmt.Sprintf("%s\nhttps://github.com/%s/%s/issues/%d", body, repositoryOwner, repositoryName, issueID),
			)

			return nil
		},
	}
}
