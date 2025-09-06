import sys
import os

# Add the orcaai package to the path
sys.path.insert(0, os.path.join(os.path.dirname(__file__), ".."))

def test_import():
    """Test that we can import the OrcaClient."""
    from orcaai import OrcaClient
    assert OrcaClient is not None

def test_exceptions_import():
    """Test that we can import the exceptions."""
    from orcaai import OrcaAIException, AuthenticationError, APIError
    assert OrcaAIException is not None
    assert AuthenticationError is not None
    assert APIError is not None

if __name__ == "__main__":
    test_import()
    test_exceptions_import()
    print("All import tests passed!")