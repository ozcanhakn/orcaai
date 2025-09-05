package main

import (
	"fmt"
	"os"

	"orcaai-cli/commands"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Global flags
	apiKey string
	apiURL string
	debug  bool
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "orcaai",
		Short: "OrcaAI CLI - Intelligent AI Orchestration",
		Long: `OrcaAI CLI is a command line interface for the OrcaAI platform.
It allows you to interact with AI providers, manage your account, and monitor usage.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Initialize viper
			initConfig()
		},
	}

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&apiKey, "api-key", "k", "", "API key for authentication")
	rootCmd.PersistentFlags().StringVarP(&apiURL, "api-url", "u", "", "API URL")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")

	// Add subcommands
	rootCmd.AddCommand(commands.VersionCmd)
	rootCmd.AddCommand(commands.ConfigCmd)
	rootCmd.AddCommand(commands.AuthCmd)
	rootCmd.AddCommand(commands.QueryCmd)
	rootCmd.AddCommand(commands.ProvidersCmd)
	rootCmd.AddCommand(commands.KeysCmd)
	rootCmd.AddCommand(commands.MetricsCmd)

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Set default values
	viper.SetDefault("api-url", "http://localhost:8080")

	// Read config from file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Find home directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("Error finding config directory: %v\n", err)
		return
	}

	viper.AddConfigPath(configDir + "/orcaai")

	// Read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if debug {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}

	// Override with command line flags if provided
	if apiKey != "" {
		viper.Set("api-key", apiKey)
	}
	if apiURL != "" {
		viper.Set("api-url", apiURL)
	}
}
