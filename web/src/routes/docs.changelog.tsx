import { createFileRoute } from "@tanstack/react-router";
import { FadeIn } from "@/components/FadeIn";
import { Callout } from "@/components/docs/Callout";
import { DocsPager } from "@/components/docs/DocsPager";

export const Route = createFileRoute("/docs/changelog")({
  head: () => ({
    meta: [
      { title: "Changelog — Mr. Browser Docs" },
      {
        name: "description",
        content: "Mr. Browser release history — what changed in each version.",
      },
      { property: "og:title", content: "Changelog — Mr. Browser Docs" },
      {
        property: "og:description",
        content: "Mr. Browser release history.",
      },
    ],
  }),
  component: Changelog,
});

interface Release {
  version: string;
  date: string;
  type: "major" | "minor" | "patch";
  changes: { kind: "feat" | "fix" | "perf" | "break"; text: string }[];
}

const RELEASES: Release[] = [
  {
    version: "0.2.0",
    date: "2026-07-01",
    type: "minor",
    changes: [
      { kind: "feat", text: "New debug command — step through workflows interactively with ENTER/skip/quit prompts." },
      { kind: "feat", text: "hover: step type — trigger mouseover/mouseenter events on any element." },
      { kind: "feat", text: "upload: step type — set file inputs for upload automation." },
      { kind: "feat", text: "extract: step type — save element text to session memory for use in later steps." },
      { kind: "feat", text: "reload: step type — reload the current page." },
      { kind: "feat", text: "wait.url — wait until the page URL contains a given substring." },
      { kind: "feat", text: "assert.element_exists — assert an element matching an intent exists on page." },
      { kind: "feat", text: "inspect --target flag — test a specific intent string and see ranked candidates." },
      { kind: "feat", text: "Python SDK: hover(), upload(), extract_text(), execute_js(), get_cookies(), set_cookies(), wait_for_url(), assert_url_contains() methods added." },
      { kind: "feat", text: "TypeScript SDK: matching methods added." },
      { kind: "perf", text: "DOM Accessibility snapshot caching — 40% faster on pages with repeated intent lookups." },
      { kind: "fix", text: "Self-heal threshold now correctly defaults to 0.85 (was 0.7 in v0.1.0)." },
      { kind: "fix", text: "Password values are now masked as *** in all logs regardless of field name." },
    ],
  },
  {
    version: "0.1.0",
    date: "2026-06-01",
    type: "minor",
    changes: [
      { kind: "feat", text: "Initial public release." },
      { kind: "feat", text: "CLI: run, screenshot, inspect commands." },
      { kind: "feat", text: "YAML steps: open, click, type, scroll, screenshot, wait, assert (text_visible, url_contains)." },
      { kind: "feat", text: "Python SDK with MrBrowser client and Page context manager." },
      { kind: "feat", text: "TypeScript SDK with MrBrowser class and Page." },
      { kind: "feat", text: "Memory Engine: SQLite fingerprint store with structural self-healing." },
      { kind: "feat", text: "Global flags: --headless, --no-sandbox, --db, --log-level, --chromium, --config." },
      { kind: "feat", text: "Docker Compose setup with Chromium, volume mounts, and health checks." },
      { kind: "feat", text: "DOM Accessibility Tree extraction via Chrome DevTools Protocol." },
      { kind: "feat", text: "Local NLP intent resolution — no cloud, no LLM latency." },
    ],
  },
];

const kindLabel: Record<string, { label: string; cls: string }> = {
  feat:  { label: "feat",  cls: "text-primary border-primary/40 bg-primary/10" },
  fix:   { label: "fix",   cls: "text-yellow-400 border-yellow-400/40 bg-yellow-400/10" },
  perf:  { label: "perf",  cls: "text-blue-400 border-blue-400/40 bg-blue-400/10" },
  break: { label: "break", cls: "text-destructive border-destructive/40 bg-destructive/10" },
};

function Changelog() {
  return (
    <FadeIn>
      <article className="max-w-3xl">
        <p className="font-mono text-xs text-primary">Project</p>
        <h1 className="mt-2 font-mono text-3xl font-bold tracking-tight sm:text-4xl">
          Changelog
        </h1>
        <p className="mt-4 text-[15px] leading-relaxed text-muted-foreground">
          All notable changes are documented here. Mr. Browser follows{" "}
          <a
            href="https://semver.org"
            target="_blank"
            rel="noreferrer"
            className="text-primary hover:underline"
          >
            Semantic Versioning
          </a>
          .
        </p>

        <Callout type="tip">
          To upgrade, pull the latest tag and rebuild:{" "}
          <code className="font-mono text-primary">git pull && make build && make install</code>.
          Check the breaking changes section before upgrading across minor versions.
        </Callout>

        <div className="mt-12 space-y-14">
          {RELEASES.map((release) => (
            <section key={release.version}>
              <div className="flex items-baseline gap-3">
                <h2 className="font-mono text-2xl font-bold text-foreground">
                  <span className="text-primary">v</span>
                  {release.version}
                </h2>
                <span className="font-mono text-xs text-muted-foreground">{release.date}</span>
                <span
                  className={`rounded-full border px-2 py-0.5 font-mono text-[10px] font-semibold ${
                    release.type === "major"
                      ? "border-destructive/40 bg-destructive/10 text-destructive"
                      : release.type === "minor"
                      ? "border-primary/40 bg-primary/10 text-primary"
                      : "border-border text-muted-foreground"
                  }`}
                >
                  {release.type}
                </span>
              </div>
              <ul className="mt-6 space-y-3">
                {release.changes.map((change, i) => {
                  const badge = kindLabel[change.kind];
                  return (
                    <li key={i} className="flex items-start gap-3">
                      <span
                        className={`mt-0.5 shrink-0 rounded border px-1.5 py-0.5 font-mono text-[10px] font-semibold ${badge.cls}`}
                      >
                        {badge.label}
                      </span>
                      <span className="text-sm leading-relaxed text-muted-foreground">
                        {change.text}
                      </span>
                    </li>
                  );
                })}
              </ul>
            </section>
          ))}
        </div>

        <DocsPager path="/docs/changelog" />
      </article>
    </FadeIn>
  );
}
