package tasks

import (
	"encoding/json"

	"gophercises/task/util"

	"github.com/boltdb/bolt"
)

//Task struct
type Task struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

//DeleteTask takes an index and deletes the task found in said index
func DeleteTask(index int, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Tasks"))
		util.Check(err)
		return b.Delete(util.Itob(index))
	})
}

//AddTask serializes a Task and inserts it into the DB
func AddTask(task *Task, db *bolt.DB) error {

	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Tasks"))
		util.Check(err)
		id, _ := b.NextSequence()
		task.ID = int(id)
		taskJSON, err := json.Marshal(task)
		util.Check(err)
		return b.Put(util.Itob(int(id)), taskJSON)
	})
}

//UpdateTask serializes a Task and updates it in the DB
func UpdateTask(task *Task, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Tasks"))
		util.Check(err)
		taskJSON, err := json.Marshal(task)
		util.Check(err)
		return b.Put(util.Itob(task.ID), taskJSON)
	})
}

//GetTask gets one task by index
func GetTask(index int, db *bolt.DB) *Task {
	var task Task
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Tasks"))
		util.Check(err)
		taskBytes := b.Get(util.Itob(index))
		json.Unmarshal(taskBytes, &task)
		return nil
	})
	util.Check(err)
	return &task
}

//GetTasks gets all the tasks in the db
func GetTasks(db *bolt.DB) []*Task {

	var tasks []*Task

	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Tasks"))
		util.Check(err)
		b.ForEach(func(k, v []byte) error {
			var task Task
			json.Unmarshal(v, &task)
			tasks = append(tasks, &task)
			return nil
		})
		return nil
	})

	if err != nil {
		return nil
	}
	return tasks
}
