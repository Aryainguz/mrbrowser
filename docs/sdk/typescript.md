# TypeScript SDK Deep Dive

The Mr. Browser TypeScript SDK (`@mrbrowser/sdk`) brings intent-driven automation to the Node.js ecosystem. 

## Design Philosophy

Similar to the Python SDK, the TypeScript SDK is designed with **Zero Dependencies**, utilizing the native `fetch` API available in Node 18+. 

It provides **strict type-safety** using modern ES modules. Every response, payload, and configuration option is strongly typed, providing excellent autocomplete in VSCode and catching errors at compile time rather than runtime.

---

## Installation

Requires Node.js >= 18.0.0.

```bash
npm install @mrbrowser/sdk
```

---

## Use Cases & Examples

### 1. Serverless Web Automation (AWS Lambda / Vercel)
Because the SDK is extremely lightweight and uses native fetch, it is perfectly suited for serverless functions that need to orchestrate a remote Mr. Browser engine (hosted on a long-running EC2 instance or Fargate task).

```typescript
// pages/api/scrape.ts
import { MrBrowser, ElementNotFoundError } from '@mrbrowser/sdk';

export default async function handler(req, res) {
  // Connect to remote engine
  const browser = new MrBrowser({ 
    host: 'engine.internal.corp', 
    port: 7331,
    apiKey: process.env.MR_BROWSER_KEY 
  });
  
  try {
    const page = await browser.open('https://dashboard.example.com');
    await page.click('Generate Report');
    await page.wait(5);
    
    const screenshotBuf = await page.screenshot();
    
    res.setHeader('Content-Type', 'image/png');
    res.status(200).send(screenshotBuf);
  } catch (error) {
    if (error instanceof ElementNotFoundError) {
      res.status(404).json({ error: 'Report button not found' });
    }
  } finally {
    await browser.close(); // Clean up remote sessions
  }
}
```

### 2. Integration with Jest / Mocha
Use Mr. Browser to write highly resilient tests in TypeScript.

```typescript
import { MrBrowser } from '@mrbrowser/sdk';

describe('Authentication Flow', () => {
  let browser: MrBrowser;

  beforeAll(() => {
    browser = new MrBrowser();
  });

  afterAll(async () => {
    await browser.close();
  });

  it('should successfully log in', async () => {
    const page = await browser.open('http://localhost:3000/login');
    
    // Fluent, awaitable API
    await page.type('Username', 'admin');
    await page.type('Password', 'secret123', { clear: true });
    await page.click('Login');
    
    await page.assertText('Welcome to the Dashboard');
  }, 30000); // Extended timeout for browser ops
});
```

---

## Complete API Reference

### `MrBrowser` Client

#### `constructor(options: MrBrowserOptions)`
- `options.host`: Engine hostname (default: `localhost`)
- `options.port`: Engine port (default: `7331`)
- `options.timeout`: Request timeout in ms (default: `30000`)
- `options.apiKey`: Optional API key for authenticated engines.

#### `ping(): Promise<boolean>`
Returns `true` if the engine is reachable and healthy.

#### `open(url: string): Promise<Page>`
Spawns a new tab and navigates to the URL.

### `Page` Context

#### `click(target: string, options?: { selector?: string }): Promise<Page>`
Instructs the engine to resolve the target and click it. Chainable.

#### `type(target: string, value: string, options?: { selector?: string, clear?: boolean, enter?: boolean }): Promise<Page>`
Resolves an input field and types the value.
- `options.clear`: Whether to clear the field first (default: `true`).
- `options.enter`: Whether to trigger a form submit (`\r`) after typing.

#### `scroll(direction: 'up' | 'down' | 'left' | 'right', pixels?: number): Promise<Page>`
Scrolls the viewport.

#### `assertText(text: string): Promise<Page>`
Throws an error if the raw HTML of the page does not contain the specified text. Useful for quick validation.

#### `inspect(visibleOnly?: boolean): Promise<any[]>`
Returns an array of all interactable elements identified by the intelligence module.

#### `close(): Promise<void>`
Terminates the specific page session on the engine, freeing up memory.
