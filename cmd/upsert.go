package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/johnnychang25678/my-words-app/common"
	"github.com/johnnychang25678/my-words-app/repository"
	"github.com/spf13/cobra"
)

// upsertCmd represents the upsert command
var upsertCmd = &cobra.Command{
	Use:   "upsert {word definition}",
	Short: "Insert or update a word and its definition to the database. Use --file flag to upsert with a csv file.",
	Long: `Insert or update a word and its definition to the database.
- Either upsert one word with terminal or multiple words at once with --file filename.
- With --file, make sure the file exists, it's a csv file, and the first column is word, second column is definition.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fileName, _ := cmd.Flags().GetString("file")
			if fileName == "" {
				printUpsertErr("Please -h to see the usage")
				return
			}
			file, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
			if err != nil {
				appErr := common.AppError{ErrorCode: common.OpenFileError, Err: err}
				appErr.PrintAppError()
				return
			}
			words, appErr := getWordsFromFile(file)
			if appErr != nil {
				appErr.PrintAppError()
				return
			}
			// bulk upsert to db
			if err := repository.WordRepo.BulkUpsert(words); err != nil {
				appErr := common.AppError{ErrorCode: common.DbError, Err: err}
				appErr.PrintAppError()
			}
		} else if len(args) != 2 {
			printUpsertErr("Please provide two args [word] [definition]")
		} else {
			// upsert one word to db
			word, definition := args[0], args[1]
			if word == "" {
				printUpsertErr("word cannot be empty")
				return
			}
			if definition == "" {
				printUpsertErr("definition cannot be empty")
				return
			}
			if err := repository.WordRepo.Upsert(word, definition); err != nil {
				appErr := common.AppError{ErrorCode: common.DbError, Err: err}
				appErr.PrintAppError()
			}
		}
	},
}

func printUpsertErr(errMessage string) {
	err := fmt.Errorf(errMessage)
	appErr := common.AppError{ErrorCode: common.UpsertError, Err: err}
	appErr.PrintAppError()
}

func getWordsFromFile(file *os.File) ([]repository.Word, *common.AppError) {
	defer file.Close()
	csvReader := csv.NewReader(file)
	var words []repository.Word
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, &common.AppError{ErrorCode: common.ReadCsvError, Err: err}
		}
		if len(record) != 2 {
			err = fmt.Errorf("csv format is incorrect. Please make sure there are exact two columns (separate with comma) for every row")
			return nil, &common.AppError{ErrorCode: common.ReadCsvError, Err: err}
		}
		words = append(words, repository.Word{Word: record[0], Definition: record[1]})
	}
	if len(words) == 0 {
		err := fmt.Errorf("There are no words in the file")
		return nil, &common.AppError{ErrorCode: common.ReadCsvError, Err: err}
	}
	return words, nil
}

func init() {
	rootCmd.AddCommand(upsertCmd)
	upsertCmd.Flags().StringP("file", "f", "", "pass csv file name as arg to upsert all the words and definitions at once")
}
