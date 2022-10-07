package repository

import (
	"database/sql"
	"fmt"
)

type wordRepository struct {
	db *sql.DB
}

var WordRepo *wordRepository

func InitWordRepo(db *sql.DB) {
	WordRepo = &wordRepository{db: db}
}

type Word struct {
	Word       string
	Definition string
}

func (w wordRepository) Upsert(word, definition string) error {
	insert := `
		INSERT INTO words (word, definition) VALUES (?, ?) 
			ON CONFLICT(word) DO UPDATE SET definition = excluded.definition;
	`
	statement, err := w.db.Prepare(insert)
	if err != nil {
		return err
	}
	_, err = statement.Exec(word, definition)
	if err == nil {
		fmt.Println("Upsert success!")
	}
	return err
}

func (w wordRepository) genBulkUpsertSql(words []Word) string {
	tuples := fmt.Sprintf("('%s', '%s')", words[0].Word, words[0].Definition)

	if len(words) > 1 {
		for i := 1; i < len(words); i++ {
			tuples += fmt.Sprintf(", ('%s', '%s')", words[i].Word, words[i].Definition)
		}
	}
	return "INSERT INTO words (word, definition) VALUES " + tuples +
		" ON CONFLICT(word) DO UPDATE SET definition = excluded.definition;"
}

func (w wordRepository) BulkUpsert(words []Word) error {
	if len(words) == 0 {
		return fmt.Errorf("len must > 0")
	}
	sql := w.genBulkUpsertSql(words)
	statement, err := w.db.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err == nil {
		fmt.Println("Upsert success!")
	}
	return err
}

func (w wordRepository) SelectAll() ([]Word, error) {
	row, err := w.db.Query("SELECT word, definition FROM wordss")
	if err != nil {
		return nil, err
	}
	defer row.Close()
	var output []Word
	for row.Next() {
		var word, def string
		row.Scan(&word, &def)
		output = append(output, Word{Word: word, Definition: def})
	}
	return output, nil
}
