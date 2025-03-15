package cmd

import (
	"github.com/spf13/cobra"
)

var TodoBucket string = "TODOList"

var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "task is a CLI for managing your TODOs.",
	Long:  "task is a CLI for managing your TODOs.",
}

func Execute() error {
	return rootCmd.Execute()
}
