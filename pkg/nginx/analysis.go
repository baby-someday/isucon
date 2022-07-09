package nginx

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/baby-someday/isucon/pkg/interaction"
	"github.com/baby-someday/isucon/pkg/output"
	"github.com/baby-someday/isucon/pkg/util"
)

func Analize(alpConfigPath string) error {
	alpConfig := ALP{}
	err := util.ParseFile(alpConfigPath, &alpConfig)
	if err != nil {
		interaction.Error(err.Error())
		return err
	}

	dirIndexString := interaction.Choose(
		"ディレクトリを選択してください。",
		len(alpConfig.Dirs),
		func(index int) (string, string) {
			return strconv.Itoa(index), alpConfig.Dirs[index].Name
		},
	)
	dirIndex, err := strconv.Atoi(dirIndexString)
	if err != nil {
		return util.HandleError(err)
	}
	dir := alpConfig.Dirs[dirIndex]

	files, err := ioutil.ReadDir(dir.Path)
	if err != nil {
		return util.HandleError(err)
	}

	fileNames := []string{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".log") {
			continue
		}
		fileNames = append(fileNames, file.Name())
	}

	fileNameIndexString := interaction.Choose(
		"ファイルを選択してください。",
		len(fileNames),
		func(index int) (string, string) {
			return strconv.Itoa(index), fileNames[index]
		},
	)
	fileNameIndex, err := strconv.Atoi(fileNameIndexString)
	if err != nil {
		return util.HandleError(err)
	}
	fileName := fileNames[fileNameIndex]

	presetIndexString := interaction.Choose(
		"プリセットを選択してください。",
		len(alpConfig.ALPPresets),
		func(index int) (string, string) {
			return strconv.Itoa(index), alpConfig.ALPPresets[index].Name
		},
	)
	presetIndex, err := strconv.Atoi(presetIndexString)
	if err != nil {
		return util.HandleError(err)
	}
	preset := alpConfig.ALPPresets[presetIndex]

	err = RunAlp(
		alpConfig.Bin,
		path.Join(dir.Path, fileName),
		preset,
	)
	if err != nil {
		return util.HandleError(err)
	}
	return nil
}

func RunAlp(bin string, logFilePath string, preset ALPPreset) error {
	args := []string{
		"ltsv",
		"--file",
		logFilePath,
	}
	if preset.M != "" {
		args = append(args, "-m", preset.M)
	}
	if preset.O != "" {
		args = append(args, "-o", preset.O)
	}
	if preset.Q != "" {
		args = append(args, "-q", preset.Q)
	}
	if preset.QsIgnoreValues {
		args = append(args, "--qs-ignore-values")
	}
	if preset.R != "" {
		args = append(args, "-r", preset.R)
	}
	if preset.ShowFooters {
		args = append(args, "--show-footers")
	}
	if preset.Sort != "" {
		args = append(args, "--sort", preset.Sort)
	}
	if preset.Extra != "" {
		args = append(args, strings.Split(preset.Extra, " ")...)
	}
	command := fmt.Sprintf(
		"%s %s",
		bin,
		strings.Join(args, " "),
	)
	interaction.Message("以下のコマンドを実行します。")
	interaction.Message(command)
	bytes, err := exec.Command(
		bin,
		args...,
	).Output()
	if err != nil {
		return util.HandleError(err)
	}
	now := time.Now()
	timestamp := now.Format("2006-01-02_15:04:05.txt")
	filename := path.Join(output.GetNginxAnalysisDirPath(), timestamp)
	err = os.MkdirAll(path.Dir(filename), os.ModePerm)
	if err != nil {
		return util.HandleError(err)
	}
	err = ioutil.WriteFile(filename, bytes, os.ModePerm)
	if err != nil {
		return util.HandleError(err)
	}
	return err
}
