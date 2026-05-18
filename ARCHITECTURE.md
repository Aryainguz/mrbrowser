# Architecture Deep Dive

Mr. Browser represents a paradigm shift in browser automation. Rather than forcing developers to write fragile, hard-coded scripts tied to specific DOM structures (CSS selectors or XPath), Mr. Browser relies on **intent-driven automation** backed by a **self-healing engine**.

This document details the architectural decisions, design patterns, and tradeoffs made in building the engine, as well as future considerations.

---

## 1. Core Principles & Philosophy

### Intent Over Structure
**Decision:** Users specify *what* they want to interact with (e.g., `target: "Login Button"`), not *how* to find it.
**Rationale:** Web interfaces are dynamic. A/B testing, redesigns, and dynamic class generation (e.g., Tailwind, CSS modules) frequently break traditional automation. By resolving targets based on semantic meaning and accessibility roles, the engine decouples the script from the UI implementation.

### Zero Mandatory AI Dependency
**Decision:** The core intent resolution is built in native Go using NLP string similarity (Jaro-Winkler) and accessibility tree parsing, rather than relying on expensive, high-latency LLM calls.
**Rationale:** Automation must be fast and deterministic. LLMs introduce latency (often 2-5 seconds per action) and non-deterministic behavior. LLMs are reserved as optional extensions (Phase 5/6) for complex reasoning tasks, not basic element discovery.

### Self-Healing via Memory
**Decision:** Every successful interaction fingerprints the target element (recording its outer HTML, bounding box, tree depth, and relative siblings) and stores it in a local SQLite database (`mrbrowser.db`).
**Rationale:** If the NLP resolver fails because a button's text changed drastically, the engine falls back to historical fingerprints. It scores current DOM elements against past fingerprints to "heal" the workflow and find where the element moved.

---

## 2. Component Architecture

The system is decoupled into five distinct layers to enforce separation of concerns.

### 2.1 The Driver Layer (`core/browser`)
**Role:** Manages the Chromium lifecycle via the Chrome DevTools Protocol (CDP).
**Implementation:** Built on top of the `chromedp` Go library.
**Decisions:** 
- We chose CDP over WebDriver (Selenium) because CDP provides deeper access to the browser's accessibility tree, raw DOM, and network stack without requiring an intermediary driver binary (like `chromedriver`).
- **Context Management:** The `Browser` struct manages the root allocation context, while each `Page` gets an isolated sub-context to prevent memory leaks during long-running scraping tasks.

### 2.2 The Intelligence Layer (`intelligence/`)
**Role:** Extracts the DOM and resolves intent.
**Implementation:** 
- **DOM Extractor:** Injects a lightweight vanilla JavaScript payload into the page context. This script walks the DOM and returns a sanitized JSON tree of interactive elements.
- **Intent Resolver:** Parses the user's intent (e.g., "Email Address") to infer the likely HTML role (e.g., `<input type="email">`). It then scores all extracted elements using weighted criteria: exact text match (highest weight), partial match, accessibility aria-labels, and placeholders.

### 2.3 The Action Layer (`core/actions`)
**Role:** Executes physical actions and verifies state mutations.
**Implementation:** 
- **Verification:** Before an action (like `Click`) is executed, the engine snapshots the URL and raw HTML. After the click, it waits 300ms and compares the new state against the snapshot.
- **Why?** Many single-page applications (SPAs) intercept clicks to trigger network requests without changing the URL. Verification ensures the engine knows if an action actually triggered a UI mutation.

### 2.4 The Execution Layer (`core/runtime`)
**Role:** Orchestrates the workflow.
**Implementation:** 
- **YAML Parser:** Translates declarative YAML steps into Go `Step` interfaces.
- **The Execution Loop:** Iterates through steps, handles timeouts, and orchestrates the fallback to the `Healer` if the `Resolver` fails.

### 2.5 The Persistence Layer (`memory/`)
**Role:** Stores the telemetry and fingerprints.
**Implementation:** Embedded SQLite (`mattn/go-sqlite3`).
**Decisions:** SQLite was chosen because it requires no separate database server, making Mr. Browser trivial to self-host. The connection pool is strictly limited to 1 (`SetMaxOpenConns(1)`) to avoid database locking issues during concurrent workflow execution.

---

## 3. Communication & SDKs

The engine is built as a standalone service exposing a REST API on port `7331`.
**Why REST instead of gRPC?**
- REST endpoints are trivial to consume from any language (Python, Node, Bash, Rust) without requiring Protobuf compilation.
- Browser automation is generally I/O bound, meaning the microsecond performance benefits of gRPC binary framing are negligible compared to network and rendering latency.

The SDKs (Python and TypeScript) are strictly "thin wrappers." They contain zero business logic. They simply construct JSON payloads and HTTP requests. This guarantees that behavior is exactly identical regardless of the language driving the browser.

---

## 4. Future Considerations & Roadmap

### 1. Vision Module (Phase 5)
Currently, intent resolution is purely DOM-based. If an element is rendered on a `<canvas>` (like a WebGL game or highly custom UI), the DOM tree provides no useful metadata. 
**Future Plan:** Integrate a lightweight local vision model (e.g., via OpenCV or a quantized local multimodal LLM) to resolve bounding boxes based on pixels rather than DOM nodes.

### 2. Multi-Tab Orchestration
The current `Executor` is scoped to a single `Page`. 
**Future Plan:** Expand the YAML schema to allow parallel execution across multiple tabs, and introduce a `tab: switch` step for complex SSO/OAuth flows that spawn new windows.

### 3. Distributed Execution
While SQLite is perfect for single-node deployment, scaling Mr. Browser across a Kubernetes cluster requires shared state.
**Future Plan:** Abstract the `memory.Store` interface to support Redis or Postgres for cluster-wide fingerprint sharing.
