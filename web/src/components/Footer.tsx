import { Link } from "@tanstack/react-router";

export function Footer() {
  return (
    <footer className="border-t border-border bg-[oklch(0.09_0_0)]">
      <div className="mx-auto flex max-w-6xl flex-col items-center justify-between gap-6 px-4 py-10 sm:flex-row sm:px-6">
        <div className="flex flex-col items-center gap-1 sm:items-start">
          <Link to="/" className="font-mono text-sm font-bold text-primary text-glow transition-opacity hover:opacity-80">
            [Mr. Browser]
          </Link>
          <span className="font-mono text-[11px] text-muted-foreground">
            © 2026 Mr. Browser. Open source under MIT.
          </span>
        </div>
        <div className="flex items-center gap-6 font-mono text-xs text-muted-foreground">
          <a href="/#features" className="transition-colors hover:text-primary">
            features
          </a>
          <a href="/#use-cases" className="transition-colors hover:text-primary">
            use_cases
          </a>
          <Link to="/docs" className="transition-colors hover:text-primary">
            docs
          </Link>
          <a
            href="https://github.com"
            target="_blank"
            rel="noreferrer"
            className="transition-colors hover:text-primary"
          >
            github
          </a>
        </div>
      </div>
    </footer>
  );
}
