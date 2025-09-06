import requests
import json
from typing import Dict, Any, Optional
from .exceptions import OrcaAIException, AuthenticationError, APIError


class OrcaClient:
    """OrcaAI Client for interacting with the OrcaAI API."""

    def __init__(self, api_key: str, base_url: str = "http://localhost:8080"):
        """
        Initialize the OrcaClient.

        Args:
            api_key (str): Your OrcaAI API key
            base_url (str): Base URL for the OrcaAI API (default: http://localhost:8080)
        """
        self.api_key = api_key
        self.base_url = base_url.rstrip('/')
        self.session = requests.Session()
        self.session.headers.update({
            "Authorization": f"Bearer {self.api_key}",
            "Content-Type": "application/json"
        })

    def query(self, prompt: str, task_type: str = "text-generation", 
              provider: Optional[str] = None, model: Optional[str] = None) -> Dict[str, Any]:
        """
        Send a query to the OrcaAI platform.

        Args:
            prompt (str): The prompt to send to the AI
            task_type (str): The type of task (default: text-generation)
            provider (str, optional): Specific provider to use
            model (str, optional): Specific model to use

        Returns:
            Dict[str, Any]: The AI response

        Raises:
            OrcaAIException: If there's an error with the request
        """
        url = f"{self.base_url}/api/v1/ai/query"
        
        payload = {
            "prompt": prompt,
            "task_type": task_type
        }
        
        if provider:
            payload["provider"] = provider
            
        if model:
            payload["model"] = model

        try:
            response = self.session.post(url, json=payload)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.HTTPError as e:
            if response.status_code == 401:  # pyright: ignore[reportPossiblyUnboundVariable]
                raise AuthenticationError("Invalid API key")
            else:
                raise APIError(f"API request failed: {e}")
        except requests.exceptions.RequestException as e:
            raise OrcaAIException(f"Request failed: {e}")

    def get_providers(self) -> Dict[str, Any]:
        """
        Get available AI providers.

        Returns:
            Dict[str, Any]: List of available providers

        Raises:
            OrcaAIException: If there's an error with the request
        """
        url = f"{self.base_url}/api/v1/ai/providers"
        
        try:
            response = self.session.get(url)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.HTTPError as e:
            if response.status_code == 401:  # pyright: ignore[reportPossiblyUnboundVariable]
                raise AuthenticationError("Invalid API key")
            else:
                raise APIError(f"API request failed: {e}")
        except requests.exceptions.RequestException as e:
            raise OrcaAIException(f"Request failed: {e}")

    def get_metrics(self) -> Dict[str, Any]:
        """
        Get usage metrics.

        Returns:
            Dict[str, Any]: Usage metrics

        Raises:
            OrcaAIException: If there's an error with the request
        """
        url = f"{self.base_url}/api/v1/metrics"
        
        try:
            response = self.session.get(url)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.HTTPError as e:
            if response.status_code == 401:  # pyright: ignore[reportPossiblyUnboundVariable]
                raise AuthenticationError("Invalid API key")
            else:
                raise APIError(f"API request failed: {e}")
        except requests.exceptions.RequestException as e:
            raise OrcaAIException(f"Request failed: {e}")

    def get_api_keys(self) -> Dict[str, Any]:
        """
        Get all API keys for the authenticated user.

        Returns:
            Dict[str, Any]: List of API keys

        Raises:
            OrcaAIException: If there's an error with the request
        """
        url = f"{self.base_url}/api/v1/keys"
        
        try:
            response = self.session.get(url)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.HTTPError as e:
            if response.status_code == 401:# pyright: ignore[reportPossiblyUnboundVariable]
                raise AuthenticationError("Invalid API key")
            else:
                raise APIError(f"API request failed: {e}")
        except requests.exceptions.RequestException as e:
            raise OrcaAIException(f"Request failed: {e}")

    def create_api_key(self, name: str) -> Dict[str, Any]:
        """
        Create a new API key.

        Args:
            name (str): Name for the new API key

        Returns:
            Dict[str, Any]: The created API key

        Raises:
            OrcaAIException: If there's an error with the request
        """
        url = f"{self.base_url}/api/v1/keys"
        
        payload = {"name": name}
        
        try:
            response = self.session.post(url, json=payload)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.HTTPError as e:
            if response.status_code == 401:# pyright: ignore[reportPossiblyUnboundVariable]
                raise AuthenticationError("Invalid API key")
            else:
                raise APIError(f"API request failed: {e}")
        except requests.exceptions.RequestException as e:
            raise OrcaAIException(f"Request failed: {e}")

    def delete_api_key(self, key_id: str) -> Dict[str, Any]:
        """
        Delete an API key.

        Args:
            key_id (str): ID of the API key to delete

        Returns:
            Dict[str, Any]: Confirmation of deletion

        Raises:
            OrcaAIException: If there's an error with the request
        """
        url = f"{self.base_url}/api/v1/keys/{key_id}"
        
        try:
            response = self.session.delete(url)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.HTTPError as e:
            if response.status_code == 401:# pyright: ignore[reportPossiblyUnboundVariable]
                raise AuthenticationError("Invalid API key")
            else:
                raise APIError(f"API request failed: {e}")
        except requests.exceptions.RequestException as e:
            raise OrcaAIException(f"Request failed: {e}")