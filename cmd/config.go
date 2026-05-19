/*
Copyright © 2026 Rayr https://rayrsn.me/LinkInBio
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		viper.Set(key, value)

		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configPath := filepath.Join(home, ".weather-cli.yaml")
		
		if err := viper.WriteConfigAs(configPath); err != nil {
			// If file doesn't exist, WriteConfigAs will create it
			return fmt.Errorf("failed to write config: %w", err)
		}

		fmt.Printf("Set %s to %s in %s\n", key, value, configPath)
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := viper.Get(key)
		fmt.Printf("%s: %v\n", key, value)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
}
