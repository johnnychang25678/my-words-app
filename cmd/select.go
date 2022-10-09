package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/johnnychang25678/my-words-app/common"
	"github.com/johnnychang25678/my-words-app/repository"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// selectCmd represents the select command
var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Print or create a csv file of the words and definitions in database. If no flags are provided, will select latest 10 words.",
	Long: `Print or create a csv file of the words and definitions in database.
- If no flags are provided, will select latest 10 words.
- A prompt will check if you want to output as a csv or print out to terminal.
- If output to csv, the file will be created under your currect directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		// if no args, check --all, --incorrect, --count flags
		var selectedWords []repository.Word
		var err error

		if len(args) == 0 {
			isAll, _ := cmd.Flags().GetBool("all")
			isIncorrect, _ := cmd.Flags().GetBool("incorrect")
			if isAll {
				selectedWords, err = repository.WordRepo.SelectAll()
			} else if isIncorrect {
				selectedWords, err = repository.WordRepo.SelectLastIncorrectWords()
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

		if len(selectedWords) == 0 {
			fmt.Println("You do not have any words in your database yet")
			return
		}

		prompt := promptui.Prompt{
			Label:     "Do you want to output a csv file? (If no, will print result to the terminal)",
			IsConfirm: true,
		}
		_, err = prompt.Run()
		if err != nil {
			// err: user select no
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"word", "definition", "create time"})

			for _, word := range selectedWords {
				t.AppendRow(table.Row{
					word.Word, word.Definition, word.CreateTime,
				})
			}
			t.Render()

		} else {
			fileName := "output-" + time.Now().Format("20060102T150405") + ".csv"
			csvFile, err := os.Create(fileName)
			defer csvFile.Close()
			if err != nil {
				appErr := common.AppError{ErrorCode: common.CreateFileError, Err: err}
				appErr.PrintAppError()
				return
			}
			csvWriter := csv.NewWriter(csvFile)
			var data [][]string
			for _, word := range selectedWords {
				row := []string{word.Word, word.Definition, word.CreateTime}
				data = append(data, row)
			}

			for _, d := range data {
				csvWriter.Write(d)
			}
			csvWriter.Flush()
			path, _ := filepath.Abs(fileName)
			fmt.Println("Your csv file is created at:", path)
		}

	},
}

func init() {
	rootCmd.AddCommand(selectCmd)
	selectCmd.Flags().IntP("count", "c", 10, "Select the latest [int] words in database.")
	selectCmd.Flags().BoolP("all", "a", false, "Select all the words in database.")
	selectCmd.Flags().Bool("incorrect", false, "Select all the words you got wrong from last quiz.")
}
