package commands

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var MetricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Get usage metrics",
	Long:  `Get usage metrics for your OrcaAI account.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get API key from config or flag
		apiKey := viper.GetString("api-key")
		if apiKey == "" {
			fmt.Println("API key not found. Please login or set API key with 'orcaai config set api-key YOUR_KEY'")
			return
		}

		// Get API URL from config or flag
		apiURL := viper.GetString("api-url")
		if apiURL == "" {
			apiURL = "http://localhost:8080"
		}

		// Create the HTTP request
		req, err := http.NewRequest("GET", apiURL+"/api/v1/metrics", nil)
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			return
		}

		// Set headers
		req.Header.Set("Authorization", "Bearer "+apiKey)

		// Make the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making metrics request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		// Read the response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response: %v\n", err)
			return
		}

		// Check if request was successful
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Metrics request failed: %s\n", string(body))
			return
		}

		// Print the response
		fmt.Println(string(body))
	},
}

func writeConfig() error {
	// Get the config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	// Create the orcaai config directory
	orcaaiConfigDir := filepath.Join(configDir, "orcaai")
	if err := os.MkdirAll(orcaaiConfigDir, 0755); err != nil {
		return err
	}

	// Set the config file name and type
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(orcaaiConfigDir)

	// Write the config file
	return viper.WriteConfig()
}
