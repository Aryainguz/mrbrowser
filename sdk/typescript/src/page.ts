import { MrBrowser } from './client';

export class Page {
  private client: MrBrowser;
  private sessionId: string;
  private _url: string;

  constructor(client: MrBrowser, sessionId: string, url: string = '') {
    this.client = client;
    this.sessionId = sessionId;
    this._url = url;
  }

  get id(): string {
    return this.sessionId;
  }

  // ──────────────────────────────────────────────────────────
  // Navigation
  // ──────────────────────────────────────────────────────────

  async navigate(url: string): Promise<Page> {
    await this.client._post(`/sessions/${this.sessionId}/navigate`, { url });
    this._url = url;
    return this;
  }

  async url(): Promise<string> {
    const resp = await this.client._get(`/sessions/${this.sessionId}/url`);
    return resp.url || this._url;
  }

  async title(): Promise<string> {
    const resp = await this.client._get(`/sessions/${this.sessionId}/title`);
    return resp.title || '';
  }

  async reload(): Promise<Page> {
    await this.client._post(`/sessions/${this.sessionId}/reload`, {});
    return this;
  }

  // ──────────────────────────────────────────────────────────
  // Intent-based actions
  // ──────────────────────────────────────────────────────────

  async click(target: string, options?: { selector?: string }): Promise<Page> {
    await this.client._post(`/sessions/${this.sessionId}/click`, {
      target,
      selector: options?.selector || '',
    });
    return this;
  }

  async type(target: string, value: string, options?: { selector?: string; clear?: boolean }): Promise<Page> {
    await this.client._post(`/sessions/${this.sessionId}/type`, {
      target,
      value,
      selector: options?.selector || '',
      clear: options?.clear ?? true,
    });
    return this;
  }

  async hover(target: string, options?: { selector?: string }): Promise<Page> {
    await this.client._post(`/sessions/${this.sessionId}/hover`, {
      target,
      selector: options?.selector || '',
    });
    return this;
  }

  async scroll(direction: 'up' | 'down' | 'left' | 'right' | 'top' | 'bottom' = 'down', pixels: number = 300): Promise<Page> {
    await this.client._post(`/sessions/${this.sessionId}/scroll`, {
      direction,
      pixels,
    });
    return this;
  }

  // ──────────────────────────────────────────────────────────
  // Data extraction
  // ──────────────────────────────────────────────────────────

  async screenshot(): Promise<Buffer> {
    const resp = await this.client._post(`/sessions/${this.sessionId}/screenshot`, {});
    return Buffer.from(resp.data, 'base64');
  }

  async getHtml(): Promise<string> {
    const resp = await this.client._get(`/sessions/${this.sessionId}/html`);
    return resp.html || '';
  }

  async executeJs(script: string): Promise<any> {
    const resp = await this.client._post(`/sessions/${this.sessionId}/js`, { script });
    return resp.result;
  }

  async inspect(visibleOnly: boolean = false): Promise<any[]> {
    const params = visibleOnly ? '?visible_only=true' : '';
    const resp = await this.client._get(`/sessions/${this.sessionId}/elements${params}`);
    return resp.elements || [];
  }

  // ──────────────────────────────────────────────────────────
  // Waiting
  // ──────────────────────────────────────────────────────────

  async waitForSelector(selector: string, timeout: number = 10.0): Promise<Page> {
    await this.client._post(`/sessions/${this.sessionId}/wait`, {
      selector,
      timeout,
    });
    return this;
  }

  async wait(seconds: number): Promise<Page> {
    await new Promise((resolve) => setTimeout(resolve, seconds * 1000));
    return this;
  }

  // ──────────────────────────────────────────────────────────
  // Assertions
  // ──────────────────────────────────────────────────────────

  async assertText(text: string): Promise<Page> {
    const html = await this.getHtml();
    if (!html.includes(text)) {
      throw new Error(`Assertion failed: Text not found on page: "${text}"`);
    }
    return this;
  }

  async close(): Promise<void> {
    await this.client._delete(`/sessions/${this.sessionId}`);
  }
}
