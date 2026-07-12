import { createFileRoute } from "@tanstack/react-router";
import { FadeIn } from "@/components/FadeIn";
import { CodeBlock } from "@/components/CodeBlock";
import { Callout } from "@/components/docs/Callout";
import { DocsPager } from "@/components/docs/DocsPager";

export const Route = createFileRoute("/docs/use-cases")({
  head: () => ({
    meta: [
      { title: "Use Cases — Mr. Browser Docs" },
      {
        name: "description",
        content:
          "How teams use Mr. Browser for QA automation, RPA, web scraping, and CI/CD pipelines.",
      },
      { property: "og:title", content: "Use Cases — Mr. Browser Docs" },
      {
        property: "og:description",
        content: "QA automation, RPA, scraping — real use cases for Mr. Browser.",
      },
    ],
  }),
  component: UseCases,
});

const QA_EXAMPLE = `# conftest.py
import pytest
from mrbrowser import MrBrowser

@pytest.fixture(scope="module")
def browser():
    with MrBrowser(host="localhost", port=7331) as b:
        yield b

# test_auth.py
def test_login_flow(browser):
    page = browser.open("https://app.example.com/login")
    page.type("email field", "qa@example.com")
    page.type("password field", "Test@123")
    page.click("sign in button")
    page.wait(1)
    page.assert_text("Welcome back")
    page.assert_url_contains("/dashboard")

def test_login_invalid_password(browser):
    page = browser.open("https://app.example.com/login")
    page.type("email field", "qa@example.com")
    page.type("password field", "wrong")
    page.click("sign in button")
    page.assert_text("Invalid credentials")`;

const RPA_EXAMPLE = `name: download_weekly_invoices
description: "Log in to the billing portal and download all unpaid invoices"
steps:
  - open:
      url: "https://billing.corp.internal/login"
  - type:
      target: "Username"
      value: "$BILLING_USER"
  - type:
      target: "Password"
      value: "$BILLING_PASS"
      enter: true
  - wait:
      seconds: 2
  - click:
      target: "Invoices"
  - click:
      target: "Filter by Unpaid"
  - screenshot:
      output: "reports/invoices_snapshot.png"`;

const SCRAPING_EXAMPLE = `from mrbrowser import MrBrowser

with MrBrowser() as browser:
    page = browser.open("https://news.ycombinator.com")

    # Extract specific elements by plain-English description
    top_story = page.extract_text("the first story title")
    score = page.extract_text("the score of the first story")
    print(f"{top_story} — {score}")

    # Get all interactive elements for exploration
    links = [
        el for el in page.inspect(visible_only=True)
        if el.get("role") == "link"
    ]
    print(f"Found {len(links)} links on the page")`;

function UseCases() {
  return (
    <FadeIn>
      <article className="max-w-3xl">
        <p className="font-mono text-xs text-primary">Guides</p>
        <h1 className="mt-2 font-mono text-3xl font-bold tracking-tight sm:text-4xl">
          Use Cases
        </h1>
        <p className="mt-4 text-[15px] leading-relaxed text-muted-foreground">
          Mr. Browser is used across three main domains: QA/E2E testing, Robotic Process
          Automation (RPA), and intelligent web scraping. Each benefits from the same core
          guarantee — describe what you want in plain English and the engine figures out the
          rest.
        </p>

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> QA & E2E Testing
        </h2>
        <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
          QA teams lose days to failing tests caused by trivial UI changes — a class rename, a
          button re-labeled, a form restructured. Mr. Browser's{" "}
          <strong className="text-foreground">Memory Engine</strong> fingerprints every element
          after a successful run. When the UI changes, the engine finds where the element moved
          instead of failing. Tests only break when{" "}
          <em>actual business logic</em> breaks.
        </p>
        <CodeBlock className="mt-4" code={QA_EXAMPLE} lang="python" title="test_auth.py" />
        <Callout type="tip">
          Commit <code className="font-mono text-primary">./mrbrowser.db</code> to your repo
          (it's a compact SQLite file). Fresh CI runners then inherit the full fingerprint
          history from the last passing run — no cold-cache startup failures.
        </Callout>

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Robotic Process Automation (RPA)
        </h2>
        <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
          Automating legacy enterprise web apps has always meant brittle selectors and constant
          maintenance. Mr. Browser lets you write YAML workflows in plain English — accessible
          to non-developers and resilient to minor UI changes. No HTML knowledge required.
        </p>
        <CodeBlock
          className="mt-4"
          code={RPA_EXAMPLE}
          lang="yaml"
          title="download_invoices.yaml"
        />
        <p className="mt-4 text-sm leading-relaxed text-muted-foreground">
          Run this on a schedule with a cron job:{" "}
          <code className="font-mono text-primary">
            0 8 * * MON mr-browser run download_weekly_invoices.yaml
          </code>
        </p>
        <Callout type="warning">
          Never hardcode credentials in YAML files. Use{" "}
          <code className="font-mono text-destructive">$ENV_VAR</code> references — they are
          resolved at runtime and are never logged or stored in fingerprints.
        </Callout>

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Intelligent Web Scraping
        </h2>
        <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
          Modern websites obfuscate their HTML — React/Webpack apps emit class names like{" "}
          <code className="font-mono text-primary">.x9f2a</code> that change on every build.
          Mr. Browser matches against the browser's{" "}
          <strong className="text-foreground">Accessibility Tree</strong> instead, which
          remains stable across UI rebuilds.
        </p>
        <CodeBlock className="mt-4" code={SCRAPING_EXAMPLE} lang="python" title="scraper.py" />
        <Callout type="tip">
          Use <code className="font-mono text-primary">page.inspect(visible_only=True)</code>{" "}
          as a discovery tool when writing a new scraper. It shows every element the engine can
          see with its role, text, and confidence-ready metadata.
        </Callout>

        <DocsPager path="/docs/use-cases" />
      </article>
    </FadeIn>
  );
}
