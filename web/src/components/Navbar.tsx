import { Link } from "@tanstack/react-router";

export function Navbar() {
  return (
    <header className="fixed inset-x-0 top-0 z-50 glass border-x-0 border-t-0">
      <nav className="mx-auto flex h-14 max-w-6xl items-center justify-between px-4 sm:px-6">
        <Link
          to="/"
          className="font-mono text-sm font-bold text-primary text-glow transition-opacity hover:opacity-80"
        >
          [Mr. Browser]
        </Link>
        <div className="flex items-center gap-4 sm:gap-7">
          <a
            href="/#features"
            className="hidden font-mono text-xs text-muted-foreground transition-colors hover:text-primary sm:block"
          >
            Features
          </a>
          <a
            href="/#use-cases"
            className="hidden font-mono text-xs text-muted-foreground transition-colors hover:text-primary sm:block"
          >
            Use Cases
          </a>
          <Link
            to="/docs"
            className="font-mono text-xs text-muted-foreground transition-colors hover:text-primary"
          >
            Docs
          </Link>
          <a
            href="https://github.com/mrbrowser/mrbrowser"
            target="_blank"
            rel="noreferrer"
            className="flex items-center gap-2 rounded-md border border-primary/60 bg-primary/10 px-4 py-1.5 font-mono text-xs font-semibold text-primary shadow-glow transition-all hover:bg-primary hover:text-primary-foreground hover:shadow-glow-strong"
          >
            GitHub
          </a>
        </div>
      </nav>
    </header>
  );
}
