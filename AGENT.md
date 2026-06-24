# Mr. Browser — Agent Guide

Welcome to the Mr. Browser repository! This document serves as a guide for AI agents contributing to or utilizing the Mr. Browser codebase.

## 🎯 What is Mr. Browser?

Mr. Browser is an **intent-driven, self-healing browser automation engine** built in Go, with native zero-dependency SDKs for Python and TypeScript.

Unlike Selenium or Playwright, users do not write CSS selectors. They write intents (e.g. `target: "Login Button"`). The engine injects native JS to extract an accessibility/interactive tree, scores targets using Jaro-Winkler NLP similarity, clicks the best match, and stores the DOM fingerprint in SQLite (`mrbrowser.db`). If the site changes and the NLP match fails later, the engine falls back to historical fingerprints to self-heal the workflow.

## 📸 Tool in Action

Here is an example of Mr. Browser executing a YAML workflow autonomously.

```yaml
steps:
  - open:
      url: https://en.wikipedia.org/wiki/Main_Page
  - type:
      target: "Search Wikipedia"
      value: "Browser automation"
      enter: true
  - wait:
      seconds: 1
  - screenshot:
      output: "search_results.png"
```

The engine parsed the target `"Search Wikipedia"`, automatically mapped it to the Wikipedia search input field (`#searchInput`), typed the value, pressed enter, and captured the following screenshot:

![search_results.png](/Users/aryainguzsupertramp/Developer/otel-java/mrbrowser/search_results.png)

*(The engine correctly inferred the input field without any DOM selectors being provided by the user!)*

## 🏗️ Codebase Layout

- **`core/runtime`**: Contains the YAML parser, Step definitions, and the `Executor` orchestration loop.
- **`core/browser`**: The CDP driver wrapping the `chromedp` library. Handles raw tab lifecycle (`Browser`, `Page`, `Element`).
- **`core/actions`**: Physical primitive actions (`Click`, `Type`, `Scroll`). Contains before/after DOM verification logic to ensure actions actually trigger state changes.
- **`intelligence/`**: The brain of the operation.
  - `dom/extractor.go`: Injects JS to map interactive DOM elements.
  - `accessibility/tree.go`: Fetches the Chrome AX tree.
  - `resolver/`: The NLP string-similarity engine.
- **`memory/`**: SQLite persistence (`mattn/go-sqlite3`). Tracks telemetry and DOM fingerprints for self-healing.
- **`server/`**: The Go REST API (`:7331`) that listens for SDK commands.
- **`cli/`**: The `mr-browser` cobra CLI tool for executing YAML files directly.
- **`sdk/`**: Python and TypeScript SDKs. Both are *zero external dependency* wrappers around the REST API.
- **`docs/`**: A VitePress documentation site.

## 🛠️ Contribution Guidelines for Agents

1. **Prioritize Native Go**: If adding features to the engine, rely on standard library or established tools (`chromedp`). Do not add heavy NLP/ML dependencies to the core binary.
2. **REST First**: Any new feature added to the core engine must be exposed via `server/handlers.go` and implemented identically in both the TypeScript and Python SDKs.
3. **No Selectors**: Do not add CSS selector parsing logic to the core intent resolver. The system is fundamentally semantic. Wait until Phase 5 (Vision) to add bounding box matching.
4. **Documentation**: If you change an API or SDK method, immediately update the VitePress markdown files in `docs/` and run `npm run docs:build` to verify formatting.
