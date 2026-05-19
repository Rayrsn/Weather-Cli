/*
Copyright © 2026 Rayr https://rayrsn.me/LinkInBio/
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "weather-Cli",
	Short: "An app made to get weather information through the terminal",
	Long: `You can use this app to get weather information through the terminal.
it can also be used with the --raw flag to get the response in json format.`,
	Version: "2.0.0",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.weather-cli.yaml)")

	rootCmd.Flags().BoolP("version", "v", false, "Print the version number")
	rootCmd.SetVersionTemplate("Hey there, I'm version {{.Version}}")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".weather-cli")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		// Quietly continue if config exists
	}
}
