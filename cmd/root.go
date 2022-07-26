/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "weather-Cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Version: "0.8.0",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.MarkFlagRequired("city")
	rootCmd.Flags().BoolP("raw", "r", false, "Show raw data")
	rootCmd.Flags().BoolP("version", "v", false, "Print the version number")
	rootCmd.SetVersionTemplate("Hey there, I'm version {{.Version}}")

}
