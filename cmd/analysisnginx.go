package cmd

import (
	"github.com/baby-someday/isucon/pkg/interaction"
	"github.com/baby-someday/isucon/pkg/nginx"
	"github.com/spf13/cobra"
)

var analysisNginxCmd = &cobra.Command{
	Use:   "nginx",
	Short: "nginx",
	Long:  `nginx`,
	Run:   runAnalysisNginxCommand,
}

func init() {
	analysisNginxCmd.Flags().String(
		FLAG_ALP_PATH,
		FLAG_ALP_PATH_DEFAULT,
		"",
	)

	analysisCmd.AddCommand(analysisNginxCmd)
}

func runAnalysisNginxCommand(cmd *cobra.Command, args []string) {
	interaction.Message("ALPの実行を開始します。")
	filePath, err := cmd.Flags().GetString(FLAG_ALP_PATH)
	if err != nil {
		interaction.Error(err.Error())
		return
	}

	err = nginx.Analize(filePath)
	if err != nil {
		interaction.Error(err.Error())
		return
	}
	interaction.Message("ALPの実行が完了しました。")
}
