import { createFileRoute } from "@tanstack/react-router";
import { FadeIn } from "@/components/FadeIn";
import { CodeBlock } from "@/components/CodeBlock";
import { Callout } from "@/components/docs/Callout";
import { DocsPager } from "@/components/docs/DocsPager";

export const Route = createFileRoute("/docs/yaml-workflows")({
  head: () => ({
    meta: [
      { title: "YAML Workflow Reference — Mr. Browser Docs" },
      {
        name: "description",
        content:
          "Complete reference for Mr. Browser YAML workflow files — all actions, options, and environment variable support.",
      },
      { property: "og:title", content: "YAML Workflow Reference — Mr. Browser Docs" },
      {
        property: "og:description",
        content: "Complete YAML schema for Mr. Browser automation workflows.",
      },
    ],
  }),
  component: YamlWorkflows,
});

const FULL_EXAMPLE = `name: full-demo
description: "Demonstrates every available step type"
config:
  timeout_seconds: 30   # per-step timeout (default: 30)
  stop_on_error: true   # halt on first failure (default: true)

steps:
  # Navigate to a URL
  - open:
      url: "https://example.com/login"

  # Type text into an element
  - type:
      target: "Email"          # plain-English element description
      value: "$EMAIL"          # $ENV_VAR references are resolved at runtime
      clear: true              # clear field first (default: true)
      enter: false             # press Enter after typing (default: false)

  - type:
      target: "Password"
      value: "$PASSWORD"
      enter: true              # submit the form

  # Wait for navigation or content
  - wait:
      seconds: 2               # fixed pause

  - wait:
      url: "/dashboard"        # wait until URL contains this substring

  - wait:
      selector: "#main-content" # wait until CSS selector is visible

  # Click an element
  - click:
      target: "Download Report"

  # Hover over an element
  - hover:
      target: "the user avatar menu"

  # Scroll the page
  - scroll:
      direction: down          # up / down / left / right / top / bottom
      pixels: 800              # ignored for top/bottom

  # Take a screenshot
  - screenshot:
      output: "reports/dashboard.png"   # defaults to screenshot_<timestamp>.png

  # Extract text and save to session memory
  - extract:
      target: "the account balance"
      save_as: "balance"       # available as $balance in later steps

  # Assert conditions
  - assert:
      text_visible: "Welcome"       # text must appear on page

  - assert:
      url_contains: "/dashboard"    # URL must contain substring

  - assert:
      element_exists: "Sign Out"    # element matching intent must exist

  # Upload a file
  - upload:
      target: "the file upload input"
      file: "/tmp/report.csv"

  # Reload the page
  - reload: {}`;

const ENV_EXAMPLE = `# .env (never commit this)
EMAIL=admin@corp.com
PASSWORD=supersecret

# Run with env vars from shell
EMAIL=admin@corp.com PASSWORD=supersecret mr-browser run login.yaml

# Or export and run
export MRBROWSER_LOG_LEVEL=debug
mr-browser run login.yaml`;

const CLICK_SHORTHAND = `steps:
  # Long form (with optional selector override)
  - click:
      target: "Submit button"
      selector: "#submit-btn"   # optional: bypass resolver

  # Type long form
  - type:
      target: "Search box"
      value: "hello world"`;

function YamlWorkflows() {
  return (
    <FadeIn>
      <article className="max-w-3xl">
        <p className="font-mono text-xs text-primary">Guides</p>
        <h1 className="mt-2 font-mono text-3xl font-bold tracking-tight sm:text-4xl">
          YAML Workflow Reference
        </h1>
        <p className="mt-4 text-[15px] leading-relaxed text-muted-foreground">
          A workflow is a plain YAML file that describes a sequence of browser actions. Run it
          with{" "}
          <code className="font-mono text-primary">mr-browser run &lt;file.yaml&gt;</code>.
          The engine resolves every <code className="font-mono text-primary">target</code>{" "}
          string by searching the page's Accessibility Tree for the best match.
        </p>

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Complete schema
        </h2>
        <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
          Every available step type with all options documented inline:
        </p>
        <CodeBlock className="mt-4" code={FULL_EXAMPLE} lang="yaml" title="full-demo.yaml" />

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Environment variables
        </h2>
        <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
          Any <code className="font-mono text-primary">value</code> field supports{" "}
          <code className="font-mono text-primary">$VAR_NAME</code> references. They are
          resolved from the shell environment at runtime and are never logged or stored.
        </p>
        <CodeBlock className="mt-4" code={ENV_EXAMPLE} lang="bash" title="terminal" />
        <Callout type="warning">
          Using <code className="font-mono text-destructive">$ENV_VAR</code> syntax is the
          only safe way to pass credentials. Never hardcode passwords in YAML files — they
          will appear in git history.
        </Callout>

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Selector fallback
        </h2>
        <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
          Both <code className="font-mono text-primary">click</code> and{" "}
          <code className="font-mono text-primary">type</code> accept an optional{" "}
          <code className="font-mono text-primary">selector</code> field. When provided, the
          engine skips intent resolution and uses the CSS selector directly — useful for
          machine-generated IDs that are stable.
        </p>
        <CodeBlock className="mt-4" code={CLICK_SHORTHAND} lang="yaml" title="selector-override.yaml" />

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> All step types
        </h2>
        <div className="mt-4 overflow-x-auto rounded-lg border border-border">
          <table className="w-full font-mono text-sm">
            <thead>
              <tr className="border-b border-border bg-muted/30">
                <th className="px-4 py-3 text-left font-semibold text-foreground">Step</th>
                <th className="px-4 py-3 text-left font-semibold text-foreground">Required fields</th>
                <th className="px-4 py-3 text-left font-semibold text-foreground">Optional fields</th>
              </tr>
            </thead>
            <tbody className="text-muted-foreground">
              {[
                ["open", "url", "—"],
                ["click", "target", "selector"],
                ["type", "target, value", "clear, enter, selector"],
                ["hover", "target", "selector"],
                ["scroll", "direction", "pixels"],
                ["wait", "seconds | selector | url", "—"],
                ["screenshot", "—", "output"],
                ["extract", "target, save_as", "—"],
                ["assert", "text_visible | url_contains | element_exists", "—"],
                ["upload", "target, file", "selector"],
                ["reload", "—", "—"],
              ].map(([step, req, opt]) => (
                <tr key={step} className="border-b border-border/50 hover:bg-muted/20">
                  <td className="px-4 py-3 text-primary">{step}:</td>
                  <td className="px-4 py-3">{req}</td>
                  <td className="px-4 py-3 text-muted-foreground/60">{opt}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        <DocsPager path="/docs/yaml-workflows" />
      </article>
    </FadeIn>
  );
}
