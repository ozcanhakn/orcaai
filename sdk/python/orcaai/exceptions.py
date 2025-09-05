class OrcaAIException(Exception):
    """Base exception for OrcaAI SDK."""
    pass


class AuthenticationError(OrcaAIException):
    """Exception raised for authentication errors."""
    pass


class APIError(OrcaAIException):
    """Exception raised for API errors."""
    pass