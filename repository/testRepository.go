package repository

import (
	"database/sql"
	"fmt"

	"github.com/johnnychang25678/my-words-app/common"
)

type Result struct {
	Word          string
	IsCorrect     bool
	Definition    string
	UserSelection string
}

func (r Result) boolToInt() int {
	if r.IsCorrect {
		return 1
	}
	return 0
}

type Test struct {
	Id             int
	CorrectCount   int
	IncorrectCount int
	Date           string
}

func (t Test) GetScore() float64 {
	return float64(100 * t.CorrectCount / (t.CorrectCount + t.IncorrectCount))
}

type testRepository struct {
	db *sql.DB
}

var TestRepo *testRepository

func InitTestRepo(db *sql.DB) {
	TestRepo = &testRepository{db: db}
}

// when inserting test, should use a transaction to insert testWord
// so we don't need testWordRepository.

func (t testRepository) Insert(results []Result, correctCount int, incorrectCount int) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}

	insertTestSql := fmt.Sprintf("INSERT INTO tests (correct_count, incorrect_count) VALUES (%d, %d);",
		correctCount, incorrectCount)
	statement1Result, err := t.handleTransaction(tx, insertTestSql)

	lastInsertId, err := statement1Result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	insertTestWordSql := "INSERT INTO test_word (test_id, word, is_correct, definition, user_selection) VALUES "
	for i, result := range results {
		if i == len(results)-1 {
			insertTestWordSql += fmt.Sprintf("(%d, '%s', %d, '%s', '%s');",
				lastInsertId, result.Word, result.boolToInt(), result.Definition, result.UserSelection)
		} else {
			insertTestWordSql += fmt.Sprintf("(%d, '%s', %d, '%s', '%s'), ",
				lastInsertId, result.Word, result.boolToInt(), result.Definition, result.UserSelection)
		}

	}

	if _, err = t.handleTransaction(tx, insertTestWordSql); err != nil {
		return err
	}

	return tx.Commit()
}

func (t testRepository) SelectAllTests() ([]Test, error) {
	row, err := t.db.Query("SELECT rowid, correct_count, incorrect_count, date FROM tests;")
	if err != nil {
		return nil, err
	}
	defer row.Close()
	var tests []Test
	for row.Next() {
		var id, correct, incorrect int
		var date string
		row.Scan(&id, &correct, &incorrect, &date)

		tests = append(tests, Test{
			Id:             id,
			CorrectCount:   correct,
			IncorrectCount: incorrect,
			Date:           common.ToLocalTimeString(date),
		})
	}
	return tests, nil
}
func (t testRepository) SelectTestById(testId int) (*Test, error) {
	sql := "SELECT rowid, correct_count, incorrect_count, date FROM tests WHERE rowid = ?;"
	row := t.db.QueryRow(sql, testId)
	var id, correct, incorrect int
	var date string
	err := row.Scan(&id, &correct, &incorrect, &date)
	if err != nil {
		return nil, err
	}
	return &Test{
		Id:             id,
		CorrectCount:   correct,
		IncorrectCount: incorrect,
		Date:           common.ToLocalTimeString(date),
	}, err
}
func (t testRepository) SelectLatestTest() (*Test, error) {
	sql := "SELECT rowid, correct_count, incorrect_count, date FROM tests ORDER BY rowid DESC LIMIT 1;"
	row := t.db.QueryRow(sql)
	var id, correct, incorrect int
	var date string
	err := row.Scan(&id, &correct, &incorrect, &date)
	if err != nil {
		return nil, err
	}
	return &Test{
		Id:             id,
		CorrectCount:   correct,
		IncorrectCount: incorrect,
		Date:           common.ToLocalTimeString(date),
	}, err
}

// func (t testRepository) SelectLatestTestResult() (*Result, error) {
// 	sql := "SELECT word, is_correct, definition, user_selection FROM test_word WHERE test_id = (SELECT rowid FROM tests ORDER BY rowid DESC LIMIT 1);"
// 	row := t.db.QueryRow(sql)
// 	var word string
// 	var isCorrectInt int
// 	var def string
// 	var userSelect string
// 	err := row.Scan(&word, &isCorrectInt, &def, &userSelect)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var isCorrect bool
// 	if isCorrectInt == 1 {
// 		isCorrect = true
// 	}
// 	return &Result{
// 		Word:          word,
// 		IsCorrect:     isCorrect,
// 		Definition:    def,
// 		UserSelection: userSelect,
// 	}, err

// }

func (t testRepository) SelectTestResultById(testId int) ([]Result, error) {
	sql := fmt.Sprintf("SELECT word, is_correct, definition, user_selection FROM test_word WHERE test_id = %d;", testId)
	row, err := t.db.Query(sql)
	if err != nil {
		return nil, err
	}
	var results []Result
	for row.Next() {
		var word string
		var isCorrectInt int
		var def string
		var userSelect string
		row.Scan(&word, &isCorrectInt, &def, &userSelect)

		var isCorrect bool
		if isCorrectInt == 1 {
			isCorrect = true
		}
		results = append(results, Result{
			Word:          word,
			IsCorrect:     isCorrect,
			Definition:    def,
			UserSelection: userSelect,
		})
	}
	return results, nil
}

// handle rollback
func (t testRepository) handleTransaction(tx *sql.Tx, statement string) (sql.Result, error) {
	s, err := tx.Prepare(statement)
	defer s.Close()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	res, err := s.Exec()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return res, nil
}
