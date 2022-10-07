package main

import (
	"log"

	"github.com/johnnychang25678/my-words-app/cmd"
	"github.com/johnnychang25678/my-words-app/db"
)

func main() {
	_, err := db.ConnectToDB()
	if err != nil {
		log.Fatal(err) // close the app if fail to connect to db
	}
	cmd.Execute()
}
