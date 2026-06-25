"""
Exceptions raised by the Mr. Browser Python SDK.
"""


class MrBrowserError(Exception):
    """Base exception for all Mr. Browser errors."""
    pass


class ConnectionError(MrBrowserError):
    """Raised when the SDK cannot connect to the Mr. Browser engine."""

    def __init__(self, host: str, port: int, cause: Exception = None):
        msg = f"Cannot connect to Mr. Browser at {host}:{port}"
        if cause:
            msg += f" — {cause}"
        super().__init__(msg)
        self.host = host
        self.port = port
        self.cause = cause


class ElementNotFoundError(MrBrowserError):
    """Raised when no element matching the intent could be found."""

    def __init__(self, target: str, url: str = None):
        msg = f"Element not found: {target!r}"
        if url:
            msg += f" on {url}"
        super().__init__(msg)
        self.target = target
        self.url = url


class ActionError(MrBrowserError):
    """Raised when a browser action fails."""

    def __init__(self, action: str, target: str, reason: str):
        super().__init__(f"Action '{action}' on {target!r} failed: {reason}")
        self.action = action
        self.target = target
        self.reason = reason


class TaskError(MrBrowserError):
    """Raised when a task workflow fails."""

    def __init__(self, task_name: str, step: int, reason: str):
        super().__init__(f"Task '{task_name}' failed at step {step}: {reason}")
        self.task_name = task_name
        self.step = step
        self.reason = reason
