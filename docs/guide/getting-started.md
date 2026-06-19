# Mr. Browser Guide

Welcome to the definitive guide for Mr. Browser. This document will walk you through the core concepts, common use-cases, and how to harness the full power of intent-driven automation.

---

## What is Mr. Browser?

Traditional automation tools (Selenium, Playwright, Puppeteer) require you to find elements using CSS selectors or XPaths (`#login-form > div:nth-child(2) > input`). When the website updates its UI, the selectors break, and your automation fails.

**Mr. Browser solves this by understanding the page.**
You provide an "intent" (e.g., `"Email field"` or `"Submit"`). Mr. Browser parses the page's accessibility tree and interactive elements, mathematically scores them against your intent, and interacts with the best match. If the element later changes so drastically that the intent fails, Mr. Browser's **Memory Engine** uses historical DOM fingerprints to locate where the element moved, healing the script automatically.

---

## Use Cases

### 1. Resilient E2E Testing
QA teams spend countless hours fixing broken end-to-end tests due to trivial UI changes (like a developer changing a button from `<button id="submit">` to `<button class="btn-primary">`). Mr. Browser's self-healing ensures tests only fail when business logic actually breaks, drastically reducing maintenance overhead.

### 2. Intelligent Web Scraping
Targeting data on heavily obfuscated websites is difficult. With Mr. Browser, you can navigate and extract data semantically without needing to reverse-engineer minified React/Webpack class names.

### 3. RPA (Robotic Process Automation)
Automating repetitive internal tasks across legacy enterprise software. Since you only need to describe the element (e.g., `"Download Invoice"`), non-technical operators can write YAML workflows without knowing HTML or CSS.

---

## Writing YAML Workflows

The simplest way to use Mr. Browser is via the built-in YAML execution engine. Workflows are designed to be entirely readable.

### Schema Overview

A workflow file contains a name, description, and an array of `steps`.

```yaml
name: invoice-automation
description: "Logs in and downloads the latest invoice"

steps:
  - open:
      url: https://billing.example.com
      
  - type:
      target: "Email Address"
      value: "admin@example.com"
      
  - type:
      target: "Password"
      value: "supersecret"
      enter: true            # Automatically presses the Enter key after typing
      
  - wait:
      seconds: 2             # Explicit wait for slow SPA transitions
      
  - click:
      target: "Download Latest Invoice"
      
  - assert:
      text_visible: "Download Complete"
```

### Advanced Actions

- **Extracting Data:**
  You can extract text from the screen dynamically using the `extract` step, which saves it to the session memory.
  ```yaml
  - extract:
      target: "Account Balance"
      save_as: "balance"
  ```

- **Scrolling:**
  ```yaml
  - scroll:
      direction: down
      pixels: 800
  ```

## Running the Engine

### Starting the Server
Mr. Browser operates as a standalone server. The easiest way to run it is via the provided Docker Compose configuration, which guarantees a clean, isolated Chromium environment.

```bash
docker-compose up -d
```
*The API will be available on `http://localhost:7331`.*

### Executing Workflows
Use the Mr. Browser CLI to run your YAML files against the engine:

```bash
# Run headless (default)
go run ./cli run my-workflow.yaml

# Run with headed browser to watch it happen
go run ./cli run my-workflow.yaml --headed
```
