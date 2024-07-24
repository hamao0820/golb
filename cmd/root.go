/*
Copyright © 2024 hamao
*/
package cmd

import (
	"golb/golb"
	"os"

	"github.com/spf13/cobra"
)

var libPackage string
var rootDir string

var rootCmd = &cobra.Command{
	Use:   "golb",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		b := golb.NewBundler("golb/testdata/src/a/main.go", libPackage, rootDir)
		err := b.Bundle()
		return err
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&libPackage, "lib", "l", "github.com/hamao0820/ac-library-go", "target library package")
	rootCmd.PersistentFlags().StringVarP(&rootDir, "root", "r", "golb/testdata", "root directory")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
