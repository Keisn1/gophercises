package cmd

import (
	"fmt"
	"strings"
	"task/db"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task to your TODO list",
	Long:  "Add a new task to your TODO list",
	Args:  cobra.MatchAll(cobra.ArbitraryArgs),
	Run: func(cmd *cobra.Command, args []string) {
		task := strings.Join(args, " ")
		id, err := db.CreateTask(task)
		if err != nil {
			fmt.Println("Something went wrong", err)
			return
		}
		fmt.Printf("Created task '%s' with id %d\n", task, id)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
