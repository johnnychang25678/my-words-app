/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/johnnychang25678/my-words-app/common"
	"github.com/johnnychang25678/my-words-app/repository"
	"github.com/spf13/cobra"
)

// transcriptCmd represents the transcript command
var transcriptCmd = &cobra.Command{
	Use:   "transcript",
	Short: "Show transcript of a quiz taken",
	Long:  `Show transcript of a quiz taken`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			appErr := common.AppError{ErrorCode: common.TrasncriptError, Err: fmt.Errorf("Too many arguments. Please -h to see usage.")}
			appErr.PrintAppError()
			return
		}
		if len(args) == 0 {
			isHisotry, _ := cmd.Flags().GetBool("history")
			if isHisotry {
				// print history
				tests, err := repository.TestRepo.SelectAllTests()
				if err != nil {
					appErr := common.AppError{ErrorCode: common.DbError, Err: err}
					appErr.PrintAppError()
					return
				}
				t := table.NewWriter()
				t.SetOutputMirror(os.Stdout)
				t.AppendHeader(table.Row{"id", "correct", "incorrect", "score", "date"})
				for _, test := range tests {
					t.AppendRow(table.Row{
						test.Id, test.CorrectCount, test.IncorrectCount, fmt.Sprintf("%.f%%", test.GetScore()), test.Date,
					})
				}
				t.Render()
			} else {
				test, err := repository.TestRepo.SelectLatestTest()
				if err != nil {
					appErr := common.AppError{ErrorCode: common.DbError, Err: err}
					appErr.PrintAppError()
					return
				}

				latestTestId := test.Id

				var results []repository.Result

				results, err = repository.TestRepo.SelectTestResultById(latestTestId)
				if err != nil {
					appErr := common.AppError{ErrorCode: common.DbError, Err: err}
					appErr.PrintAppError()
					return
				}
				fmt.Println("Quiz date:", test.Date)
				PrintQuizSummary(results, test.CorrectCount, test.IncorrectCount)
			}
		} else {
			// with arg, select result by id
			id, err := strconv.Atoi(args[0])
			if err != nil {
				appErr := common.AppError{ErrorCode: common.TrasncriptError, Err: err}
				appErr.PrintAppError()
				return
			}
			test, err := repository.TestRepo.SelectTestById(id)
			if err != nil {
				appErr := common.AppError{ErrorCode: common.DbError, Err: err}
				appErr.PrintAppError()
				return
			}
			results, err := repository.TestRepo.SelectTestResultById(id)
			if err != nil {
				appErr := common.AppError{ErrorCode: common.DbError, Err: err}
				appErr.PrintAppError()
				return
			}
			fmt.Println("Quiz date:", test.Date)
			PrintQuizSummary(results, test.CorrectCount, test.IncorrectCount)
		}

	},
}

func init() {
	rootCmd.AddCommand(transcriptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// transcriptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// transcriptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	transcriptCmd.Flags().Bool("history", false, "show all tests record")
}
