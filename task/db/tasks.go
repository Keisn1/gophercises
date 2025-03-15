package db

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"strconv"
	"time"
)

var taskBucket = []byte("tasks")
var db *bolt.DB

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}

type Task struct {
	Key   int
	Value string
}

func CreateTask(task string) (int, error) {
	db, err := bolt.Open("instance/my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var id int
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(taskBucket))

		id64, _ := b.NextSequence()
		id = int(id64)
		key := itob(id)
		return b.Put(key, []byte(task))
	})
	if err != nil {
		return -1, err
	}
	return id, nil
}

func AllTasks() ([]Task, error) {
	db, err := bolt.Open("instance/my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var tasks []Task
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(taskBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			tasks = append(tasks, Task{
				Key:   btoi(k),
				Value: string(v),
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func DeleteTask(key int) error {
	db, err := bolt.Open("instance/my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(taskBucket))
		return b.Delete(itob(key))
	})
}

func DoTask(args []string) {
	tNum, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal("Couldn't parse argument to integer")
	}
	if tNum < 1 {
		log.Fatal("Not a valid position")
	}

	db, err := bolt.Open("instance/my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(taskBucket))
		c := b.Cursor()

		count := 1
		k, v := c.First()
		for count < tNum && k != nil {
			k, v = c.Next()
			count++
		}
		if k == nil {
			log.Fatalf("You have seeked to a non-valid position %d", tNum)
		}

		var task Task
		err := json.Unmarshal(v, &task)
		if err != nil {
			log.Fatal(err)
		}

		b.Delete(k)
		fmt.Printf("You have completed the '%s' task.\n", task.Value)
		return nil
	})
}

func Init(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(taskBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}
