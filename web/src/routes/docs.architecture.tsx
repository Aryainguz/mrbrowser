import { createFileRoute } from "@tanstack/react-router";
import { FadeIn } from "@/components/FadeIn";
import { CodeBlock } from "@/components/CodeBlock";
import { DocsPager } from "@/components/docs/DocsPager";

export const Route = createFileRoute("/docs/architecture")({
  head: () => ({
    meta: [
      { title: "Architecture — Mr. Browser Docs" },
      {
        name: "description",
        content:
          "How Mr. Browser works under the hood: Chrome DevTools Protocol, accessibility tree extraction, intent resolution, and the Memory Engine.",
      },
      { property: "og:title", content: "Architecture — Mr. Browser Docs" },
      {
        property: "og:description",
        content: "The layered architecture powering Mr. Browser's intent-driven automation.",
      },
    ],
  }),
  component: Architecture,
});

const ARCH_DIAGRAM = `┌─────────────────────────────────────────────────────────┐
│                      Client Layer                       │
│     Python SDK (mrbrowser)   TypeScript SDK (@mrbrowser/sdk) │
│     YAML Workflows (mr-browser run)                     │
└──────────────────────┬──────────────────────────────────┘
                       │  HTTP REST  (localhost:7331)
┌──────────────────────▼──────────────────────────────────┐
│                   CLI / Server (Go)                     │
│   cobra commands: run / debug / screenshot / inspect    │
└──────────────────────┬──────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────┐
│                   Runtime / Executor                    │
│   task.go: YAML parsing & validation                    │
│   executor.go: step dispatch, error handling            │
└──────────────────────┬──────────────────────────────────┘
                       │
          ┌────────────┴────────────┐
          │                        │
┌─────────▼──────────┐  ┌──────────▼──────────────────────┐
│  Actions Engine    │  │   Intelligence / Resolution      │
│  actions.go:       │  │   dom/extractor.go:              │
│  Click / Type /    │  │   Build Accessibility Tree       │
│  Hover / Scroll /  │  │                                 │
│  Upload / Screenshot│  │   intent/resolver.go:           │
│  verification.go:  │  │   Score elements against intent │
│  DOM-change checks │  │   (local NLP, no cloud)         │
└─────────┬──────────┘  └──────────┬──────────────────────┘
          │                        │
          └──────────┬─────────────┘
                     │
┌────────────────────▼──────────────────────────────────┐
│               Memory Engine (SQLite)                  │
│   memory.go: store & retrieve element fingerprints   │
│   self-heal: compare live DOM to historical snapshot  │
│   database: ./mrbrowser.db                           │
└────────────────────┬──────────────────────────────────┘
                     │
┌────────────────────▼──────────────────────────────────┐
│       Chrome DevTools Protocol (chromedp)             │
│   Navigate / Click / Type / Screenshot / JS execute   │
│   Accessibility tree snapshot (CDP Accessibility API) │
└────────────────────┬──────────────────────────────────┘
                     │
               [Chromium Browser]`;

const FINGERPRINT = `// memory/fingerprint (stored in mrbrowser.db)
{
  "intent":    "Login button",
  "resolved":  { "role": "button", "name": "Login" },
  "fingerprint": {
    "ancestors": ["form[auth]", "main", "body"],
    "siblings":  ["input[Email]", "input[Password]"],
    "position":  { "region": "center", "order": 3 },
    "text_hash": "b2f9a1…"
  },
  "confidence": 0.99,
  "last_seen":  "2026-07-01T09:14:22Z"
}`;

const RESOLUTION = `// Resolution pipeline (intelligence/dom + intent packages)
//
// 1. CDP Accessibility snapshot → []PageElement
// 2. For each element, compute score:
//      textScore     = fuzzy match(intent, element.text + aria-label)
//      roleScore     = match(role-hints in intent, element.role)
//      regionScore   = match(position-hints, element.boundingBox)
//      overallScore  = weighted average
// 3. Sort by score, take top candidate
// 4. If score < threshold → check Memory Engine fingerprints
// 5. If fingerprint match > heal_threshold → self-heal & update DB
// 6. If no match → fail step with descriptive error`;

