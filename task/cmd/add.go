package cmd

import (
	"gophercises/task/tasks"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

func GetAddCmd(db *bolt.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "adds a task to the database",
		Run: func(cmd *cobra.Command, args []string) {
			var task tasks.Task
			task.Text = strings.Join(args, " ")
			task.Completed = false
			tasks.AddTask(&task, db)
		},
	}
}
