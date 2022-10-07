package repository

import (
	"database/sql"
)

type testWordRepository struct {
	db *sql.DB
}

var TestWordRepo *testWordRepository

func InitTestWordRepo(db *sql.DB) {
	TestWordRepo = &testWordRepository{db: db}
}
