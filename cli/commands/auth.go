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

var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
	Long:  `Authentication commands for managing your OrcaAI account.`,
}

var authLoginCmd = &cobra.Command{
	Use:   "login [email] [password]",
	Short: "Login to your OrcaAI account",
	Long:  `Login to your OrcaAI account and store the API key for future use.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		email := args[0]
		password := args[1]

		// Prepare the login request
		loginData := map[string]string{
			"email":    email,
			"password": password,
		}

		jsonData, err := json.Marshal(loginData)
		if err != nil {
			fmt.Printf("Error marshaling login data: %v\n", err)
			return
		}

		// Make the login request
		apiURL := viper.GetString("api-url")
		if apiURL == "" {
			apiURL = "http://localhost:8080"
		}

		resp, err := http.Post(apiURL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error making login request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		// Read the response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response: %v\n", err)
			return
		}

		// Check if login was successful
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Login failed: %s\n", string(body))
			return
		}

		// Parse the response
		var loginResponse struct {
			Token string `json:"token"`
		}
		if err := json.Unmarshal(body, &loginResponse); err != nil {
			fmt.Printf("Error parsing login response: %v\n", err)
			return
		}

		// Store the API key in config
		viper.Set("api-key", loginResponse.Token)
		if err := writeConfig(); err != nil {
			fmt.Printf("Error saving API key: %v\n", err)
			return
		}

		fmt.Println("Login successful! API key saved to config.")
	},
}

var authRegisterCmd = &cobra.Command{
	Use:   "register [email] [password] [name]",
	Short: "Register a new OrcaAI account",
	Long:  `Register a new OrcaAI account.`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		email := args[0]
		password := args[1]
		name := args[2]

		// Prepare the registration request
		regData := map[string]string{
			"email":    email,
			"password": password,
			"name":     name,
		}

		jsonData, err := json.Marshal(regData)
		if err != nil {
			fmt.Printf("Error marshaling registration data: %v\n", err)
			return
		}

		// Make the registration request
		apiURL := viper.GetString("api-url")
		if apiURL == "" {
			apiURL = "http://localhost:8080"
		}

		resp, err := http.Post(apiURL+"/api/v1/auth/register", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error making registration request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		// Read the response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response: %v\n", err)
			return
		}

		// Check if registration was successful
		if resp.StatusCode != http.StatusCreated {
			fmt.Printf("Registration failed: %s\n", string(body))
			return
		}

		fmt.Println("Registration successful!")
	},
}

func init() {
	AuthCmd.AddCommand(authLoginCmd)
	AuthCmd.AddCommand(authRegisterCmd)
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
