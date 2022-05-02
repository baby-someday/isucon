package cmd

import (
	"bytes"
	"io/ioutil"
	"log"
	"os/exec"
	"path"
	"time"

	"github.com/baby-someday/isucon/internal/benchbot"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var benchbotCmd = &cobra.Command{
	Use:   "benchbot",
	Short: "benchbot",
	Long:  `benchbot`,
	Args:  validateBenchbotArgs,
	Run:   runBenchbotCommand,
}

func init() {
	benchbotCmd.Flags().String(
		FLAG_CONFIG_PATH,
		"",
		"config yaml file path",
	)
	benchbotCmd.MarkFlagRequired(FLAG_CONFIG_PATH)
	rootCmd.AddCommand(benchbotCmd)
}

func validateBenchbotArgs(cmd *cobra.Command, args []string) error {
	return nil
}

func runBenchbotCommand(cmd *cobra.Command, args []string) {
	configFilePath, err := cmd.Flags().GetString(FLAG_CONFIG_PATH)
	if err != nil {
		log.Fatal(err.Error())
	}

	configFileBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	config := benchbot.Config{}
	err = yaml.Unmarshal(configFileBytes, &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	uuidObj, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err.Error())
	}

	repositoryDir := path.Join(config.TMP, uuidObj.String())
	err = exec.Command("git", "clone", config.Repository.URL, repositoryDir).Run()
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		for _, branch := range config.Repository.Branches {
			var stderr bytes.Buffer
			var stdout bytes.Buffer
			exec.Command("cd", repositoryDir, "&&", "/usr/bin/git", "fetch", "&&", "/usr/bin/git", "checkout", branch)
			command := exec.Command("cd", repositoryDir, "&&", "/usr/bin/git", "log", "--oneline", "-1", "--pretty=%H")
			err = command.Run()
			log.Println(err)
			if err != nil {
				log.Fatal(err.Error())
			}
			command.Stderr = &stderr
			command.Stdout = &stdout
			println(stderr.String())
			println(stdout.String())
		}

		time.Sleep(time.Second * time.Duration(config.Sleep))
	}

	// // TODO: YAMLにする
	// // TODO: ログファイル等もSlackに送る
	// values := url.Values{}
	// // TODO: 引数で渡す
	// values.Set("token", "xoxb-3096051984995-3467105680852-5OGaN27oCTPWRT1ORXJJCmrg")
	// values.Set("channel", "#2022練習")
	// values.Set("text", "benchbot test")

	// request, err := http.NewRequest(
	// 	"POST",
	// 	"https://slack.com/api/chat.postMessage",
	// 	strings.NewReader(values.Encode()),
	// )
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// client := &http.Client{}
	// response, err := client.Do(request)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// defer response.Body.Close()
}
