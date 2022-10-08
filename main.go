package main

import (
	"database/sql"
	"log"

	"github.com/johnnychang25678/my-words-app/cmd"
	"github.com/johnnychang25678/my-words-app/db"
	"github.com/johnnychang25678/my-words-app/repository"
	"github.com/joho/godotenv"
)

func initRepos(db *sql.DB) {
	repository.InitWordRepo(db)
	repository.InitTestRepo(db)
	repository.InitTestWordRepo(db)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	database, err := db.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	initRepos(database)
	cmd.Execute()
}
