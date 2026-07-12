import { createFileRoute } from "@tanstack/react-router";
import { FadeIn } from "@/components/FadeIn";
import { CodeBlock } from "@/components/CodeBlock";
import { Callout } from "@/components/docs/Callout";
import { DocsPager } from "@/components/docs/DocsPager";

export const Route = createFileRoute("/docs/cli")({
  head: () => ({
    meta: [
      { title: "CLI Reference — Mr. Browser Docs" },
      {
        name: "description",
        content:
          "Complete reference for the mr-browser CLI: run, debug, screenshot, inspect commands and all global flags.",
      },
      { property: "og:title", content: "CLI Reference — Mr. Browser Docs" },
      {
        property: "og:description",
        content: "mr-browser CLI commands: run, debug, screenshot, inspect — full reference.",
      },
    ],
  }),
  component: CliReference,
});

const GLOBAL_FLAGS = `mr-browser [flags] <command>

Global flags (available on every command):
  --headless          Run browser headlessly (default: true)
  --headless=false    Show the browser window
  --chromium <path>   Path to Chromium executable (auto-detected by default)
  --no-sandbox        Disable Chromium sandbox (required inside Docker)
  --db <path>         Path to SQLite fingerprint database (default: ./mrbrowser.db)
  --log-level <lvl>   Log verbosity: debug, info, warn, error (default: info)
  --config <file>     Config file path (default: $HOME/.mrbrowser.yaml or ./.mrbrowser.yaml)

Environment variables (MRBROWSER_ prefix):
  MRBROWSER_HEADLESS=true
  MRBROWSER_NO_SANDBOX=true
  MRBROWSER_CHROMIUM_PATH=/usr/bin/chromium
  MRBROWSER_DB_PATH=/app/data/mrbrowser.db
  MRBROWSER_LOG_LEVEL=info`;

const RUN_CMD = `# Run a YAML workflow headlessly (default)
mr-browser run login.yaml

# Show the browser window
mr-browser run login.yaml --headless=false

# Use a custom database path
mr-browser run login.yaml --db ./ci/memory.db

# Verbose logging
mr-browser run login.yaml --log-level debug

# Output:
#   Task: login_admin
#   Steps: 3/3 passed
#   Duration: 2.4s
#   Status: SUCCESS`;

const DEBUG_CMD = `# Step through a workflow interactively
mr-browser debug login.yaml --headless=false

# At each step you are prompted:
#   ── Step 1/3 [open] ─────────────────
#     URL: https://corp-portal.internal/login
#   Press ENTER to execute, 's' to skip, 'q' to quit: 

# Keyboard shortcuts during debug:
#   ENTER  — execute the current step
#   s      — skip the current step
#   q      — quit the debug session`;

const SCREENSHOT_CMD = `# Capture a screenshot of a URL
mr-browser screenshot https://example.com

# Save to a specific file
mr-browser screenshot https://example.com --output page.png
mr-browser screenshot https://example.com -o reports/snapshot.png

# With headed browser (useful for pages that need interaction first)
mr-browser screenshot https://app.example.com --headless=false`;

const INSPECT_CMD = `# List all interactive elements on a page
mr-browser inspect https://example.com

# Show only visible elements
mr-browser inspect https://example.com --visible-only

# Resolve a specific intent and show ranked candidates
mr-browser inspect https://example.com --target "login button"

# Output as JSON (useful for piping to jq)
mr-browser inspect https://example.com --json | jq '.[0]'

# Example output:
#   🌐 https://example.com
#   📄 Example Domain
#   Found 12 elements
#
#   ── INPUTS ─────────────────────
#   👁  <email>      text      Email address [clickable]
#   👁  <password>   password  Password [clickable]
#   ── BUTTONS ────────────────────
#   👁  <button>     button    Sign In [clickable]
#      <button>     button    Forgot password?`;

function CliReference() {
  return (
    <FadeIn>
      <article className="max-w-3xl">
        <p className="font-mono text-xs text-primary">Guides</p>
        <h1 className="mt-2 font-mono text-3xl font-bold tracking-tight sm:text-4xl">
          CLI Reference
        </h1>
        <p className="mt-4 text-[15px] leading-relaxed text-muted-foreground">
          The <code className="font-mono text-primary">mr-browser</code> binary provides four
          subcommands. Install it with{" "}
          <code className="font-mono text-primary">make install</code> or build with{" "}
          <code className="font-mono text-primary">make build</code> (produces{" "}
          <code className="font-mono text-primary">bin/mr-browser</code>).
        </p>

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Global flags
        </h2>
        <CodeBlock className="mt-4" code={GLOBAL_FLAGS} lang="bash" title="mr-browser --help" />

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> run
        </h2>
        <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
          Execute a YAML workflow from start to finish. Exits with code{" "}
          <code className="font-mono text-primary">0</code> on success,{" "}
          <code className="font-mono text-destructive">1</code> on failure.
        </p>
        <CodeBlock className="mt-4" code={RUN_CMD} lang="bash" title="mr-browser run" />

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> debug
        </h2>
        <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
          Step through a workflow interactively — pause after each step, inspect the result,
          and decide to continue, skip, or abort. Perfect for writing new flows.
        </p>
        <CodeBlock className="mt-4" code={DEBUG_CMD} lang="bash" title="mr-browser debug" />
        <Callout type="tip">
          Always use <code className="font-mono text-primary">--headless=false</code> with{" "}
          <code className="font-mono text-primary">debug</code> so you can visually watch
          what the engine is doing at each step.
        </Callout>

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> screenshot
        </h2>
        <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
          Navigate to a URL and capture a full-page PNG screenshot. Useful as a quick sanity
          check or for visual regression baselines.
        </p>
        <CodeBlock
          className="mt-4"
          code={SCREENSHOT_CMD}
          lang="bash"
          title="mr-browser screenshot"
        />

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> inspect
        </h2>
        <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
          Navigate to a URL and print all elements the engine can see. Use this to explore a
          new page and craft accurate intent strings for your workflows.{" "}
          <code className="font-mono text-primary">--target</code> lets you test a specific
          intent string and see which element it resolves to and with what confidence.
        </p>
        <CodeBlock className="mt-4" code={INSPECT_CMD} lang="bash" title="mr-browser inspect" />
        <Callout type="warning">
          <code className="font-mono text-destructive">inspect</code> adds a 500 ms settle
          delay after navigation to let JS-heavy SPAs finish rendering. For faster pages you
          can increase <code className="font-mono text-destructive">--log-level debug</code>{" "}
          to see timing details.
        </Callout>

        <DocsPager path="/docs/cli" />
      </article>
    </FadeIn>
  );
}
