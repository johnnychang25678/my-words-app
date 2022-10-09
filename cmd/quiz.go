package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/jedib0t/go-pretty/v6/table"
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
	Use:   "quiz [count]",
	Short: "Take a quiz of [count] words from database. You need minimum 10 words in database to run this command.",
	Long: `Take a quiz of [count] words from database. 
- If no [count] arg, default 10 words per quiz.
- Use --incorrect flag to take a quiz on only the words you got wrong from last quiz.
- You need minimum 10 words in database to run this command.`,
	Run: func(cmd *cobra.Command, args []string) {
		// if word < 10 in database, cannot start quiz
		wordCount, err := repository.WordRepo.TotalWordCount()
		if err != nil {
			appErr := common.AppError{ErrorCode: common.DbError, Err: err}
			appErr.PrintAppError()
			return
		}
		if wordCount < 10 {
			appErr := common.AppError{ErrorCode: common.QuizError, Err: fmt.Errorf("Need at least 10 words in db to take quiz")}
			appErr.PrintAppError()
			return
		}

		var words []repository.Word

		if len(args) == 0 {
			// handle flag
			isIncorrectFlag, _ := cmd.Flags().GetBool("incorrect")
			if isIncorrectFlag {
				// select last quiz's incorrect words for quiz
				words, err = repository.WordRepo.SelectLastIncorrectWords()
				if err != nil {
					appErr := common.AppError{ErrorCode: common.DbError, Err: err}
					appErr.PrintAppError()
					return
				}
				if len(words) == 0 {
					fmt.Println("You don't have any incorrects in last quiz!")
					return
				}
			} else {
				// quizCount = 10 if no args
				quizCount := 10
				words, err = repository.WordRepo.RandomSelectWords(quizCount)
				if err != nil {
					appErr := common.AppError{ErrorCode: common.DbError, Err: err}
					appErr.PrintAppError()
					return
				}
			}
		} else if len(args) > 1 {
			appErr := common.AppError{ErrorCode: common.QuizError, Err: fmt.Errorf("only accept one argument")}
			appErr.PrintAppError()
			return
		} else {
			// get quizCount from arg
			quizCount, err := strconv.Atoi(args[0])
			if err != nil {
				appErr := common.AppError{ErrorCode: common.QuizError, Err: fmt.Errorf("arg must be a number")}
				appErr.PrintAppError()
				return
			}
			if quizCount > wordCount {
				appErr := common.AppError{ErrorCode: common.QuizError,
					Err: fmt.Errorf("The number you input > word count in db. Please provide a smaller number. db word count: %d", wordCount)}
				appErr.PrintAppError()
				return
			}

			words, err = repository.WordRepo.RandomSelectWords(quizCount)
			if err != nil {
				appErr := common.AppError{ErrorCode: common.DbError, Err: err}
				appErr.PrintAppError()
				return
			}
		}

		// generate questions
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

		var correctCount int
		var incorrectCount int
		var results []repository.Result

		// prompt user
		fmt.Printf("---- TOTAL QUESTIONS: %d ----\n", len(questions))
		qIndex := 0
		for qIndex < len(questions) {
			question := questions[qIndex]
			currentWord := question.word
			options := question.options
			answer := question.answer

			result := repository.Result{
				Word:          currentWord,
				IsCorrect:     false,
				Definition:    answer,
				UserSelection: "",
			}

			prompt := promptui.Select{
				Label:        fmt.Sprintf("Q%d: %s", qIndex+1, currentWord),
				Items:        options,
				HideSelected: true,
			}

			_, userSelectedDef, err := prompt.Run()
			if err != nil {
				appErr := common.AppError{ErrorCode: common.QuizError, Err: err}
				appErr.PrintAppError()
				return
			}
			result.UserSelection = userSelectedDef
			if userSelectedDef != answer {
				incorrectCount++
				result.IsCorrect = false
				fmt.Printf("You're wrong!\nYou select: '%s', but the definition of the word '%s' is: '%s'\n",
					userSelectedDef, currentWord, answer)
			} else {
				correctCount++
				result.IsCorrect = true
				fmt.Printf("You're correct!\nThe definition of the word '%s' is: '%s'\n",
					currentWord, userSelectedDef)
			}
			results = append(results, result)
			qIndex++
		}

		// write quiz result to db
		if err := repository.TestRepo.Insert(results, correctCount, incorrectCount); err != nil {
			appErr := common.AppError{ErrorCode: common.DbError, Err: err}
			appErr.PrintAppError()
			return
		}
		// print out the result table
		PrintQuizSummary(results, correctCount, incorrectCount)
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

// this function share with transcript command
func PrintQuizSummary(results []repository.Result, correctCount int, incorrectCount int) {
	table1 := table.NewWriter()
	table1.SetOutputMirror(os.Stdout)
	table1.AppendHeader(table.Row{"word", "definition"})

	table2 := table.NewWriter()
	table2.SetOutputMirror(os.Stdout)
	table2.AppendHeader(table.Row{"word", "you selected", "correct definition"})

	for _, result := range results {
		if result.IsCorrect {
			table1.AppendRow(table.Row{
				result.Word, result.Definition,
			})
		} else {
			table2.AppendRow(table.Row{
				result.Word, result.UserSelection, result.Definition,
			})
		}
	}
	fmt.Println("---- CORRECTS ----")
	table1.Render()
	fmt.Println("---- INCORRECTS ----")
	table2.Render()

	fmt.Println("---- SUMMARY ----")
	table3 := table.NewWriter()
	table3.SetOutputMirror(os.Stdout)
	table3.AppendRow(table.Row{"correct", correctCount})
	table3.AppendRow(table.Row{"incorrect", incorrectCount})
	table3.Render()
	fmt.Printf("Score: %.f%%\n", float64(100*correctCount/(correctCount+incorrectCount)))
}

func init() {
	rootCmd.AddCommand(quizCmd)
	quizCmd.Flags().Bool("incorrect", false, "Quiz incorrect words on last quiz")
}
