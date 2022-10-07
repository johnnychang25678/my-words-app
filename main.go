package main

import (
	"database/sql"
	"log"

	"github.com/johnnychang25678/my-words-app/cmd"
	"github.com/johnnychang25678/my-words-app/db"
	"github.com/johnnychang25678/my-words-app/repository"
)

func initRepos(db *sql.DB) {
	repository.InitWordRepo(db)
	repository.InitTestRepo(db)
	repository.InitTestWordRepo(db)
}

func main() {
	database, err := db.ConnectToDB()
	if err != nil {
		log.Fatal(err) // close the app if fail to connect to db
	}
	initRepos(database)
	cmd.Execute()
}
