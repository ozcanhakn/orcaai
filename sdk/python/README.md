# OrcaAI Python SDK

[![PyPI version](https://badge.fury.io/py/orcaai.svg)](https://badge.fury.io/py/orcaai)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

The official Python SDK for the OrcaAI platform - Intelligent AI Orchestration.

OrcaAI is an intelligent AI orchestration platform that routes AI requests to the best provider based on cost, latency, quality, and availability. This SDK provides a simple interface to interact with the OrcaAI platform from Python applications.

## Features

- **Smart Routing**: AI-powered provider selection based on cost, latency, and quality
- **Caching**: Transparent caching for improved performance
- **Fallback**: Automatic failover when providers are unavailable
- **Metrics**: Usage tracking and cost optimization
- **Multi-user Support**: API key management and role-based access control

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
result = client.query("Explain what artificial intelligence is in simple terms")
print(result["content"])

# Get available providers
providers = client.get_providers()
print(providers)
```

## Documentation

For full documentation, visit [https://docs.orcaai.com](https://docs.orcaai.com)

## License

MIT License

## Support

For support, email support@orcaai.com or file an issue on GitHub.