# OrcaAI Go SDK

The official Go SDK for the OrcaAI platform - Intelligent AI Orchestration.

OrcaAI is an intelligent AI orchestration platform that routes AI requests to the best provider based on cost, latency, quality, and availability. This SDK provides a simple interface to interact with the OrcaAI platform from Go applications.

## Features

- **Smart Routing**: AI-powered provider selection based on cost, latency, and quality
- **Caching**: Transparent caching for improved performance
- **Fallback**: Automatic failover when providers are unavailable
- **Metrics**: Usage tracking and cost optimization
- **Multi-user Support**: API key management and role-based access control

## Installation

```bash
go get github.com/ozcanhakn/orcaai-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "os"
    "time"

    "github.com/ozcanhakn/orcaai-go/client"
)

func main() {
    // Create client
    c := client.NewClient(client.Config{
        APIKey:  os.Getenv("ORCAAI_API_KEY"),
        BaseURL: "http://localhost:8080",
        Timeout: 30 * time.Second,
    })

    // Send a query
    queryReq := client.QueryRequest{
        Prompt:   "Explain what artificial intelligence is in simple terms",
        TaskType: "text-generation",
    }

    resp, err := c.Query(queryReq)
    if err != nil {
        log.Fatalf("Failed to send query: %v", err)
    }

    fmt.Printf("Response: %s\n", resp.Content)
}
```

## Documentation

For full documentation, visit [https://docs.orcaai.com](https://docs.orcaai.com)

## License

MIT License

## Support

For support, email support@orcaai.com or file an issue on GitHub.