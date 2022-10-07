package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func ConnectToDB() (*sql.DB, error) {
	var err error
	path := "./sqlite-database.db"
	if _, err := os.Stat(path); err != nil {
		file, err := os.Create(path)
		defer file.Close()
		if err != nil {
			return nil, err
		}
	}
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		fmt.Println("connect to db error:", err)
		return nil, err
	}
	fmt.Println("connected to db")

	return db, nil
}

func createTable(sql ...string) error {
	for _, s := range sql {
		statement, err := db.Prepare(s)
		if err != nil {
			return err
		}
		_, err = statement.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateAppTables() error {
	words := `
	 	CREATE TABLE IF NOT EXISTS words (
	 		word VARCHAR(100) UNIQUE NOT NULL CHECK (word <> ''),
	 		definition TEXT NOT NULL CHECK (definition <> ''),
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP
	 	);
	 `

	tests := `
		CREATE TABLE IF NOT EXISTS tests (
			date DATETIME DEFAULT CURRENT_TIMESTAMP,
			correct_count INTEGER NOT NULL,
			incorrect_count INTEGER NOT NULL
		);
	`

	// questions := `
	// 	CREATE TABLE IF NOT EXISTS questions (
	// 		test_id INTEGER NOT NULL,
	// 		correct_count INTEGER NOT NULL,
	// 		incorrect_count INTEGER NOT NULL
	// 	);
	// `
	// questionOptions := `
	// 	CREATE TABLE IF NOT EXISTS question_options (
	// 		question_id INTEGER NOT NULL,
	// 		correct_count INTEGER NOT NULL,
	// 		incorrect_count INTEGER NOT NULL
	// 	);
	// `

	testWord := `
		CREATE TABLE IF NOT EXISTS test_word (
			test_id INTEGER NOT NULL, 
			word VARCHAR(100) NOT NULL,
			is_correct TINYINT NOT NULL,
			definition TEXT NOT NULL,
			user_selection TEXT NOT NULL
		);
	`
	return createTable(words, tests, testWord)
}
