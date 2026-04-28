import { Page } from './page';
import { MrBrowserError, ConnectionError, ElementNotFoundError, ActionError } from './exceptions';

export interface MrBrowserOptions {
  host?: string;
  port?: number;
  timeout?: number;
  apiKey?: string;
}

export class MrBrowser {
  private host: string;
  private port: number;
  private timeout: number;
  private apiKey?: string;
  private baseUrl: string;
  private pages: Page[] = [];

  constructor(options: MrBrowserOptions = {}) {
    this.host = options.host || 'localhost';
    this.port = options.port || 7331;
    this.timeout = options.timeout || 30000;
    this.apiKey = options.apiKey;
    this.baseUrl = `http://${this.host}:${this.port}/api/v1`;
  }

  async ping(): Promise<boolean> {
    try {
      const resp = await this._get('/health', true);
      return resp.status === 'ok';
    } catch {
      return false;
    }
  }

  async open(url: string): Promise<Page> {
    const resp = await this._post('/sessions', { url });
    const sessionId = resp.session_id || '';
    const page = new Page(this, sessionId, url);
    this.pages.push(page);
    return page;
  }

  async runWorkflow(yamlContent: string): Promise<any> {
    return this._post('/workflows/run', { yaml: yamlContent });
  }

  async close(): Promise<void> {
    await this._post('/sessions/close-all', {});
    this.pages = [];
  }

  async version(): Promise<string> {
    const resp = await this._get('/version', true);
    return resp.version || 'unknown';
  }

  // ──────────────────────────────────────────────────────────
  // Internal HTTP helpers
  // ──────────────────────────────────────────────────────────

  async _get(path: string, bypassBase: boolean = false): Promise<any> {
    const url = bypassBase ? `http://${this.host}:${this.port}${path}` : `${this.baseUrl}${path}`;
    return this._do(url, { method: 'GET' });
  }

  async _post(path: string, body: any): Promise<any> {
    return this._do(`${this.baseUrl}${path}`, {
      method: 'POST',
      body: JSON.stringify(body),
    });
  }

  async _delete(path: string): Promise<any> {
    return this._do(`${this.baseUrl}${path}`, { method: 'DELETE' });
  }

  private async _do(url: string, init: RequestInit): Promise<any> {
    const headers: Record<string, string> = {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
    };
    if (this.apiKey) {
      headers['X-API-Key'] = this.apiKey;
    }

    const controller = new AbortController();
    const id = setTimeout(() => controller.abort(), this.timeout);

    try {
      const response = await fetch(url, { ...init, headers, signal: controller.signal });
      const raw = await response.text();
      
      if (!raw) return {};

      let data;
      try {
        data = JSON.parse(raw);
      } catch (e) {
        throw new MrBrowserError(`Invalid JSON response from server`);
      }

      if (data.error) {
        this._raiseApiError(data.error);
      }

      return data;
    } catch (error: any) {
      if (error.name === 'AbortError') {
        throw new ConnectionError(this.host, this.port, new Error('Request timed out'));
      }
      if (error instanceof MrBrowserError) {
        throw error;
      }
      throw new ConnectionError(this.host, this.port, error);
    } finally {
      clearTimeout(id);
    }
  }

  private _raiseApiError(err: any): void {
    const code = err.code || 'unknown';
    const msg = err.message || JSON.stringify(err);

    if (code === 'ELEMENT_NOT_FOUND') {
      throw new ElementNotFoundError(err.target || '?');
    }
    if (code === 'ACTION_FAILED') {
      throw new ActionError(err.action || '?', err.target || '?', msg);
    }
    throw new MrBrowserError(`API error [${code}]: ${msg}`);
  }
}
