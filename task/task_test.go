package main

import (
	"gophercises/task/tasks"
	"log"
	"testing"

	"github.com/boltdb/bolt"
)

func TestAdd(t *testing.T) {

	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var task tasks.Task
	task.Text = "Testing"
	task.Completed = false
	err = tasks.AddTask(&task, db)

	if err != nil {
		t.Logf("Failed adding task to DB")
	}

}
