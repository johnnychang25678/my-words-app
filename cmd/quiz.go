package cmd

import (
	"context"
	"fmt"
	"sync"

	"github.com/johnnychang25678/my-words-app/common"
	"github.com/johnnychang25678/my-words-app/repository"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type Question struct {
	word    string
	answer  string
	options []string // one definition and 2 other random defintions
}

// quizCmd represents the quiz command
var quizCmd = &cobra.Command{
	Use:   "quiz",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		wordCount, err := repository.WordRepo.TotalWordCount()
		if err != nil {
			appErr := common.AppError{ErrorCode: common.DbError, Err: err}
			appErr.PrintAppError()
			return
		}
		if wordCount < 10 {
			appErr := common.AppError{ErrorCode: common.QuizError, Err: fmt.Errorf("Need at least 10 words in db to do quiz")}
			appErr.PrintAppError()
			return
		}

		// randomly select 10 words to quiz
		words, err := repository.WordRepo.RandomSelectWords(10)
		if err != nil {
			appErr := common.AppError{ErrorCode: common.DbError, Err: err}
			appErr.PrintAppError()
			return
		}
		// use for loop + go routine to generate questions
		questions := make([]Question, len(words))
		wg := &sync.WaitGroup{}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		for i, word := range words {
			wg.Add(1)
			go func(index int, currentWord repository.Word) {
				defer wg.Done()

				select {
				case <-ctx.Done():
					return
				default:
				}

				defs, err := getRandomTwoDefinitions(currentWord.Word)
				if err != nil {
					appErr := common.AppError{ErrorCode: common.DbError, Err: err}
					appErr.PrintAppError()
					cancel()
					return
				}
				defs = append(defs, currentWord.Definition)
				common.ShuffleSlice(defs)
				question := Question{
					word:    currentWord.Word,
					answer:  currentWord.Definition,
					options: defs,
				}
				questions[index] = question

			}(i, word)
		}
		wg.Wait()
		qIndex := 0
		// use promptui to generate question and options
		for qIndex < len(questions) {
			question := questions[qIndex]
			currentWord := question.word
			options := question.options
			answer := question.answer

			prompt := promptui.Select{
				Label:        fmt.Sprintf("Q%d. %s", qIndex+1, currentWord),
				Items:        options,
				HideSelected: true,
			}

			_, userSelectedDef, err := prompt.Run()
			if err != nil {
				fmt.Println("err", err)
				return
			}
			if userSelectedDef != answer {
				fmt.Printf("You're wrong!\nYou select: '%s', but the definition of the word '%s' is: '%s'\n",
					userSelectedDef, currentWord, answer)

			} else {
				fmt.Printf("You're correct!\nThe definition of the word '%s' is: '%s'\n",
					currentWord, userSelectedDef)
			}
			qIndex++
		}

		// after user select, write its result to db table

	},
}

func getRandomTwoDefinitions(currentWord string) ([]string, error) {
	output := make([]string, 2)
	words, err := repository.WordRepo.RandomSelectWords(2)
	if err != nil {
		return nil, err
	}
	for i, word := range words {
		if word.Word == currentWord {
			return getRandomTwoDefinitions(currentWord)
		}
		output[i] = word.Definition
	}
	return output, nil
}

func init() {
	rootCmd.AddCommand(quizCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// quizCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// quizCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
