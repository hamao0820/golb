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
		b := golb.NewBundler("golb/testdata/src/a/main.go", "github.com/hamao0820/ac-library-go", "golb/testdata")
		err := b.Bundle()
		return err
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
