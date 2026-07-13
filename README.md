<img width="601" height="404" alt="mrbrowser" src="https://github.com/user-attachments/assets/c035cc7b-9c68-4158-b3b8-8a13637c381f" />


---


[![Go](https://img.shields.io/badge/Go-1.22+-blue?logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Self-Hosted](https://img.shields.io/badge/self--hosted-friendly-brightgreen)]()

**Mr. Browser** is an open-source, self-hosted browser automation intelligence engine.

It replaces fragile CSS selector-based automation with intent-driven, self-healing browser control.

### See it in action (Fifa World Cup Latest YouTube Videos & Morning Tech Briefing Workflow)


https://github.com/user-attachments/assets/c9709fe8-227c-4d49-952d-ea5557bc9322



```
User intent
    │
Page understanding
    │
Element discovery  ← scored candidates, no selectors required
    │
Action             ← click / type / scroll / upload / download
    │
Verification       ← DOM diff, URL change detection
    │
Recovery           ← fingerprint-based self-healing
```

---

## Features

- 🧠 **Intent-based element resolution** — `click("login button")` instead of `click("#auth-btn-v2")`
- 🔧 **Self-healing automation** — fingerprints elements; recovers when the DOM changes
- 📄 **YAML task engine** — declarative, human-readable workflows
- 🚀 **Zero mandatory AI** — pure algorithmic scoring; AI is optional
- 🐳 **Docker-first** — minimal images, production-ready
- 🐍 **Python SDK** — clean client library for Python workflows
- 🔌 **REST API** — control the engine remotely via HTTP

---

## Quick Start

### Run with Docker

```bash
git clone https://github.com/Aryainguz/mrbrowser
cd mrbrowser
cp docker/.env.example docker/.env
docker compose -f docker/docker-compose.yml up
```

### Run locally

```bash
# Prerequisites: Go 1.22+, Chromium
go build -o mr-browser ./cli
./mr-browser screenshot https://example.com
./mr-browser run examples/login.yaml
./mr-browser inspect https://example.com
```

### Python SDK

```bash
pip install mrbrowser
```

```python
from mrbrowser import MrBrowser

browser = MrBrowser(host="localhost", port=7331)
page = browser.open("https://example.com")
page.click("login button")
page.type("email field", "user@example.com")
page.type("password field", "secret")
page.click("sign in")
screenshot = page.screenshot()
browser.close()
```

---

## CLI Commands

```bash
mr-browser run workflow.yaml          # Execute a task workflow
mr-browser inspect https://url        # Print page element tree
mr-browser screenshot https://url     # Capture a screenshot
mr-browser debug workflow.yaml        # Step-through with pause prompts
```

---

## YAML Workflow

```yaml
task:
  name: download_invoice

steps:
  - open:
      url: https://example.com

  - click:
      target: "login button"

  - type:
      target: "email field"
      value: "user@example.com"

  - type:
      target: "password field"
      value: "secret"

  - click:
      target: "download invoice"

  - screenshot:
      output: invoice_confirmation.png
```

---

## Architecture

```
Mr Browser SDK (Python / REST)
          │
    Task Runtime (Go)
          │
  ┌───────┴────────┐
  │                │
Browser         Intelligence
Controller      Engine
  │                │
Chrome CDP    ┌────┴────┐
  │           │         │
Chromium   DOM       Element
         Analyzer   Resolver
               │
          Action Executor
               │
          Memory Layer (SQLite)
```

---

## Configuration

Copy `.env.example` and configure:

```env
MRBROWSER_PORT=7331
MRBROWSER_CHROMIUM_PATH=/usr/bin/chromium
MRBROWSER_HEADLESS=true
MRBROWSER_DB_PATH=./mrbrowser.db
MRBROWSER_LOG_LEVEL=info
```

---

## Development

```bash
make build       # Compile binary
make test        # Run all tests
make test-unit   # Unit tests only
make test-int    # Integration tests (needs Chromium)
make docker      # Build Docker image
make lint        # Run linter
```

---

## License

MIT License — see [LICENSE](LICENSE).
