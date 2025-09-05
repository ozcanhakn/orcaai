# OrcaAI Python SDK

The official Python SDK for the OrcaAI platform - Intelligent AI Orchestration.

## Installation

```bash
pip install orcaai
```

## Quick Start

```python
from orcaai import OrcaClient

# Initialize the client
client = OrcaClient(api_key="your-api-key")

# Send a query
result = client.query("Summarize this report")
print(result["content"])

# Get available providers
providers = client.get_providers()
print(providers)
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