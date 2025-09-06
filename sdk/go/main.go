package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ozcanhakn/orcaai-go/sdk/go/client"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("ORCAAI_API_KEY")
	if apiKey == "" {
		fmt.Println("ORCAAI_API_KEY environment variable is required")
		return
	}

	// Create client
	c := client.NewClient(client.Config{
		APIKey:  apiKey,
		BaseURL: "http://localhost:8080",
		Timeout: 30 * time.Second,
	})

	fmt.Printf("Client created successfully: %v\n", c)
}
