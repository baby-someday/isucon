package distribute

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/baby-someday/isucon/internal/metricsnginx"
	"github.com/baby-someday/isucon/pkg/build"
	"github.com/baby-someday/isucon/pkg/github"
	"github.com/baby-someday/isucon/pkg/output"
	"github.com/baby-someday/isucon/pkg/remote"
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
			makeNginxMetricsAction(network.Servers),
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
	githubToken,
	repositoryOwner,
	repositoryName,
	repositoryURL,
	repositoryBranch,
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
			makeNginxMetricsAction(network.Servers),
			makeSaveScoreAction(
				githubToken,
				repositoryOwner,
				repositoryName,
				repositoryBranch,
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
		return err
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
	err := tryToLock(
		lock,
		network,
	)
	if err != nil {
		log.Println("ロックに失敗しました、他のベンチマークが実行中です。")
		return err
	}

	err = deploy()
	if err != nil {
		return err
	}

	processes := []process{}

	// TODO: Closeちゃんとやる
	for _, server := range network.Servers {
		authenticationMethod, err := remote.MakeAuthenticationMethod(server)
		if err != nil {
			return err
		}

		client, session, err := remote.NewSession(
			server.Host,
			server.Environments,
			authenticationMethod,
		)

		stdoutPipe, err := session.StdoutPipe()
		if err != nil {
			return err
		}
		stderrPipe, err := session.StderrPipe()
		if err != nil {
			return err
		}

		stdoutFilePath := path.Join(output.GetDistributeOutputDirPath(), server.Host, "stdout")
		err = os.MkdirAll(path.Dir(stdoutFilePath), 0755)
		if err != nil {
			return err
		}
		stdoutFile, err := os.Create(stdoutFilePath)
		if err != nil {
			return err
		}

		stderrFilePath := path.Join(output.GetDistributeOutputDirPath(), server.Host, "stderr")
		err = os.MkdirAll(path.Dir(stderrFilePath), 0755)
		if err != nil {
			return err
		}
		stderrFile, err := os.Create(stderrFilePath)
		if err != nil {
			return err
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

		go func() {
			io.Copy(stdoutFile, stdoutPipe)
		}()

		go func() {
			io.Copy(stderrFile, stderrPipe)
		}()

		go func() {
			session.Run(command)
		}()
	}

	for {
		println("🤖    操作を選んでください")
		print("👉    ")
		for index, action := range actions {
			print(fmt.Sprintf("%d:%s    ", index, action.name))
		}
		print("q:quit")
		println()

		var in string
		fmt.Scan(&in)

		if in == "q" {
			break
		}

		index, err := strconv.ParseInt(in, 10, 64)
		if err != nil || int64(len(actions)) <= index {
			continue
		}

		err = actions[index].callback()
		if err != nil {
			log.Println(err.Error())
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

	err = tryToUnlock(
		lock,
		network,
	)
	if err != nil {
		return err
	}

	return nil
}

func deloyFromLocal(ctx context.Context, network remote.Network, src, dst string, ignore []string) func() error {
	return func() error {
		zipPath := path.Join(output.GetDistributeOutputDirPath(), path.Base(src)+".zip")
		err := build.Compress(src, zipPath, ignore)
		if err != nil {
			return err
		}

		for _, server := range network.Servers {
			authenticationMethod, err := remote.MakeAuthenticationMethod(server)
			if err != nil {
				return err
			}

			err = remote.CopyFromLocal(
				ctx,
				server.Host,
				zipPath,
				dst,
				authenticationMethod,
			)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func deloyFromGitHub(ctx context.Context, network remote.Network, repositoryOwner, repositoryName, branch, dst string) func() error {
	return func() error {
		for _, server := range network.Servers {
			authenticationMethod, err := remote.MakeAuthenticationMethod(server)
			if err != nil {
				return err
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
				return err
			}
		}
		return nil
	}
}

func tryToLock(lock string, network remote.Network) error {
	for _, server := range network.Servers {
		authenticationMethod, err := remote.MakeAuthenticationMethod(server)
		// TODO Unlockちゃんとやる
		if err != nil {
			return err
		}
		err = remote.Lock(
			lock,
			server.Host,
			authenticationMethod,
		)
		// TODO Unlockちゃんとやる
		if err != nil {
			return err
		}
	}

	return nil
}

func tryToUnlock(lock string, network remote.Network) error {
	for _, server := range network.Servers {
		authenticationMethod, err := remote.MakeAuthenticationMethod(server)
		// TODO Unlockちゃんとやる
		if err != nil {
			return err
		}
		err = remote.Unlock(
			lock,
			server.Host,
			authenticationMethod,
		)
		// TODO Unlockちゃんとやる
		if err != nil {
			return err
		}
	}

	return nil
}

func makeNginxMetricsAction(servers []remote.Server) action {
	return action{
		name: "metrics-nginx",
		callback: func() error {
			err := metricsnginx.CopyFiles(servers)
			if err != nil {
				log.Println("Nginxのメトリクス取得に失敗しました。")
				return err
			}
			return nil
		},
	}
}

func makeSaveScoreAction(token, owner, repositoryName, branch string) action {
	return action{
		name: "save-score",
		callback: func() error {
			commit, err := github.GetCommit(
				token,
				owner,
				repositoryName,
				branch,
			)
			if err != nil {
				return err
			}

			println("🤖    スコアを入力してください")
			var score int
			fmt.Scan(&score)

			terminate := "baby-someday:terminate"
			println("🤖    ベンチマーク結果を入力してください")
			println(fmt.Sprintf("      ※終了する場合は %s を入力してください", terminate))
			var result = fmt.Sprintf(`
### スコア
%d
		
### ブランチ
%s
		
### コミット
%s
		
### 結果
`, score, branch, commit.Sha1)
			result += "```\n"
			for {
				scanner := bufio.NewScanner(os.Stdin)
				if !scanner.Scan() {
					break
				}
				line := scanner.Text()
				if line == terminate {
					break
				}
				result += line + "\n"
			}
			result += "\n```\n"

			err = github.PostIssue(
				token,
				owner,
				repositoryName,
				fmt.Sprintf("ベンチマーク: Score@%d Branch@%s Commit@%s", score, branch, commit.GetShortSha1()),
				result,
				[]string{github.TAG_BENCHMARK, fmt.Sprintf("branch/%s", branch), fmt.Sprintf("commit/%s", commit.GetShortSha1())},
			)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
