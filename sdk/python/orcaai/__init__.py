"""
OrcaAI Python SDK
~~~~~~~~~~~~~~~~~

A Python SDK for the OrcaAI platform.

Basic usage:

    from orcaai import OrcaClient

    client = OrcaClient(api_key="your-api-key")
    result = client.query("Summarize this report")
    print(result)

:copyright: (c) 2024 by OrcaAI Team.
:license: MIT, see LICENSE for more details.
"""

__title__ = "orcaai"
__version__ = "1.0.0"
__build__ = 100
__author__ = "OrcaAI Team"
__license__ = "MIT"
__copyright__ = "Copyright 2024 OrcaAI Team"

from .client import OrcaClient
from .exceptions import OrcaAIException, AuthenticationError, APIError

__all__ = [
    "OrcaClient",
    "OrcaAIException",
    "AuthenticationError",
    "APIError",
]