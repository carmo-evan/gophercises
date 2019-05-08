package cmd

import (
	"fmt"
	"strconv"

	"gophercises/task/tasks"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

func GetDeleteCmd(db *bolt.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "deletes your task from the database",
		Args: func(cmd *cobra.Command, args []string) error {

			if len(args) != 1 {
				return fmt.Errorf("accepts only 1 arg, received %d", len(args))
			}

			_, err := strconv.Atoi(args[0])

			if err != nil {
				return fmt.Errorf("%v is not an integer", args[0])
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			// we can ignore error because we've already checked it parses correcly
			index, _ := strconv.Atoi(args[0])
			tasks.DeleteTask(index, db)
		},
	}
}
