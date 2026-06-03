export class MrBrowserError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'MrBrowserError';
  }
}

export class ConnectionError extends MrBrowserError {
  constructor(host: string, port: number, cause?: Error) {
    super(`Cannot connect to Mr. Browser at ${host}:${port}${cause ? ` — ${cause.message}` : ''}`);
    this.name = 'ConnectionError';
  }
}

export class ElementNotFoundError extends MrBrowserError {
  constructor(target: string, url?: string) {
    super(`Element not found: "${target}"${url ? ` on ${url}` : ''}`);
    this.name = 'ElementNotFoundError';
  }
}

export class ActionError extends MrBrowserError {
  constructor(action: string, target: string, reason: string) {
    super(`Action '${action}' on "${target}" failed: ${reason}`);
    this.name = 'ActionError';
  }
}
