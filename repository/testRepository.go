package repository

import (
	"database/sql"
)

type testRepository struct {
	db *sql.DB
}

var TestRepo *testRepository

func InitTestRepo(db *sql.DB) {
	TestRepo = &testRepository{db: db}
}
