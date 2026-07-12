import { createFileRoute } from "@tanstack/react-router";
import { FadeIn } from "@/components/FadeIn";
import { CodeBlock } from "@/components/CodeBlock";
import { Callout } from "@/components/docs/Callout";
import { DocsPager } from "@/components/docs/DocsPager";

export const Route = createFileRoute("/docs/typescript")({
  head: () => ({
    meta: [
      { title: "TypeScript SDK — Mr. Browser Docs" },
      {
        name: "description",
        content:
          "Use Mr. Browser from TypeScript: fully typed intent actions, async flows, and typed semantic extraction.",
      },
      { property: "og:title", content: "TypeScript SDK — Mr. Browser Docs" },
      {
        property: "og:description",
        content: "Fully typed, intent-driven browser automation for Node and Bun.",
      },
    ],
  }),
  component: TypeScriptSdk,
});

const BASIC = `import { MrBrowser } from "@mrbrowser/sdk";

const browser = new MrBrowser({ host: "localhost", port: 7331 });
const page = await browser.open("https://shop.example.com");

await page.type("the search field", "mechanical keyboard");
await page.click("the search button");
await page.click("the first product card");

const html = await page.getHtml();
if (!html.includes("Add to cart")) {
  throw new Error("product page did not load");
}
await page.close();
await browser.close();`;

const TYPED_EXTRACT = `import { MrBrowser } from "@mrbrowser/sdk";

const browser = new MrBrowser();
const page = await browser.open("https://billing.corp.com/invoices");

// Extract text by plain-English description
const balance = await page.extractText("the account balance");

// Inspect all interactive elements
const elements = await page.inspect({ visibleOnly: true });
elements.forEach(el => console.log(el.text, el.role, el.selector));

// Run arbitrary JS
const title = await page.executeJs<string>("return document.title");

// Cookies
const cookies = await page.getCookies();
await page.setCookies([{ name: "session", value: "abc123" }]);

await browser.close();`;

const CONFIG = `import { MrBrowser } from "@mrbrowser/sdk";

const browser = new MrBrowser({
  host: "localhost",
  port: 7331,          // default engine port
  timeout: 30_000,     // ms per request (default 30s)
  apiKey: process.env.MR_BROWSER_KEY, // optional auth
});

// Available Page methods:
// page.click(target)             — click by plain-English description
// page.type(target, value)       — type into a field
// page.hover(target)             — hover over an element
// page.scroll(direction, pixels) — scroll: up/down/left/right/top/bottom
// page.navigate(url)             — navigate to a new URL
// page.reload()                  — reload current page
// page.screenshot()              — returns PNG Buffer
// page.getHtml()                 — returns full page HTML
// page.extractText(target)       — extract text by intent
// page.inspect(opts)             — list all elements
// page.waitForSelector(css)      — wait for element visibility
// page.waitForUrl(substring)     — wait for URL change
// page.assertText(text)          — assert text is in page HTML
// page.getCookies()              — get all cookies
// page.setCookies(cookies)       — set cookies
// page.executeJs(script)         — execute JavaScript
// page.close()                   — close this tab`;

function TypeScriptSdk() {
  return (
    <FadeIn>
      <article className="max-w-3xl">
        <p className="font-mono text-xs text-primary">SDKs</p>
        <h1 className="mt-2 font-mono text-3xl font-bold tracking-tight sm:text-4xl">
          TypeScript SDK
        </h1>
        <p className="mt-4 text-[15px] leading-relaxed text-muted-foreground">
          <code className="font-mono text-primary">npm install @mrbrowser/sdk</code> — a
          fully typed client for Node 18+. Instantiate{" "}
          <code className="font-mono text-primary">new MrBrowser(&#123;host, port&#125;)</code>,
          open a page, and interact with elements by plain-English description.
        </p>

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Basic usage
        </h2>
        <CodeBlock className="mt-4" code={BASIC} lang="typescript" title="quickstart.ts" />

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Typed extraction
        </h2>
        <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
          Pass a type parameter to <code className="font-mono text-primary">extract</code>{" "}
          and the result is typed end to end — your scraper breaks at compile time, not in
          production.
        </p>
        <CodeBlock className="mt-4" code={TYPED_EXTRACT} lang="typescript" title="invoices.ts" />

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Configuration & full API
        </h2>
        <CodeBlock className="mt-4" code={CONFIG} lang="typescript" title="api-reference.ts" />

        <Callout type="tip">
          The SDK requires the Mr. Browser engine to be running. Start it with{" "}
          <code className="font-mono text-primary">docker compose -f docker/docker-compose.yml up -d</code>{" "}
          or run <code className="font-mono text-primary">mr-browser --help</code> if installed from source.
          The engine listens on <code className="font-mono text-primary">localhost:7331</code> by default.
        </Callout>

        <Callout type="warning">
          If the engine is unreachable you will get{" "}
          <code className="font-mono text-destructive">ECONNREFUSED 127.0.0.1:7331</code>.
          Make sure the engine is started before calling{" "}
          <code className="font-mono text-destructive">browser.open()</code>.
        </Callout>

        <DocsPager path="/docs/typescript" />
      </article>
    </FadeIn>
  );
}
