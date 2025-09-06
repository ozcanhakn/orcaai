import pytest
import os
import sys

# Add the orcaai package to the path
sys.path.insert(0, os.path.join(os.path.dirname(__file__), ".."))

@pytest.fixture
def api_key():
    """Return a test API key."""
    return os.environ.get("ORCAAI_API_KEY", "test-api-key")

@pytest.fixture
def base_url():
    """Return a test base URL."""
    return os.environ.get("ORCAAI_BASE_URL", "http://localhost:8080")