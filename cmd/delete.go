package cmd

import (
	"fmt"

	"github.com/johnnychang25678/my-words-app/common"
	"github.com/johnnychang25678/my-words-app/repository"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete word",
	Short: "Delete a word from database",
	Long:  "Delete a word from database",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || len(args) > 1 {
			appErr := common.AppError{ErrorCode: common.DeleteError, Err: fmt.Errorf("Please pass in one arg the word you want to delete from database")}
			appErr.PrintAppError()
			return
		}
		word := args[0]
		rowsDeleted, err := repository.WordRepo.DeleteByWord(word)
		if err != nil {
			appErr := common.AppError{ErrorCode: common.DbError, Err: err}
			appErr.PrintAppError()
			return
		}
		if rowsDeleted == 0 {
			appErr := common.AppError{ErrorCode: common.DeleteError, Err: fmt.Errorf("There's no word: %s in databse", word)}
			appErr.PrintAppError()
			return
		}
		fmt.Printf("Successfully delete word: %s from database", word)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
