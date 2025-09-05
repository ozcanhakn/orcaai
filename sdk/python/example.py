#!/usr/bin/env python3
"""
Example usage of the OrcaAI Python SDK
"""

import os
import sys

# Add the orcaai package to the path
sys.path.insert(0, os.path.join(os.path.dirname(__file__), "."))

from orcaai import OrcaClient


def main():
    # Get API key from environment variable or prompt
    api_key = os.environ.get("ORCAAI_API_KEY")
    if not api_key:
        api_key = input("Enter your OrcaAI API key: ")

    # Initialize the client
    client = OrcaClient(api_key=api_key)

    try:
        # Get available providers
        print("Getting available providers...")
        providers = client.get_providers()
        print(f"Available providers: {len(providers.get('providers', []))}")

        # Send a query
        print("\nSending a query...")
        result = client.query(
            prompt="Explain what artificial intelligence is in simple terms",
            task_type="text-generation"
        )
        print(f"Response: {result['content']}")

        # Get metrics
        print("\nGetting metrics...")
        metrics = client.get_metrics()
        print(f"Metrics: {metrics}")

        print("\nExample completed successfully!")

    except Exception as e:
        print(f"Error: {e}")


if __name__ == "__main__":
    main()