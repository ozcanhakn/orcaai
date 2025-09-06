import pytest
from unittest.mock import Mock, patch
from orcaai import OrcaClient


class TestMetrics:
    @patch('requests.Session.get')
    def test_get_metrics_success(self, mock_get):
        """Test successful metrics retrieval."""
        # Mock response
        mock_response = Mock()
        mock_response.json.return_value = {
            "total_requests": 1000,
            "avg_latency": 450,
            "cost_savings": 25.50,
            "uptime": 99.9
        }
        mock_response.raise_for_status.return_value = None
        mock_get.return_value = mock_response

        client = OrcaClient(api_key="test-key")
        result = client.get_metrics()

        assert result["total_requests"] == 1000
        assert result["avg_latency"] == 450
        assert result["cost_savings"] == 25.50
        assert result["uptime"] == 99.9

    @patch('requests.Session.get')
    def test_get_metrics_authentication_error(self, mock_get):
        """Test metrics retrieval with authentication error."""
        from requests.exceptions import HTTPError

        # Mock HTTP error response
        mock_response = Mock()
        mock_response.status_code = 401
        mock_response.raise_for_status.side_effect = HTTPError()

        mock_get.return_value = mock_response

        client = OrcaClient(api_key="invalid-key")

        with pytest.raises(Exception):  # OrcaAIException or AuthenticationError
            client.get_metrics()