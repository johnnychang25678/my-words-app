package repository

import (
	"database/sql"
	"fmt"

	"github.com/johnnychang25678/my-words-app/common"
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
	CreateTime string
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
	return w.query("SELECT word, definition, create_time FROM words ORDER BY rowid DESC")
}

func (w wordRepository) SelectWithLimit(limit int) ([]Word, error) {
	return w.query(fmt.Sprintf("SELECT word, definition, create_time FROM words ORDER BY rowid DESC LIMIT %d", limit))
}

func (w wordRepository) SelectByWord(word string) ([]Word, error) {
	return w.query(fmt.Sprintf("SELECT word, definition, create_time FROM words WHERE word = '%s'", word))
}

func (w wordRepository) RandomSelectWords(limit int) ([]Word, error) {
	return w.query(fmt.Sprintf("SELECT word, definition, create_time FROM words ORDER BY RANDOM() LIMIT %d", limit))
}

func (w wordRepository) SelectInWords(wordStrings []string) ([]Word, error) {
	sql := "SELECT word, definition, create_time FROM words WHERE word IN ("
	for i, word := range wordStrings {
		if i == len(wordStrings)-1 {
			sql += fmt.Sprintf("'%s'", word)
		} else {
			sql += fmt.Sprintf("'%s', ", word)
		}
	}
	sql += ");"
	return w.query(sql)
}

func (w wordRepository) SelectLastIncorrectWords() ([]Word, error) {
	row, err := w.db.Query("SELECT word from test_word WHERE test_id = (SELECT rowid FROM tests ORDER BY rowid DESC LIMIT 1) AND is_correct = 0;")
	if err != nil {
		return nil, err
	}
	defer row.Close()
	var words []string
	for row.Next() {
		var word string
		row.Scan(&word)
		words = append(words, word)
	}
	return w.SelectInWords(words)
}

func (w wordRepository) TotalWordCount() (int, error) {
	row := w.db.QueryRow("SELECT count(1) from words")
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (w wordRepository) query(sql string) ([]Word, error) {
	row, err := w.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	var output []Word
	for row.Next() {
		var word, def, t string
		row.Scan(&word, &def, &t)
		output = append(output, Word{Word: word, Definition: def, CreateTime: common.ToLocalTimeString(t)})
	}
	return output, nil
}
