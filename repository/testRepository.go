package repository

import (
	"database/sql"
	"fmt"
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
