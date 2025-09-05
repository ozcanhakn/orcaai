# OrcaAI JavaScript SDK

The official JavaScript SDK for the OrcaAI platform - Intelligent AI Orchestration.

## Installation

```bash
npm install @orcaai/sdk
```

## Quick Start

```javascript
const { OrcaClient } = require('@orcaai/sdk');

// Initialize the client
const client = new OrcaClient({ apiKey: 'your-api-key' });

// Send a query
async function example() {
  try {
    const result = await client.query('Summarize this report');
    console.log(result.content);
    
    // Get available providers
    const providers = await client.getProviders();
    console.log(providers);
  } catch (error) {
    console.error('Error:', error.message);
  }
}

example();
```

## Features

- **Smart Routing**: AI-powered provider selection based on cost, latency, and quality
- **Caching**: Transparent caching for improved performance
- **Fallback**: Automatic failover when providers are unavailable
- **Metrics**: Usage tracking and cost optimization
- **Multi-user Support**: API key management and role-based access control

## Documentation

For full documentation, visit [https://docs.orcaai.com](https://docs.orcaai.com)

## License

MIT License