/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"os"
	"path"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initial command",
	Long:  `This command initializes the configuration file. The configuration file is stored in $HOME/.golb/config.json.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		if !existsDir(path.Join(home, ".golb")) {
			err := createDir(path.Join(home, ".golb"))
			if err != nil {
				return err
			}
		}

		if !existsFile(path.Join(home, ".golb", "config.json")) {
			config := Config{}

			f, err := os.Create(path.Join(home, ".golb", "config.json"))
			if err != nil {
				return err
			}
			defer f.Close()

			enc := json.NewEncoder(f)
			err = enc.Encode(config)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func existsFile(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func createFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}

func existsDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	return stat.IsDir()
}

func createDir(path string) error {
	err := os.Mkdir(path, 075)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
