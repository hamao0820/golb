/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFilePath string

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config command",
	Long: `This command sets the configuration file and shows the current configuration.
The configuration file is stored in $HOME/.golb/config.json.
You can set the library package and the root directory.
- lib (l): target library package
- root (r): root directory(absolute path of go.mod)`,

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
