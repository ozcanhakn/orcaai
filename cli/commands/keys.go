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

var KeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage API keys",
	Long:  `Manage API keys for your OrcaAI account.`,
}

var keysListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all API keys",
	Long:  `List all API keys associated with your account.`,
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
		req, err := http.NewRequest("GET", apiURL+"/api/v1/keys", nil)
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
			fmt.Printf("Error making keys request: %v\n", err)
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
			fmt.Printf("Keys request failed: %s\n", string(body))
			return
		}

		// Print the response
		fmt.Println(string(body))
	},
}

var keysCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new API key",
	Long:  `Create a new API key with the specified name.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

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

		// Prepare the create key request
		keyData := map[string]string{
			"name": name,
		}

		jsonData, err := json.Marshal(keyData)
		if err != nil {
			fmt.Printf("Error marshaling key data: %v\n", err)
			return
		}

		// Create the HTTP request
		req, err := http.NewRequest("POST", apiURL+"/api/v1/keys", bytes.NewBuffer(jsonData))
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
			fmt.Printf("Error making create key request: %v\n", err)
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
		if resp.StatusCode != http.StatusCreated {
			fmt.Printf("Create key request failed: %s\n", string(body))
			return
		}

		// Print the response
		fmt.Println(string(body))
	},
}

var keysDeleteCmd = &cobra.Command{
	Use:   "delete [key-id]",
	Short: "Delete an API key",
	Long:  `Delete an API key with the specified ID.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyID := args[0]

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
		req, err := http.NewRequest("DELETE", apiURL+"/api/v1/keys/"+keyID, nil)
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
			fmt.Printf("Error making delete key request: %v\n", err)
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
			fmt.Printf("Delete key request failed: %s\n", string(body))
			return
		}

		// Print the response
		fmt.Println(string(body))
	},
}

func init() {
	KeysCmd.AddCommand(keysListCmd)
	KeysCmd.AddCommand(keysCreateCmd)
	KeysCmd.AddCommand(keysDeleteCmd)
}

func writeKeysConfig() error {
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
