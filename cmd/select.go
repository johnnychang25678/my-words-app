package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/johnnychang25678/my-words-app/common"
	"github.com/johnnychang25678/my-words-app/repository"
	"github.com/spf13/cobra"
)

// selectCmd represents the select command
var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "print the words and definitions in database",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// if no args, check --all, --incorrect, --count flags
		// TODO: add prompt to check if output to csv file
		// TODO: --incorrect flags
		var selectedWords []repository.Word
		var err error

		if len(args) == 0 {
			isAll, _ := cmd.Flags().GetBool("all")
			if isAll {
				selectedWords, err = repository.WordRepo.SelectAll()
			} else {
				count, _ := cmd.Flags().GetInt("count")
				selectedWords, err = repository.WordRepo.SelectWithLimit(count)
			}

		} else {
			selectedWords, err = repository.WordRepo.SelectByWord(args[0])
		}
		if err != nil {
			appErr := common.AppError{ErrorCode: common.DbError, Err: err}
			appErr.PrintAppError()
			return
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"word", "definition", "create time"})

		if len(selectedWords) == 0 {
			t.Render()
			return
		}

		for _, word := range selectedWords {
			t.AppendRow(table.Row{
				word.Word, word.Definition, word.CreateTime,
			})
		}
		t.Render()

	},
}

func init() {
	rootCmd.AddCommand(selectCmd)
	selectCmd.Flags().IntP("count", "c", 10, "Select the latest [int] words in database.")
	selectCmd.Flags().BoolP("all", "a", false, "Select all the words in database.")
}
