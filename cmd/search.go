package cmd

import (
	"fmt"

	"github.com/johnnychang25678/my-words-app/client"
	"github.com/johnnychang25678/my-words-app/common"
	"github.com/johnnychang25678/my-words-app/repository"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search word",
	Short: "Search word definition from api",
	Long:  `Search word definition from api`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || len(args) > 1 {
			appErr := common.AppError{ErrorCode: common.SearchError, Err: fmt.Errorf("Please provide exactly one argument")}
			appErr.PrintAppError()
			return
		}
		word := args[0]

		client := client.NewDictionaryClient()
		fmt.Println("fetching data from the internet...")
		definitions, err := client.GetDefinitions(word)
		if err != nil {
			appErr := common.AppError{ErrorCode: common.ApiError, Err: err}
			appErr.PrintAppError()
			return
		}

		fmt.Printf("These are the definitions of word: %s\n", word)
		for i, d := range definitions {
			fmt.Printf("%d. %s\n", i+1, d)
		}
		// prompt user if they want to add to database
		prompt := promptui.Prompt{
			Label:     "do you want to add or update this word to the database?",
			IsConfirm: true,
		}
		_, err = prompt.Run()
		if err != nil { // if user doesn't respone with y or Y, will return error
			fmt.Println("ok!")
			return
		}
		var items []string
		for _, d := range definitions {
			items = append(items, d)
		}
		prompt2 := promptui.Select{
			Label: "Pick the definition",
			Items: items,
		}
		_, pick, err := prompt2.Run()
		if err != nil {
			appErr := common.AppError{ErrorCode: common.SearchError, Err: err}
			appErr.PrintAppError()
			return
		}
		// upsert to db
		if err := repository.WordRepo.Upsert(word, pick); err != nil {
			appErr := common.AppError{ErrorCode: common.DbError, Err: err}
			appErr.PrintAppError()
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}