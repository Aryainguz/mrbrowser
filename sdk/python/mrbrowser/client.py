"""
MrBrowser client — the main entry point for the Python SDK.
"""

from __future__ import annotations

import json
from typing import Any, Dict, Optional
from urllib.request import urlopen, Request
from urllib.error import URLError
from urllib.parse import urlencode

from .page import Page
from .exceptions import ConnectionError, ElementNotFoundError, ActionError, MrBrowserError


class MrBrowser:
    """
    Client for the Mr. Browser automation engine.

    Connects to a running Mr. Browser server (default: localhost:7331).

    Example::

        browser = MrBrowser(host="localhost", port=7331)
        page = browser.open("https://example.com")
        page.click("login button")
        page.type("email field", "user@example.com")
        browser.close()

    Context manager::

        with MrBrowser() as browser:
            page = browser.open("https://example.com")
            page.click("login button")
    """

    DEFAULT_HOST = "localhost"
    DEFAULT_PORT = 7331
    DEFAULT_TIMEOUT = 30.0

    def __init__(
        self,
        host: str = DEFAULT_HOST,
        port: int = DEFAULT_PORT,
        timeout: float = DEFAULT_TIMEOUT,
        api_key: Optional[str] = None,
    ):
        """
        Create a new MrBrowser client.

        Args:
            host: Hostname of the Mr. Browser server.
            port: Port of the Mr. Browser server.
            timeout: Request timeout in seconds.
            api_key: Optional API key for authentication.
        """
        self.host = host
        self.port = port
        self.timeout = timeout
        self._api_key = api_key
        self._base_url = f"http://{host}:{port}/api/v1"
        self._pages: list[Page] = []

    def ping(self) -> bool:
        """Check whether the Mr. Browser server is reachable."""
        try:
            resp = self._get("/health")
            return resp.get("status") == "ok"
        except MrBrowserError:
            return False

    def open(self, url: str) -> Page:
        """
        Open a new browser tab and navigate to the given URL.

        Args:
            url: The URL to navigate to.

        Returns:
            A :class:`Page` object for the new tab.

        Example::

            page = browser.open("https://example.com")
        """
        resp = self._post("/sessions", {"url": url})
        session_id = resp.get("session_id", "")
        page = Page(self, session_id, url)
        self._pages.append(page)
        return page

    def run_workflow(self, yaml_path: str) -> Dict[str, Any]:
        """
        Execute a YAML workflow file.

        Args:
            yaml_path: Path to the YAML workflow file.

        Returns:
            A dict with keys: task_name, success, step_results, duration.
        """
        with open(yaml_path) as f:
            content = f.read()
        return self._post("/workflows/run", {"yaml": content})

    def close(self) -> None:
        """Close all open pages and the browser session."""
        self._post("/sessions/close-all", {})
        self._pages.clear()

    def version(self) -> str:
        """Return the Mr. Browser server version."""
        resp = self._get("/version")
        return resp.get("version", "unknown")

    # ──────────────────────────────────────────────────────────
    # Internal HTTP helpers
    # ──────────────────────────────────────────────────────────

    def _get(self, path: str) -> Dict[str, Any]:
        """Make a GET request to the API."""
        url = self._base_url + path
        req = Request(url, method="GET")
        self._add_headers(req)
        return self._do(req)

    def _post(self, path: str, body: Dict[str, Any]) -> Dict[str, Any]:
        """Make a POST request to the API."""
        url = self._base_url + path
        data = json.dumps(body).encode("utf-8")
        req = Request(url, data=data, method="POST")
        req.add_header("Content-Type", "application/json")
        self._add_headers(req)
        return self._do(req)

    def _delete(self, path: str) -> Dict[str, Any]:
        """Make a DELETE request to the API."""
        url = self._base_url + path
        req = Request(url, method="DELETE")
        self._add_headers(req)
        return self._do(req)

    def _add_headers(self, req: Request) -> None:
        req.add_header("Accept", "application/json")
        if self._api_key:
            req.add_header("X-API-Key", self._api_key)

    def _do(self, req: Request) -> Dict[str, Any]:
        """Execute an HTTP request and parse the JSON response."""
        try:
            with urlopen(req, timeout=self.timeout) as resp:
                raw = resp.read().decode("utf-8")
                if not raw:
                    return {}
                data = json.loads(raw)
                if data.get("error"):
                    self._raise_api_error(data)
                return data
        except URLError as e:
            raise ConnectionError(self.host, self.port, e) from e
        except json.JSONDecodeError as e:
            raise MrBrowserError(f"Invalid JSON response: {e}") from e

    @staticmethod
    def _raise_api_error(data: Dict[str, Any]) -> None:
        """Convert an API error response into the appropriate SDK exception."""
        err = data.get("error", {})
        code = err.get("code", "unknown")
        msg = err.get("message", str(err))

        if code == "ELEMENT_NOT_FOUND":
            raise ElementNotFoundError(err.get("target", "?"))
        if code == "ACTION_FAILED":
            raise ActionError(
                err.get("action", "?"),
                err.get("target", "?"),
                msg,
            )
        raise MrBrowserError(f"API error [{code}]: {msg}")

    # ──────────────────────────────────────────────────────────
    # Context manager
    # ──────────────────────────────────────────────────────────

    def __enter__(self) -> "MrBrowser":
        return self

    def __exit__(self, *args) -> None:
        self.close()

    def __repr__(self) -> str:
        return f"MrBrowser(host={self.host!r}, port={self.port})"
