# Mr. Browser Python SDK

The Mr. Browser Python SDK (`mrbrowser`) provides a Pythonic, fluent API for orchestrating the Mr. Browser automation engine. It is designed for maximum portability and developer experience, with **Zero External Dependencies**.

## Installation

Currently, the SDK is available directly from the repository. We recommend installing it locally from source.

### Local Installation

Clone the repository and install the package using `pip`:

```bash
git clone https://github.com/aryainguz/mrbrowser.git
cd mrbrowser/sdk/python
pip install .
```

### Publishing to PyPI

If you wish to publish this package internally or to PyPI, you can build and upload it using `twine`:

```bash
pip install build twine
python -m build
twine upload dist/*
```

## Quick Start

```python
from mrbrowser import MrBrowser

# Example Usage
with MrBrowser(host="localhost", port=7331) as browser:
    page = browser.open("https://example.com/login")
    page.type("Email Field", "user@test.com")
    page.type("Password Input", "secure123")
    page.click("Login")
```

For more details, see the full [documentation](https://github.com/aryainguz/mrbrowser/tree/main/docs/sdk/python.md).
