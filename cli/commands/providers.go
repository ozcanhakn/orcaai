package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ProvidersCmd = &cobra.Command{
	Use:   "providers",
	Short: "List available AI providers",
	Long:  `List all available AI providers and their capabilities.`,
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
		req, err := http.NewRequest("GET", apiURL+"/api/v1/ai/providers", nil)
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
			fmt.Printf("Error making providers request: %v\n", err)
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
			fmt.Printf("Providers request failed: %s\n", string(body))
			return
		}

		// Parse the response
		var providersResponse struct {
			Providers []struct {
				Name      string   `json:"name"`
				ID        string   `json:"id"`
				Models    []string `json:"models"`
				CostPer1K float64  `json:"cost_per_1k"`
				MaxTokens int      `json:"max_tokens"`
				Status    string   `json:"status"`
			} `json:"providers"`
		}
		if err := json.Unmarshal(body, &providersResponse); err != nil {
			fmt.Printf("Error parsing providers response: %v\n", err)
			return
		}

		// Print the providers
		fmt.Println("Available AI Providers:")
		for _, provider := range providersResponse.Providers {
			fmt.Printf("\n%s (%s)\n", provider.Name, provider.ID)
			fmt.Printf("  Status: %s\n", provider.Status)
			fmt.Printf("  Models: %v\n", provider.Models)
			fmt.Printf("  Cost per 1K tokens: $%.4f\n", provider.CostPer1K)
			fmt.Printf("  Max tokens: %d\n", provider.MaxTokens)
		}
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
