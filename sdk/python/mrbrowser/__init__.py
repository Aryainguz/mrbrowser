"""
Mr. Browser Python SDK
======================

A clean Python client for the Mr. Browser automation engine.

Usage::

    from mrbrowser import MrBrowser

    browser = MrBrowser(host="localhost", port=7331)
    page = browser.open("https://example.com")
    page.click("login button")
    page.type("email field", "user@example.com")
    screenshot = page.screenshot()
    browser.close()
"""

from .client import MrBrowser
from .page import Page
from .exceptions import MrBrowserError, ConnectionError, ElementNotFoundError, ActionError

__version__ = "0.1.0"
__all__ = [
    "MrBrowser",
    "Page",
    "MrBrowserError",
    "ConnectionError",
    "ElementNotFoundError",
    "ActionError",
]