function Architecture() {
  return (
    <FadeIn>
      <article className="max-w-3xl">
        <p className="font-mono text-xs text-primary">Reference</p>
        <h1 className="mt-2 font-mono text-3xl font-bold tracking-tight sm:text-4xl">
          Architecture
        </h1>
        <p className="mt-4 text-[15px] leading-relaxed text-muted-foreground">
          Mr. Browser is a layered Go application. The CLI and SDKs sit at the top; Chrome
          DevTools Protocol (CDP) via{" "}
          <code className="font-mono text-primary">chromedp</code> sits at the bottom. In
          between, the engine resolves plain-English intent against the browser's Accessibility
          Tree, executes actions, and stores structural fingerprints for future self-healing.
        </p>

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> System diagram
        </h2>
        <CodeBlock className="mt-4" code={ARCH_DIAGRAM} lang="bash" title="architecture" />

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Key packages
        </h2>
        <div className="mt-4 overflow-x-auto rounded-lg border border-border">
          <table className="w-full font-mono text-sm">
            <thead>
              <tr className="border-b border-border bg-muted/30">
                <th className="px-4 py-3 text-left font-semibold text-foreground">Package</th>
                <th className="px-4 py-3 text-left font-semibold text-foreground">Responsibility</th>
              </tr>
            </thead>
            <tbody className="text-muted-foreground">
              {[
                ["cli/cmd", "Cobra subcommands: run, debug, screenshot, inspect"],
                ["core/runtime", "YAML task parsing (task.go) and step execution (executor.go)"],
                ["core/actions", "Browser action primitives: Click, Type, Hover, Scroll, Upload, Screenshot"],
                ["core/browser", "chromedp session wrapper: Navigate, TypeSelector, ClickSelector, ExecuteJS"],
                ["intelligence/dom", "CDP Accessibility Tree extraction → []PageElement"],
                ["intelligence/intent", "Scoring engine: matches intent strings to PageElements"],
                ["memory", "SQLite fingerprint store — read/write/heal operations"],
                ["telemetry", "Structured logging (slog-based) with step/success/error helpers"],
              ].map(([pkg, desc]) => (
                <tr key={pkg} className="border-b border-border/50 hover:bg-muted/20">
                  <td className="px-4 py-3 text-primary">{pkg}</td>
                  <td className="px-4 py-3">{desc}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Intent resolution pipeline
        </h2>
        <CodeBlock className="mt-4" code={RESOLUTION} lang="typescript" title="resolution algorithm" />

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Memory Engine fingerprint
        </h2>
        <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
          Every successful resolution is persisted to a SQLite database (
          <code className="font-mono text-primary">./mrbrowser.db</code> by default). The
          fingerprint captures the structural context of the element — not its class name or
          ID, but where it lives in the DOM tree and what surrounds it.
        </p>
        <CodeBlock className="mt-4" code={FINGERPRINT} lang="typescript" title="fingerprint schema" />
        <p className="mt-4 text-sm leading-relaxed text-muted-foreground">
          When the intent resolver fails on a future run, the Memory Engine compares the
          current DOM against all stored fingerprints for that intent and finds the element
          that most structurally resembles the historical record — even if it has moved,
          been re-labeled, or wrapped in additional containers.
        </p>

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Technology choices
        </h2>
        <div className="mt-4 space-y-3 text-sm leading-relaxed text-muted-foreground">
          <p>
            <strong className="text-foreground">Go</strong> — chosen for performance, a
            single binary distribution, and excellent concurrency primitives for managing
            Chromium sessions.
          </p>
          <p>
            <strong className="text-foreground">chromedp</strong> — a pure-Go CDP client that
            gives direct access to the Accessibility Tree snapshot API, which other automation
            frameworks don't expose at the same depth.
          </p>
          <p>
            <strong className="text-foreground">SQLite (go-sqlite3)</strong> — zero-dependency
            embedded database. The fingerprint store is a single file, easily committed to
            version control alongside test code.
          </p>
          <p>
            <strong className="text-foreground">Local NLP only</strong> — no cloud calls, no
            API keys, no latency from LLM round-trips. Resolution is deterministic and
            completes in under 30 ms on a modern laptop.
          </p>
        </div>

        <DocsPager path="/docs/architecture" />
      </article>
    </FadeIn>
  );
}
