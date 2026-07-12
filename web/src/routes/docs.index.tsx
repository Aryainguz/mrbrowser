import { createFileRoute } from "@tanstack/react-router";
import { FadeIn } from "@/components/FadeIn";
import { CodeBlock } from "@/components/CodeBlock";
import { Callout } from "@/components/docs/Callout";
import { DocsPager } from "@/components/docs/DocsPager";

export const Route = createFileRoute("/docs/")({
  head: () => ({
    meta: [
      { title: "Getting Started — Mr. Browser Docs" },
      {
        name: "description",
        content:
          "Install Mr. Browser and run your first plain-English automation flow in under five minutes.",
      },
      { property: "og:title", content: "Getting Started — Mr. Browser Docs" },
      {
        property: "og:description",
        content: "Install Mr. Browser and run your first flow in five minutes.",
      },
    ],
  }),
  component: GettingStarted,
});

const INSTALL = `# Docker (recommended — includes engine + Chromium)
docker compose -f docker/docker-compose.yml up -d

# Build from source (requires Go 1.22+)
git clone https://github.com/mrbrowser/mrbrowser.git
cd mrbrowser
make build           # produces bin/mr-browser
make install         # copies to /usr/local/bin/mr-browser

# Python SDK
pip install mrbrowser

# TypeScript SDK
npm install @mrbrowser/sdk`;

const FIRST_FLOW = `# login.yaml
name: login_admin
steps:
  - open:
      url: "https://corp-portal.internal/login"
  - type:
      target: "Email"
      value: "admin@corp.com"
  - type:
      target: "Password"
      value: "$SECRET_PASS"
      enter: true
  - wait:
      seconds: 1
  - assert:
      text_visible: "Dashboard"`;

const RUN = `$ mr-browser run login.yaml

[intent] resolving "Email"      → <input aria-label="Email">     (0.98)
[intent] resolving "Password"   → <input type="password">        (0.97)
[memory] 2 fingerprints stored  → ./mrbrowser.db
✔ login_admin — 2/2 steps passed in 2.4s

# Run with headed browser (watch it happen)
$ mr-browser run login.yaml --headless=false

# Step through interactively
$ mr-browser debug login.yaml --headless=false`;

function GettingStarted() {
  return (
    <FadeIn>
      <article className="max-w-3xl">
        <p className="font-mono text-xs text-primary">Introduction</p>
        <h1 className="mt-2 font-mono text-3xl font-bold tracking-tight sm:text-4xl">
          Getting Started
        </h1>
        <p className="mt-4 text-[15px] leading-relaxed text-muted-foreground">
          Mr. Browser is an open-source browser automation engine that resolves elements
          from <strong className="text-foreground">plain-English intent</strong> instead of
          CSS selectors. This guide gets you from zero to a passing flow in about five
          minutes.
        </p>

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Installation
        </h2>
        <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
          The engine ships as a Docker image with a bundled dashboard, or as a lightweight
          SDK for Python and TypeScript.
        </p>
        <CodeBlock className="mt-4" code={INSTALL} lang="bash" title="terminal" />

        <Callout type="tip">
          The Docker image bundles a Chromium installation and writes data to
          a volume at <code className="font-mono text-primary">/app/data</code>.
          The engine listens on{" "}
          <code className="font-mono text-primary">http://localhost:7331</code>.
          The SDKs talk to it over that port — nothing ever leaves your machine.
        </Callout>

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Your first flow
        </h2>
        <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
          Flows are plain YAML. Each step declares an <em>intent</em> — a human description
          of the element — and the engine resolves it against the DOM Accessibility Tree.
        </p>
        <CodeBlock className="mt-4" code={FIRST_FLOW} lang="yaml" title="login.yaml" />

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Run it
        </h2>
        <CodeBlock className="mt-4" code={RUN} lang="bash" title="terminal" />
        <p className="mt-4 text-sm leading-relaxed text-muted-foreground">
          Each resolution logs a confidence score and stores a structural fingerprint in
          the Memory Engine. If the portal's UI changes next quarter, those fingerprints
          are what let your flow heal itself.
        </p>

        <Callout type="warning">
          Never hardcode credentials in flow files. Use{" "}
          <code className="font-mono text-destructive">$ENV_VAR</code> references — they are
          resolved at runtime and redacted from all logs and fingerprints.
        </Callout>

        <DocsPager path="/docs" />
      </article>
    </FadeIn>
  );
}
