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
	Use:   "transcript [quizId]",
	Short: "Show transcript of [quizId] quiz, if no quizId is provided, will show the latest quiz. Use --history to get the quizId.",
	Long: `Show transcript of [quizId] quiz
- If no quizId is provided, will show the latest quiz. 
- Use --history to get your quiz history, and get the quizIds.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			appErr := common.AppError{ErrorCode: common.TranscriptError, Err: fmt.Errorf("Too many arguments. Please -h to see usage.")}
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
				appErr := common.AppError{ErrorCode: common.TranscriptError, Err: err}
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
	transcriptCmd.Flags().Bool("history", false, "show all tests record")
}
