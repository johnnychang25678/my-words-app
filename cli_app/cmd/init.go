package cmd

import (
	"fmt"
	"log"

	"github.com/johnnychang25678/my_words_app_2.0/db"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the app. Run this if it's your first time using the app!!!",
	Long:  "Initialize the app Run this if it's your first time using the app!!!",
	Run: func(cmd *cobra.Command, args []string) {
		if err := db.CreateAppTables(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("App successfully initialized!")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
