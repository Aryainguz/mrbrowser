# Mr. Browser TypeScript SDK

The Mr. Browser TypeScript SDK (`@mrbrowser/sdk`) brings intent-driven automation to the Node.js ecosystem. It provides **strict type-safety** and has **Zero Dependencies**, utilizing the native `fetch` API.

## Installation

Requires Node.js >= 18.0.0.

Currently, the SDK is available directly from the repository. We recommend installing it locally from source.

### Local Installation

Clone the repository and build the package:

```bash
git clone https://github.com/aryainguz/mrbrowser.git
cd mrbrowser/sdk/typescript
npm install
npm run build
```

Then, you can link it to your project:
```bash
npm link
# Inside your project folder
npm link @mrbrowser/sdk
```

### Publishing to NPM

If you want to publish the package to NPM, ensure you are logged in and run:

```bash
npm publish --access public
```

## Quick Start

```typescript
import { MrBrowser } from '@mrbrowser/sdk';

const browser = new MrBrowser();
const page = await browser.open('http://localhost:3000/login');

await page.type('Username', 'admin');
await page.type('Password', 'secret123', { clear: true });
await page.click('Login');

await browser.close();
```

For more details, see the full [documentation](https://github.com/aryainguz/mrbrowser/tree/main/docs/sdk/typescript.md).
