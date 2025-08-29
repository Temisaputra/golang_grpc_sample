package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "be-example",
	Short: "BE Example with Clean Architecture",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// register subcommands
	rootCmd.AddCommand(grpcCmd)
	rootCmd.AddCommand(migrateCmd)
}
