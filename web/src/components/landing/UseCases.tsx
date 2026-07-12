import { FlaskConical, Bot, Braces } from "lucide-react";
import { CodeBlock } from "@/components/CodeBlock";
import { FadeIn } from "@/components/FadeIn";

const PYTEST_SNIPPET = `import mrbrowser

def test_checkout_survives_redesign(browser):
    browser.goto("https://shop.example.com")
    browser.click("the first product card")
    browser.click("Add to cart")
    browser.click("the checkout button")
    assert browser.sees("Order summary")`;

const RPA_SNIPPET = `flow: invoice_export
schedule: "0 6 * * MON"
steps:
  - goto: "https://legacy-erp.corp/invoices"
  - select: { target: "Status filter", value: "Unpaid" }
  - click: "Export to CSV"
  - wait_for: "Download complete"`;

const SCRAPE_SNIPPET = `const rows = await browser.extract({
  intent: "every product row in the results table",
  fields: {
    name: "the product title",
    price: "the listed price",
    stock: "availability status",
  },
}); // works on minified class names like .x9f2a`;

const CASES = [
  {
    icon: FlaskConical,
    tag: "qa_e2e",
    title: "QA & E2E Testing",
    body: "Write tests that survive UI redesigns. Stop fixing broken selectors. When the frontend team ships a rewrite, the Memory Engine re-resolves every element by structural fingerprint — your suite stays green.",
    lang: "python",
    file: "test_checkout.py",
    code: PYTEST_SNIPPET,
  },
  {
    icon: Bot,
    tag: "rpa",
    title: "Robotic Process Automation",
    body: "Automate legacy enterprise web apps with simple YAML workflows. No coding required. Schedule flows, chain steps, and let intent resolution deal with the 2009-era markup.",
    lang: "yaml",
    file: "invoice_export.yaml",
    code: RPA_SNIPPET,
  },
  {
    icon: Braces,
    tag: "scraping",
    title: "Intelligent Web Scraping",
    body: "Extract data semantically from heavily obfuscated React/Webpack applications. Hashed class names and shifting DOM structures are irrelevant when you target meaning, not markup.",
    lang: "typescript",
    file: "scrape.ts",
    code: SCRAPE_SNIPPET,
  },
];

export function UseCases() {
  return (
    <section
      id="use-cases"
      className="scanlines scroll-mt-20 border-y border-border bg-[oklch(0.09_0_0)] py-24"
    >
      <div className="mx-auto max-w-6xl px-4 sm:px-6">
        <FadeIn>
          <p className="font-mono text-xs text-primary">$ mrbrowser --use-cases</p>
          <h2 className="mt-3 font-mono text-3xl font-bold tracking-tight sm:text-4xl">
            One engine. <span className="text-primary text-glow">Three weapons.</span>
          </h2>
        </FadeIn>

        <div className="mt-16 space-y-20">
          {CASES.map((c, i) => (
            <FadeIn key={c.tag}>
              <div
                className={`grid items-center gap-8 lg:grid-cols-2 ${
                  i % 2 === 1 ? "lg:[&>*:first-child]:order-2" : ""
                }`}
              >
                <div>
                  <div className="flex items-center gap-3">
                    <div className="flex h-10 w-10 items-center justify-center rounded-lg border border-primary/30 bg-primary/10">
                      <c.icon size={20} className="text-primary" />
                    </div>
                    <span className="font-mono text-xs text-muted-foreground">
                      // {c.tag}
                    </span>
                  </div>
                  <h3 className="mt-4 font-mono text-xl font-semibold sm:text-2xl">
                    {c.title}
                  </h3>
                  <p className="mt-3 max-w-md text-sm leading-relaxed text-muted-foreground">
                    {c.body}
                  </p>
                </div>
                <CodeBlock code={c.code} lang={c.lang} title={c.file} />
              </div>
            </FadeIn>
          ))}
        </div>
      </div>
    </section>
  );
}
