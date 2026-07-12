import { createFileRoute } from "@tanstack/react-router";
import { FadeIn } from "@/components/FadeIn";
import { CodeBlock } from "@/components/CodeBlock";
import { Callout } from "@/components/docs/Callout";
import { DocsPager } from "@/components/docs/DocsPager";

export const Route = createFileRoute("/docs/core-concepts")({
  head: () => ({
    meta: [
      { title: "Core Concepts — Mr. Browser Docs" },
      {
        name: "description",
        content:
          "How intent resolution, the DOM Accessibility Tree, and the self-healing Memory Engine work inside Mr. Browser.",
      },
      { property: "og:title", content: "Core Concepts — Mr. Browser Docs" },
      {
        property: "og:description",
        content: "Intent resolution, accessibility trees, and the Memory Engine explained.",
      },
    ],
  }),
  component: CoreConcepts,
});

const INTENT_EXAMPLE = `- click: "the red delete button in the sidebar"

# The engine decomposes this into:
#   action:    click
#   role hint: button
#   text hint: "delete"
#   modifiers: color=red, region=sidebar`;

const FINGERPRINT = `{
  "intent": "Login",
  "resolved": { "role": "button", "name": "Login" },
  "fingerprint": {
    "ancestors": ["form[auth]", "main", "body"],
    "siblings":  ["input[Email]", "input[Password]"],
    "position":  { "region": "center", "order": 3 },
    "text_hash": "b2f9…"
  },
  "confidence": 0.99,
  "last_seen": "2026-07-01T09:14:22Z"
}`;

const HEALING = `$ mrbrowser run login.yaml

[intent] resolving "Login" → no direct match (0.41 < threshold)
[memory] consulting fingerprint b2f9…
[memory] structural match found: button moved
         form[auth] → dialog[auth-modal]   (0.94)
[memory] ✔ self-healed, fingerprint updated
✔ flow login_admin passed in 3.1s`;

function CoreConcepts() {
  return (
    <FadeIn>
      <article className="max-w-3xl">
        <p className="font-mono text-xs text-primary">Introduction</p>
        <h1 className="mt-2 font-mono text-3xl font-bold tracking-tight sm:text-4xl">
          Core Concepts
        </h1>
        <p className="mt-4 text-[15px] leading-relaxed text-muted-foreground">
          Three ideas make Mr. Browser different: intent resolution, the accessibility
          tree as the source of truth, and the self-healing Memory Engine.
        </p>

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> Intent-Driven Resolution
        </h2>
        <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
          An intent is a natural-language description of an element. A local NLP algorithm
          decomposes it into role, text, and contextual hints, then scores every node in
          the page against them. No cloud, no LLM required — resolution is deterministic
          and typically completes in under 30ms.
        </p>
        <CodeBlock className="mt-4" code={INTENT_EXAMPLE} lang="yaml" title="intent decomposition" />

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> The Accessibility Tree
        </h2>
        <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
          Instead of the raw DOM, the engine matches against the browser's{" "}
          <strong className="text-foreground">Accessibility Tree</strong> — the same
          semantic structure screen readers use. Roles, accessible names, and states are
          stable across redesigns even when class names are minified into{" "}
          <code className="font-mono text-primary">.x9f2a</code> hashes.
        </p>
        <Callout type="tip">
          Apps with good accessibility markup resolve faster and with higher confidence.
          Intent automation quietly rewards teams for writing accessible HTML.
        </Callout>

        <h2 className="mt-12 font-mono text-xl font-semibold text-foreground">
          <span className="text-primary">#</span> The Memory Engine
        </h2>
        <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
          Every successful resolution stores a structural fingerprint: ancestor chain,
          sibling context, region, and a text hash.
        </p>
        <CodeBlock className="mt-4" code={FINGERPRINT} lang="typescript" title=".mrbrowser/memory.db (entry)" />
        <p className="mt-4 text-sm leading-relaxed text-muted-foreground">
          When a later run fails to resolve an intent directly, the engine compares the
          live page against historical fingerprints to find where the element{" "}
          <em>moved</em> — then updates the fingerprint and continues:
        </p>
        <CodeBlock className="mt-4" code={HEALING} lang="bash" title="self-healing in action" />

        <Callout type="warning">
          Self-healing accepts a match only above the{" "}
          <code className="font-mono text-destructive">heal_threshold</code> (default
          0.85). Below it, the flow fails loudly instead of guessing — silent wrong-element
          clicks are worse than a red build.
        </Callout>

        <DocsPager path="/docs/core-concepts" />
      </article>
    </FadeIn>
  );
}
