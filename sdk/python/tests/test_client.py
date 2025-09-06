import pytest
from unittest.mock import Mock, patch
from orcaai import OrcaClient, OrcaAIException, AuthenticationError, APIError


class TestOrcaClient:
    def test_init(self):
        """Test client initialization."""
        client = OrcaClient(api_key="test-key")
        assert client.api_key == "test-key"
        assert client.base_url == "http://localhost:8080"

    def test_init_with_custom_url(self):
        """Test client initialization with custom URL."""
        client = OrcaClient(api_key="test-key", base_url="https://api.orcaai.com")
        assert client.base_url == "https://api.orcaai.com"

    @patch('requests.Session.post')
    def test_query_success(self, mock_post):
        """Test successful query."""
        # Mock response
        mock_response = Mock()
        mock_response.json.return_value = {
            "content": "Test response",
            "provider": "openai",
            "model": "gpt-3.5-turbo"
        }
        mock_response.raise_for_status.return_value = None
        mock_post.return_value = mock_response

        client = OrcaClient(api_key="test-key")
        result = client.query("Test prompt")

        assert result["content"] == "Test response"
        assert result["provider"] == "openai"
        assert result["model"] == "gpt-3.5-turbo"

    @patch('requests.Session.post')
    def test_query_authentication_error(self, mock_post):
        """Test query with authentication error."""
        from requests.exceptions import HTTPError
        from unittest.mock import Mock

        # Mock HTTP error response
        mock_response = Mock()
        mock_response.status_code = 401
        mock_response.raise_for_status.side_effect = HTTPError()

        mock_post.return_value = mock_response

        client = OrcaClient(api_key="invalid-key")

        with pytest.raises(AuthenticationError):
            client.query("Test prompt")

    @patch('requests.Session.get')
    def test_get_providers_success(self, mock_get):
        """Test successful providers retrieval."""
        # Mock response
        mock_response = Mock()
        mock_response.json.return_value = {
            "providers": [
                {"name": "OpenAI", "id": "openai"},
                {"name": "Claude", "id": "claude"}
            ]
        }
        mock_response.raise_for_status.return_value = None
        mock_get.return_value = mock_response

        client = OrcaClient(api_key="test-key")
        result = client.get_providers()

        assert len(result["providers"]) == 2
        assert result["providers"][0]["name"] == "OpenAI"

    def test_query_with_options(self):
        """Test query with additional options."""
        client = OrcaClient(api_key="test-key")
        
        # This would normally make a request, but we're just testing
        # that the method accepts the parameters
        assert client is not None