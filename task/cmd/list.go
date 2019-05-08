package cmd

import (
	"fmt"
	"gophercises/task/tasks"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

func GetListCmd(db *bolt.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "lists your tasks",
		Run: func(cmd *cobra.Command, args []string) {
			tasks := tasks.GetTasks(db)
			for _, task := range tasks {
				fmt.Println(strconv.Itoa(task.ID) + " - " + task.Text + " - Done: " + strconv.FormatBool(task.Completed))
			}

		},
	}
}
