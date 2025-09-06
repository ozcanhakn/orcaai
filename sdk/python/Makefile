.PHONY: install test dist clean

# Install the package in development mode
install:
	pip install -e .

# Run tests
test:
	python -m pytest tests/ -v

# Run tests with coverage
test-cov:
	python -m pytest tests/ --cov=orcaai --cov-report=html

# Format code with black
format:
	black .

# Check code style with flake8
lint:
	flake8 .

# Build distribution packages
dist:
	python -m build

# Clean build artifacts
clean:
	rm -rf build/
	rm -rf dist/
	rm -rf *.egg-info/
	rm -rf htmlcov/
	rm -rf .pytest_cache/

# Install development dependencies
dev-install:
	pip install -r requirements-dev.txt

# Run all checks
check: lint test

# Run all checks with coverage
check-cov: lint test-cov