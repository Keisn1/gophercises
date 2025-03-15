package cmd

import (
	"fmt"
	"log"
	"strconv"
	"task/db"

	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task on your TODO list as complete",
	Long:  "Mark a task on your TODO list as complete",
	Args:  cobra.MatchAll(cobra.ArbitraryArgs),
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				log.Fatal("Couldn't parse argument to integer")
			}
			ids = append(ids, id)
		}

		tasks, err := db.AllTasks()
		if err != nil {
			log.Fatal(err)
		}

		if len(tasks) == 0 {
			fmt.Println("No tasks to delete.")
			return
		}

		for _, id := range ids {
			if id < 1 || id > len(tasks) {
				fmt.Printf("%d Out of range\n", id)
				continue
			}

			task := tasks[id-1]
			err := db.DeleteTask(task.Key)
			if err != nil {
				fmt.Printf("Failed to mark '%d' as completed. Error: %s\n", id, err)
			}
			fmt.Printf("You have deleted the '%s' task \n", task.Value)
		}
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
