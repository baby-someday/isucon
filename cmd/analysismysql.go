package cmd

import (
	"github.com/baby-someday/isucon/pkg/interaction"
	"github.com/baby-someday/isucon/pkg/mysql"
	"github.com/spf13/cobra"
)

var analysisMySQLCmd = &cobra.Command{
	Use:   "mysql",
	Short: "mysql",
	Long:  `mysql`,
	Run:   runAnalysisMySQLCommand,
}

func init() {
	analysisMySQLCmd.Flags().String(
		FLAG_PT_QUERY_DIGEST_PATH,
		FLAG_PT_QUERY_DIGEST_PATH_DEFAULT,
		"",
	)

	analysisCmd.AddCommand(analysisMySQLCmd)
}

func runAnalysisMySQLCommand(cmd *cobra.Command, args []string) {
	interaction.Message("pt-query-digestの実行を開始します。")
	filePath, err := cmd.Flags().GetString(FLAG_PT_QUERY_DIGEST_PATH)
	if err != nil {
		interaction.Error(err.Error())
		return
	}

	err = mysql.Analize(filePath)
	if err != nil {
		interaction.Error(err.Error())
		return
	}
	interaction.Message("pt-query-digestの実行が完了しました。")
}
