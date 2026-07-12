import { createFileRoute } from "@tanstack/react-router";
import { FadeIn } from "@/components/FadeIn";
import { CodeBlock } from "@/components/CodeBlock";
import { Callout } from "@/components/docs/Callout";
import { DocsPager } from "@/components/docs/DocsPager";

export const Route = createFileRoute("/docs/contributing")({
  head: () => ({
    meta: [
      { title: "Contributing — Mr. Browser Docs" },
      {
        name: "description",
        content:
          "How to contribute to Mr. Browser: project setup, running tests, submitting PRs, and code style guidelines.",
      },
      { property: "og:title", content: "Contributing — Mr. Browser Docs" },
      {
        property: "og:description",
        content: "Contribution guide for Mr. Browser — setup, tests, PRs.",
      },
    ],
  }),
  component: Contributing,
});

const SETUP = `# Prerequisites: Go 1.22+, Chromium, make

git clone https://github.com/aryainguz/mrbrowser.git
cd mrbrowser

# Download Go dependencies
make deps

# Build the binary
make build
# → bin/mr-browser

# Run unit tests (no Chromium needed)
make test-unit

# Run integration tests (requires Chromium in PATH)
make test-int

# Run all tests
make test

# Format + vet
make fmt
make vet`;

const STRUCTURE = `mrbrowser/
├── cli/             # Cobra CLI commands (main entrypoint)
│   └── cmd/         # run.go, debug.go, screenshot.go, inspect.go, root.go
├── core/
│   ├── actions/     # Browser action primitives (Click, Type, Scroll…)
│   ├── browser/     # chromedp session wrapper
│   └── runtime/     # YAML task parsing (task.go) & execution (executor.go)
├── intelligence/
│   ├── dom/         # CDP Accessibility Tree extraction
│   └── intent/      # Intent resolution & scoring engine
├── memory/          # SQLite fingerprint store & self-healing logic
├── telemetry/       # Structured logging helpers
├── sdk/
│   ├── python/      # Python SDK (mrbrowser package)
│   └── typescript/  # TypeScript SDK (@mrbrowser/sdk)
├── docker/          # Dockerfile & docker-compose.yml
├── examples/        # Example YAML workflows
├── tests/
│   ├── unit/        # Fast unit tests
│   └── integration/ # Tests that require a live Chromium
└── docs/            # VitePress documentation site`;

const NEW_STEP = `// core/runtime/task.go
// 1. Add your step struct
type ConditionalStep struct {
    Target    string \`yaml:"target"\`
    Condition string \`yaml:"if"\`
}

// 2. Add it to the Step union
type Step struct {
    // ... existing fields ...
    Conditional *ConditionalStep \`yaml:"conditional,omitempty"\`
}

// 3. Add Kind() case
case s.Conditional != nil:
    return "conditional"

// core/runtime/executor.go
// 4. Add execution case
case step.Conditional != nil:
    return e.execConditional(step.Conditional)`;

function Contributing() {
  return (
    <FadeIn>
      <article className="max-w-3xl">
        <p className="font-mono text-xs text-primary">Project</p>
        <h1 className="mt-2 font-mono text-3xl font-bold tracking-tight sm:text-4xl">
          Contributing
        </h1>
        <p className="mt-4 text-[15px] leading-relaxed text-muted-foreground">
          Mr. Browser is open source under the MIT license. Contributions of all kinds are
          welcome — bug reports, documentation improvements, new step types, and performance
          work.
        </p>

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Local setup
        </h2>
        <CodeBlock className="mt-4" code={SETUP} lang="bash" title="setup" />
        <Callout type="tip">
          Unit tests run without Chromium and are fast (&lt;5 s). Integration tests require
          Chromium in your <code className="font-mono text-primary">PATH</code>. On macOS,{" "}
          <code className="font-mono text-primary">brew install chromium</code>; on Ubuntu,{" "}
          <code className="font-mono text-primary">apt install chromium-browser</code>.
        </Callout>

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Repository structure
        </h2>
        <CodeBlock className="mt-4" code={STRUCTURE} lang="bash" title="directory layout" />

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Adding a new step type
        </h2>
        <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
          New YAML step types follow a four-file pattern:
        </p>
        <CodeBlock className="mt-4" code={NEW_STEP} lang="typescript" title="adding a step" />

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Pull request checklist
        </h2>
        <ul className="mt-4 space-y-2 font-mono text-sm text-muted-foreground">
          {[
            "make fmt && make vet pass with no output",
            "make test-unit passes",
            "New code has corresponding unit tests in tests/unit/",
            "New step types have an example in examples/",
            "YAML schema changes are reflected in docs/sdk/yaml-workflows",
            "No new mandatory dependencies (stdlib only for SDKs)",
            "PR description explains the motivation and links to any related issue",
          ].map((item) => (
            <li key={item} className="flex items-start gap-2">
              <span className="mt-0.5 text-primary">›</span>
              <span>{item}</span>
            </li>
          ))}
        </ul>

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Code style
        </h2>
        <div className="mt-4 space-y-3 text-sm leading-relaxed text-muted-foreground">
          <p>
            <strong className="text-foreground">Go:</strong> Standard{" "}
            <code className="font-mono text-primary">gofmt</code> formatting. Exported
            functions and types must have doc comments. Error strings are lowercase, no
            trailing punctuation. Wrap errors with{" "}
            <code className="font-mono text-primary">fmt.Errorf("context: %w", err)</code>.
          </p>
          <p>
            <strong className="text-foreground">Python SDK:</strong> PEP 8. All public methods
            must have Google-style docstrings. No external dependencies — stdlib only.
          </p>
          <p>
            <strong className="text-foreground">TypeScript SDK:</strong> ESM modules, strict
            TypeScript. No runtime dependencies — native fetch only (Node 18+).
          </p>
        </div>
        <Callout type="warning">
          The <code className="font-mono text-destructive">--no-sandbox</code> flag must only
          be set by the Docker container, never hardcoded in tests or default config. It
          disables a Chromium security boundary and should never be the default for developer
          machines.
        </Callout>

        <DocsPager path="/docs/contributing" />
      </article>
    </FadeIn>
  );
}
