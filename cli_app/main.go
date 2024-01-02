package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/johnnychang25678/my-words-app/db"

	"github.com/johnnychang25678/my-words-app/cmd"
	"github.com/johnnychang25678/my-words-app/repository"
	"github.com/joho/godotenv"
)

func initRepos(db *sql.DB) {
	repository.InitWordRepo(db)
	repository.InitTestRepo(db)
}

func main() {
	if err := godotenv.Load(); err != nil {
		if os.IsNotExist(err) {
			if _, err := os.Create(".env"); err != nil {
				log.Fatal(err)
			}
			fmt.Println(".env file not exist. It's created now, please try again. Check README to see .env settings.")
			os.Exit(0)
		}
		log.Fatal(err)
	}
	database, err := db.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	initRepos(database)
	cmd.Execute()
}
