package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var QueryCmd = &cobra.Command{
	Use:   "query [prompt]",
	Short: "Send a query to the OrcaAI platform",
	Long:  `Send a query to the OrcaAI platform and get an AI-generated response.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prompt := args[0]

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

		// Prepare the query request
		queryData := map[string]string{
			"prompt": prompt,
		}

		jsonData, err := json.Marshal(queryData)
		if err != nil {
			fmt.Printf("Error marshaling query data: %v\n", err)
			return
		}

		// Create the HTTP request
		req, err := http.NewRequest("POST", apiURL+"/api/v1/ai/query", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			return
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		// Make the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making query request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		// Read the response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response: %v\n", err)
			return
		}

		// Check if query was successful
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Query failed: %s\n", string(body))
			return
		}

		// Parse the response
		var queryResponse struct {
			Content string `json:"content"`
		}
		if err := json.Unmarshal(body, &queryResponse); err != nil {
			fmt.Printf("Error parsing query response: %v\n", err)
			return
		}

		// Print the response
		fmt.Println(queryResponse.Content)
	},
}

func writeQueryConfig() error {
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
