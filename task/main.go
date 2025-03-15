package main

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"os"
	"task/cmd"
	"task/db"
)

func main() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	dbPath := home + "/workspace/gophercises/task/instance/my.db"
	err = db.Init(dbPath)
	Must(err)
	Must(cmd.Execute())
}

func Must(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
