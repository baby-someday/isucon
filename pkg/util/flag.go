package util

import (
	"io/ioutil"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func ParseFlag(cmd *cobra.Command, flag string, object interface{}) error {
	filePath, err := cmd.Flags().GetString(flag)
	if err != nil {
		return HandleError(err)
	}

	return ParseFile(
		filePath,
		object,
	)
}

func ParseFile(filePath string, object interface{}) error {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return HandleError(err)
	}

	return yaml.Unmarshal(fileBytes, object)
}
