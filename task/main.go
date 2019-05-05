package main

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
	"gophercises/task/cmd"
)

func main() {

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var rootCmd = &cobra.Command{
		Use:   "task",
		Short: "Task is a task manager for your command line",
	}

	var doCmd = cmd.GetDoCmd(db)
	doCmd.SetUsageTemplate("Usage:\ntask do [index] [flags]\nFlags:\n-h, --help   help for do")
	rootCmd.AddCommand(doCmd)

	var deleteCmd = cmd.GetDeleteCmd(db)
	deleteCmd.SetUsageTemplate("Usage:\ntask delete [index] [flags]\nFlags:\n-h, --help   help for delete")
	rootCmd.AddCommand(deleteCmd)

	var listCmd = cmd.GetListCmd(db)
	rootCmd.AddCommand(listCmd)

	var addCmd = cmd.GetAddCmd(db)
	addCmd.SetUsageTemplate("Usage:\ntask add [text] [flags]\nFlags:\n-h, --help   help for add")
	rootCmd.AddCommand(addCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
