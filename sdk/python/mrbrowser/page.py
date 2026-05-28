"""
Page represents a single browser tab in Mr. Browser.
"""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, Optional

if TYPE_CHECKING:
    from .client import MrBrowser


class Page:
    """
    Represents a single browser tab/page.

    All navigation, element interaction, and data extraction happen through a Page.

    Example::

        page = browser.open("https://example.com")
        page.click("login button")
        page.type("email field", "user@example.com")
        page.type("password field", "secret")
        page.click("sign in")
        page.screenshot(save_to="confirmation.png")
    """

    def __init__(self, client: "MrBrowser", session_id: str, url: str = ""):
        self._client = client
        self._session_id = session_id
        self._url = url

    @property
    def session_id(self) -> str:
        """The session ID for this page."""
        return self._session_id

    # ──────────────────────────────────────────────────────────
    # Navigation
    # ──────────────────────────────────────────────────────────

    def navigate(self, url: str) -> "Page":
        """Navigate to a URL."""
        self._client._post(f"/sessions/{self._session_id}/navigate", {"url": url})
        self._url = url
        return self

    @property
    def url(self) -> str:
        """Return the current page URL."""
        resp = self._client._get(f"/sessions/{self._session_id}/url")
        return resp.get("url", self._url)

    @property
    def title(self) -> str:
        """Return the current page title."""
        resp = self._client._get(f"/sessions/{self._session_id}/title")
        return resp.get("title", "")

    def reload(self) -> "Page":
        """Reload the current page."""
        self._client._post(f"/sessions/{self._session_id}/reload", {})
        return self

    # ──────────────────────────────────────────────────────────
    # Intent-based actions
    # ──────────────────────────────────────────────────────────

    def click(self, target: str, *, selector: str = "") -> "Page":
        """
        Click an element identified by intent.

        Args:
            target: Natural-language description of the element.
                    Examples: "login button", "sign up link", "accept all"
            selector: Optional explicit CSS selector (bypasses resolver).

        Returns:
            self (for chaining)

        Example::

            page.click("login button")
            page.click("accept cookies")
        """
        self._client._post(f"/sessions/{self._session_id}/click", {
            "target": target,
            "selector": selector,
        })
        return self

    def type(self, target: str, value: str, *, selector: str = "", clear: bool = True) -> "Page":
        """
        Type text into an element identified by intent.

        Args:
            target: Natural-language description of the input field.
            value: The text to type.
            selector: Optional explicit CSS selector.
            clear: Clear the field before typing (default: True).

        Example::

            page.type("email field", "user@example.com")
            page.type("password", "secret123")
        """
        self._client._post(f"/sessions/{self._session_id}/type", {
            "target": target,
            "value": value,
            "selector": selector,
            "clear": clear,
        })
        return self

    def hover(self, target: str, *, selector: str = "") -> "Page":
        """Hover over an element."""
        self._client._post(f"/sessions/{self._session_id}/hover", {
            "target": target,
            "selector": selector,
        })
        return self

    def scroll(self, direction: str = "down", pixels: int = 300) -> "Page":
        """
        Scroll the page.

        Args:
            direction: "up", "down", "left", "right", "top", "bottom"
            pixels: Pixels to scroll (ignored for top/bottom).
        """
        self._client._post(f"/sessions/{self._session_id}/scroll", {
            "direction": direction,
            "pixels": pixels,
        })
        return self

    def upload(self, target: str, file_path: str, *, selector: str = "") -> "Page":
        """Upload a file to a file input element."""
        self._client._post(f"/sessions/{self._session_id}/upload", {
            "target": target,
            "file_path": file_path,
            "selector": selector,
        })
        return self

    # ──────────────────────────────────────────────────────────
    # Data extraction
    # ──────────────────────────────────────────────────────────

    def screenshot(self, save_to: Optional[str] = None) -> bytes:
        """
        Capture a PNG screenshot.

        Args:
            save_to: Optional file path to save the screenshot.

        Returns:
            PNG image bytes.
        """
        resp = self._client._post(f"/sessions/{self._session_id}/screenshot", {})
        png_bytes = bytes(resp.get("data", []))
        if save_to:
            with open(save_to, "wb") as f:
                f.write(png_bytes)
        return png_bytes

    def get_html(self) -> str:
        """Return the full HTML of the current page."""
        resp = self._client._get(f"/sessions/{self._session_id}/html")
        return resp.get("html", "")

    def execute_js(self, script: str) -> Any:
        """Execute JavaScript in the page context and return the result."""
        resp = self._client._post(f"/sessions/{self._session_id}/js", {"script": script})
        return resp.get("result")

    def extract_text(self, target: str) -> str:
        """
        Extract visible text from an element identified by intent.

        Returns:
            The text content of the matched element.
        """
        resp = self._client._post(f"/sessions/{self._session_id}/extract", {
            "target": target,
        })
        return resp.get("text", "")

    def inspect(self, visible_only: bool = False) -> list[dict]:
        """
        Return all detected elements on the page.

        Returns:
            List of element dicts with keys: tag, type, text, role, visible, position, selector.
        """
        params = "?visible_only=true" if visible_only else ""
        resp = self._client._get(f"/sessions/{self._session_id}/elements{params}")
        return resp.get("elements", [])

    # ──────────────────────────────────────────────────────────
    # Waiting
    # ──────────────────────────────────────────────────────────

    def wait_for_selector(self, selector: str, timeout: float = 10.0) -> "Page":
        """Wait until a CSS selector is visible."""
        self._client._post(f"/sessions/{self._session_id}/wait", {
            "selector": selector,
            "timeout": timeout,
        })
        return self

    def wait_for_url(self, url_contains: str, timeout: float = 10.0) -> "Page":
        """Wait until the page URL contains the given string."""
        self._client._post(f"/sessions/{self._session_id}/wait", {
            "url": url_contains,
            "timeout": timeout,
        })
        return self

    def wait(self, seconds: float) -> "Page":
        """Wait for a fixed duration."""
        import time
        time.sleep(seconds)
        return self

    # ──────────────────────────────────────────────────────────
    # Assertions
    # ──────────────────────────────────────────────────────────

    def assert_text(self, text: str) -> "Page":
        """Assert that text is present on the page. Raises AssertionError if not."""
        html = self.get_html()
        if text not in html:
            raise AssertionError(f"Text not found on page: {text!r}")
        return self

    def assert_url_contains(self, substring: str) -> "Page":
        """Assert that the current URL contains the given substring."""
        current = self.url
        if substring not in current:
            raise AssertionError(f"URL does not contain {substring!r}: got {current!r}")
        return self

    # ──────────────────────────────────────────────────────────
    # Cookies
    # ──────────────────────────────────────────────────────────

    def get_cookies(self) -> list[dict]:
        """Return all cookies for this page."""
        resp = self._client._get(f"/sessions/{self._session_id}/cookies")
        return resp.get("cookies", [])

    def set_cookies(self, cookies: list[dict]) -> "Page":
        """Set cookies on this page."""
        self._client._post(f"/sessions/{self._session_id}/cookies", {"cookies": cookies})
        return self

    def close(self) -> None:
        """Close this page/tab."""
        self._client._delete(f"/sessions/{self._session_id}")

    def __repr__(self) -> str:
        return f"Page(session={self._session_id!r}, url={self._url!r})"

    # ──────────────────────────────────────────────────────────
    # Context manager support
    # ──────────────────────────────────────────────────────────

    def __enter__(self) -> "Page":
        return self

    def __exit__(self, *args) -> None:
        self.close()
