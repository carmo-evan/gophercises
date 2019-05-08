package cmd

import (
	"gophercises/task/util"
	"strconv"

	"gophercises/task/tasks"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

func GetDoCmd(db *bolt.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "do",
		Short: "completes your task",
		Long:  `completes your task`,
		Run: func(cmd *cobra.Command, args []string) {
			index, err := strconv.Atoi(args[0])
			util.Check(err)
			task := tasks.GetTask(index, db)
			task.Completed = true
			tasks.UpdateTask(task, db)
		},
	}
}
