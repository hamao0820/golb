/*
Copyright © 2024 hamao
*/
package cmd

import (
	"golb/golb"
	"os"

	"github.com/spf13/cobra"
)

var input string
var output string
var libPackage string
var rootDir string

var rootCmd = &cobra.Command{
	Use:   "golb",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		b := golb.NewBundler(input, libPackage, rootDir)
		code, err := b.Bundle()
		if err != nil {
			return err
		}

		// 書き込む
		f, err := os.Create(output)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.WriteString(code)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&input, "input", "i", "", "target source file")
	rootCmd.MarkPersistentFlagRequired("src")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output file")
	rootCmd.MarkPersistentFlagRequired("output")
	rootCmd.PersistentFlags().StringVarP(&libPackage, "lib", "l", "github.com/hamao0820/ac-library-go", "target library package")
	rootCmd.PersistentFlags().StringVarP(&rootDir, "root", "r", "golb/testdata", "root directory")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
