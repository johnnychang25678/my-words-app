package cmd

import (
	"fmt"
	"log"

	"github.com/johnnychang25678/my-words-app/db"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the app",
	Long:  "Initialize the app",
	Run: func(cmd *cobra.Command, args []string) {
		if err := db.CreateAppTables(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("App successfully initialized!")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
