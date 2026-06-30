# Python SDK Deep Dive

The Mr. Browser Python SDK (`mrbrowser`) provides a Pythonic, fluent API for orchestrating the engine. It is designed for maximum portability and developer experience.

## Design Philosophy

The primary design goal of the Python SDK is **Zero External Dependencies**. It is built entirely on the Python Standard Library (`urllib`, `json`, `dataclasses`). This means you can drop it into any legacy environment, AWS Lambda, or strict corporate CI/CD pipeline without needing to battle `pip` dependency conflicts or install massive binaries.

The SDK acts as a strict thin-client. All heavy lifting, DOM extraction, NLP resolution, and CDP orchestration occurs in the Go engine. The Python client simply manages session state and maps method calls to REST payloads.

---

## Installation

```bash
pip install mrbrowser
```

---

## Use Cases & Examples

### 1. Robust Data Extraction Scraper
Instead of hardcoding brittle `BeautifulSoup` or `lxml` selectors, use intent to find the data container.

```python
from mrbrowser import MrBrowser, MrBrowserError

def scrape_financials(ticker_url):
    with MrBrowser(host="localhost", port=7331) as browser:
        page = browser.open(ticker_url)
        
        try:
            # Wait for the SPA to load the financial data
            page.wait(3)
            
            # The engine will dynamically find the table based on semantic intent
            elements = page.inspect(visible_only=True)
            
            # Alternatively, extract directly if the API supports it
            # price = page.extract_text("Current Stock Price")
            
            page.screenshot(save_to=f"reports/{ticker_url}_snapshot.png")
            return elements
        except MrBrowserError as e:
            print(f"Extraction failed: {e}")
```

### 2. E2E Test Suite Integration (PyTest)
Mr. Browser integrates beautifully with `pytest`. 

```python
# test_login.py
import pytest
from mrbrowser import MrBrowser

@pytest.fixture(scope="module")
def browser():
    with MrBrowser() as b:
        yield b

def test_user_can_login(browser):
    page = browser.open("https://example.com/login")
    
    page.type("Email Field", "user@test.com")
    page.type("Password Input", "secure123")
    page.click("Login")
    
    # Assert state mutation
    page.wait(1)
    page.assert_text("Welcome back, User!")
```

---

## Complete API Reference

### `MrBrowser` Client

#### `__init__(host="localhost", port=7331, timeout=30)`
Initializes the client. Does *not* connect immediately.
- `timeout` (int): Maximum seconds to wait for any single HTTP request to the engine.

#### `ping() -> bool`
Checks if the engine is alive. Useful for health checks before starting a massive scrape job.

#### `open(url: str) -> Page`
Requests a new isolated browsing context from the engine and navigates to the URL. Returns a stateful `Page` object.

#### `run_workflow(yaml_path: str) -> dict`
Bypasses the programmatic API and tells the engine to execute a complete YAML workflow file directly. Excellent for kicking off batch jobs where the steps are defined by non-developers.

### `Page` Context

The `Page` object represents a single browser tab. Methods can be chained when applicable.

#### `click(target: str, selector: str = None) -> 'Page'`
Resolves the `target` string to a DOM element and triggers a physical mouse click. 
- **Under the hood:** The engine verifies the click caused a DOM mutation (e.g., URL change or HTML structure change).

#### `type(target: str, value: str, enter: bool = False) -> 'Page'`
Types text into an input field.
- `enter` (bool): If True, appends a carriage return to submit the form immediately.

#### `screenshot(save_to: str) -> bytes`
Takes a full-page snapshot and writes it to disk. Returns the raw PNG bytes.

#### `inspect(visible_only: bool = False) -> list[dict]`
Returns the raw JSON accessibility tree extracted by the engine. Useful for debugging why an intent string isn't matching, or for dumping structured data.
