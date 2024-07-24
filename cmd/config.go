/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFilePath string

func init() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	configFilePath = filepath.Join(homedir, ".golb", "config.json")
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.SetConfigFile(configFilePath)
		viper.SetConfigType("json")
		viper.ReadInConfig()

		if err := viper.ReadInConfig(); err != nil {
			return err
		}

		if libPackage != "" {
			viper.Set("libPackage", libPackage)
		}

		if rootDir != "" {
			viper.Set("rootDir", rootDir)
		}

		if err := viper.WriteConfig(); err != nil {
			return err
		}

		fmt.Println("libPackage:", viper.Get("libPackage"))
		fmt.Println("rootDir:", viper.Get("rootDir"))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringVarP(&libPackage, "lib", "l", "", "target library package")
	configCmd.Flags().StringVarP(&rootDir, "root", "r", "", "root directory")
}
