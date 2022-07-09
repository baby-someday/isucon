package mysql

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

func Analize(ptQueryDigestConfigPath string) error {
	ptQueryDigestConfig := PtQueryDigest{}
	err := util.ParseFile(ptQueryDigestConfigPath, &ptQueryDigestConfig)
	if err != nil {
		interaction.Error(err.Error())
		return err
	}

	dirIndexString := interaction.Choose(
		"ディレクトリを選択してください。",
		len(ptQueryDigestConfig.Dirs),
		func(index int) (string, string) {
			return strconv.Itoa(index), ptQueryDigestConfig.Dirs[index].Name
		},
	)
	dirIndex, err := strconv.Atoi(dirIndexString)
	if err != nil {
		return util.HandleError(err)
	}
	dir := ptQueryDigestConfig.Dirs[dirIndex]

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
		len(ptQueryDigestConfig.Presets),
		func(index int) (string, string) {
			return strconv.Itoa(index), ptQueryDigestConfig.Presets[index].Name
		},
	)
	presetIndex, err := strconv.Atoi(presetIndexString)
	if err != nil {
		return util.HandleError(err)
	}
	preset := ptQueryDigestConfig.Presets[presetIndex]

	err = RunPtQuyeryDigest(
		ptQueryDigestConfig.Bin,
		path.Join(dir.Path, fileName),
		preset,
	)
	if err != nil {
		return util.HandleError(err)
	}
	return nil
}

func RunPtQuyeryDigest(bin string, logFilePath string, preset PtQueryDigestPreset) error {
	args := []string{
		logFilePath,
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
	filename := path.Join(output.GetMySQLAnalysisDirPath(), timestamp)
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
