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
	fmt.Println("connected to db!")

	return db, nil
}
