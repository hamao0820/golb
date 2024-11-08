/*
Copyright © 2024 hamao
*/
package cmd

import (
	"errors"

	"os"
	"path/filepath"

	"github.com/hamao0820/golb/golb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	LibPackage string `json:"libPackage"`
	RootDir    string `json:"rootDir"`
}

var input string
var output string
var libPackage string
var rootDir string

func init() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	configFilePath = filepath.Join(homedir, ".golb", "config.json")
}

var rootCmd = &cobra.Command{
	Use:   "golb",
	Short: "golb is a tool for bundling go source code",
	Long: `This tool bundles go source code into a single file.
You can specify the target library package and the root directory.
- lib (l): target library package
- root (r): root directory(absolute path of go.mod)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.SetConfigFile(configFilePath)
		viper.SetConfigType("json")
		if err := viper.ReadInConfig(); err != nil {
			return err
		}

		if libPackage == "" && viper.Get("libPackage") != nil {
			libPackage = viper.Get("libPackage").(string)
		}

		if rootDir == "" && viper.Get("rootDir") != nil {
			rootDir = viper.Get("rootDir").(string)
		}

		if libPackage == "" || rootDir == "" {
			return errors.New("libPackage and rootDir are required")
		}

		b := golb.NewBundler(libPackage, rootDir)
		absPath, err := filepath.Abs(input)
		if err != nil {
			return err
		}
		code, err := b.Bundle(absPath)
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
	rootCmd.Flags().StringVarP(&input, "input", "i", "", "target source file")
	rootCmd.MarkFlagRequired("src")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "output file")
	rootCmd.MarkFlagRequired("output")
	rootCmd.Flags().StringVarP(&libPackage, "lib", "l", "", "target library package")
	rootCmd.Flags().StringVarP(&rootDir, "root", "r", "", "root directory")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
