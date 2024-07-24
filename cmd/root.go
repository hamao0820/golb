/*
Copyright © 2024 hamao
*/
package cmd

import (
	"golb/golb"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "golb",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := golb.Bundle("golb/testdata/src/a/main.go")
		return err
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
