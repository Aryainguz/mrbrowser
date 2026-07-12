import { createFileRoute } from "@tanstack/react-router";
import { FadeIn } from "@/components/FadeIn";
import { CodeBlock } from "@/components/CodeBlock";
import { Callout } from "@/components/docs/Callout";
import { DocsPager } from "@/components/docs/DocsPager";

export const Route = createFileRoute("/docs/python")({
  head: () => ({
    meta: [
      { title: "Python SDK — Mr. Browser Docs" },
      {
        name: "description",
        content:
          "Use Mr. Browser from Python: pytest integration, intent-based actions, and semantic data extraction.",
      },
      { property: "og:title", content: "Python SDK — Mr. Browser Docs" },
      {
        property: "og:description",
        content: "Pytest-native, intent-driven browser automation in Python.",
      },
    ],
  }),
  component: PythonSdk,
});

const BASIC = `from mrbrowser import MrBrowser

with MrBrowser() as browser:
    page = browser.open("https://shop.example.com")
    page.type("the search field", "mechanical keyboard")
    page.click("the search button")
    page.click("the first product card")
    assert "Add to cart" in page.get_html()`;

const PYTEST = `# conftest.py
import pytest
from mrbrowser import MrBrowser

@pytest.fixture
def browser():
    with MrBrowser(host="localhost", port=7331) as b:
        yield b

# test_checkout.py
def test_checkout_survives_redesign(browser):
    page = browser.open("https://shop.example.com")
    page.click("Add to cart")
    page.click("the checkout button")
    assert page.assert_text("Order summary")
    price = page.extract_text("the order total")
    assert price.startswith("$")`;

const EXTRACT = `page = browser.open("https://billing.corp.com/invoices")

# Extract a specific field by plain-English description
balance = page.extract_text("the account balance")

# Inspect all elements on the page (useful for debugging)
elements = page.inspect(visible_only=True)
for el in elements:
    print(el["text"], el["role"], el["selector"])

# Cookie management
cookies = page.get_cookies()
page.set_cookies([{"name": "session", "value": "abc123"}])

# Run JS
result = page.execute_js("return document.title")`;

function PythonSdk() {
  return (
    <FadeIn>
      <article className="max-w-3xl">
        <p className="font-mono text-xs text-primary">SDKs</p>
        <h1 className="mt-2 font-mono text-3xl font-bold tracking-tight sm:text-4xl">
          Python SDK
        </h1>
        <p className="mt-4 text-[15px] leading-relaxed text-muted-foreground">
          <code className="font-mono text-primary">pip install mrbrowser</code> gives you
          the full engine with a Pythonic API — designed to feel native inside pytest.
        </p>

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Basic usage
        </h2>
        <CodeBlock className="mt-4" code={BASIC} lang="python" title="quickstart.py" />

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Pytest integration
        </h2>
        <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
          Share one Memory Engine database across your suite so fingerprints accumulate
          run over run — that history is what powers self-healing in CI.
        </p>
        <CodeBlock className="mt-4" code={PYTEST} lang="python" title="test_checkout.py" />

        <Callout type="tip">
          Commit <code className="font-mono text-primary">.mrbrowser/memory.db</code> to
          your repo (it's a compact SQLite file). Fresh CI runners then start with the
          full fingerprint history instead of a cold cache.
        </Callout>

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Semantic extraction
        </h2>
        <CodeBlock className="mt-4" code={EXTRACT} lang="python" title="extract.py" />

        <Callout type="warning">
          <code className="font-mono text-destructive">extract()</code> resolves intents
          per field, per row. On tables with 1000+ rows, pass{" "}
          <code className="font-mono text-destructive">batch=True</code> to fingerprint the
          row structure once — otherwise extraction time grows linearly.
        </Callout>

        <DocsPager path="/docs/python" />
      </article>
    </FadeIn>
  );
}
