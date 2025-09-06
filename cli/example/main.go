package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ozcanhakn/orcaai-go/client"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("ORCAAI_API_KEY")
	if apiKey == "" {
		log.Fatal("ORCAAI_API_KEY environment variable is required")
	}

	// Create client
	c := client.NewClient(client.Config{
		APIKey:  apiKey,
		BaseURL: "http://localhost:8080",
		Timeout: 30 * time.Second,
	})

	// Get providers
	fmt.Println("Getting providers...")
	providers, err := c.GetProviders()
	if err != nil {
		log.Fatalf("Failed to get providers: %v", err)
	}

	fmt.Printf("Found %d providers:\n", len(providers.Providers))
	for _, provider := range providers.Providers {
		fmt.Printf("  - %s (%s)\n", provider.Name, provider.ID)
	}

	// Send a query
	fmt.Println("\nSending query...")
	queryReq := client.QueryRequest{
		Prompt:   "Explain what artificial intelligence is in simple terms",
		TaskType: "text-generation",
	}

	resp, err := c.Query(queryReq)
	if err != nil {
		log.Fatalf("Failed to send query: %v", err)
	}

	fmt.Printf("Response: %s\n", resp.Content)
	fmt.Printf("Provider: %s\n", resp.Provider)
	fmt.Printf("Model: %s\n", resp.Model)
	fmt.Printf("Cost: $%.4f\n", resp.Cost)
	fmt.Printf("Latency: %d ms\n", resp.Latency)

	// Get metrics
	fmt.Println("\nGetting metrics...")
	metrics, err := c.GetMetrics()
	if err != nil {
		log.Fatalf("Failed to get metrics: %v", err)
	}

	fmt.Printf("Total Requests: %d\n", metrics.TotalRequests)
	fmt.Printf("Average Latency: %d ms\n", metrics.AvgLatency)
	fmt.Printf("Cost Savings: $%.2f\n", metrics.CostSavings)
	fmt.Printf("Uptime: %.2f%%\n", metrics.Uptime)
}
