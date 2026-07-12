# Core Concepts

Mr. Browser introduces a paradigm shift in how we approach browser automation. This page covers the underlying concepts that make the engine intelligent, resilient, and fast.

---

## 1. Intent-Driven Resolution

Unlike traditional tools that rely on fragile CSS selectors or XPath expressions, Mr. Browser uses **Intent-Driven Resolution**.

When you provide a target string (e.g., `"Sign Up Button"`), the engine:
1. Captures the full DOM and Accessibility Tree of the current page.
2. Identifies all interactive elements (buttons, links, inputs).
3. Uses a local NLP matching algorithm to score elements based on their text content, ARIA labels, and visual proximity.
4. Selects the highest-scoring element that matches the intent.

Because this algorithm runs locally in Go without any LLM round-trips, it is extremely fast and completely private.

---

## 2. The Memory Engine (Self-Healing)

Websites change constantly. A developer might change an element's class, wrap it in a new `<div>`, or change the button text slightly.

When an element is successfully interacted with for the first time, Mr. Browser's **Memory Engine** saves a structural fingerprint of that element to `mrbrowser.db` (a local SQLite database). 

If a workflow runs again and the intent string no longer confidently matches any element (e.g., because the UI was heavily redesigned), the Memory Engine falls back to the saved fingerprint. It analyzes the new DOM structure and finds the element that most closely resembles the historical fingerprint, allowing the script to "self-heal" and continue executing.

---

## 3. Thin-Client SDK Architecture

Mr. Browser's architecture is separated into a **Core Engine** (written in Go) and **Thin-Client SDKs** (Python, TypeScript).

- **Core Engine:** Handles Chrome DevTools Protocol (CDP) orchestration, element resolution, self-healing memory, and screenshot generation.
- **SDKs:** Lightweight wrappers around HTTP/REST calls. They maintain session state but contain zero heavy dependencies, meaning you don't have to deal with installing Chromium or massive Playwright/Puppeteer binaries in your client application.

---

## 4. Workflows

Workflows are YAML files that define a sequence of automation steps. They are perfect for non-developers who want to automate tasks, or for orchestrating simple, linear data-extraction jobs without needing to write full Python or TypeScript code.

The Engine parses these YAML files and converts them into sequential CDP commands automatically.
