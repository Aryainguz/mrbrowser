import { useEffect, useRef, useState } from "react";
import { useNavigate } from "@tanstack/react-router";
import { ChevronRight } from "lucide-react";
import { DOC_PAGES } from "./docsNav";

export function DocsSearch() {
  const [query, setQuery] = useState("");
  const [focused, setFocused] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === "k") {
        e.preventDefault();
        inputRef.current?.focus();
      }
      if (e.key === "Escape") inputRef.current?.blur();
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, []);

  const results = query
    ? DOC_PAGES.filter((p) =>
        `${p.title} ${p.section}`.toLowerCase().includes(query.toLowerCase()),
      )
    : [];

  return (
    <div className="relative w-full max-w-md">
      <div className="glass flex items-center gap-2 rounded-lg px-3 py-2 font-mono text-sm transition-all focus-within:border-primary/50 focus-within:shadow-glow">
        <span className="text-primary">&gt;</span>
        <input
          ref={inputRef}
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onFocus={() => setFocused(true)}
          onBlur={() => setTimeout(() => setFocused(false), 150)}
          placeholder="search_docs..."
          className="w-full bg-transparent text-foreground placeholder:text-muted-foreground focus:outline-none"
        />
        <kbd className="hidden shrink-0 rounded border border-border bg-muted px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground sm:block">
          ⌘K
        </kbd>
      </div>
      {focused && results.length > 0 && (
        <div className="glass absolute top-full z-50 mt-2 w-full overflow-hidden rounded-lg">
          {results.map((r) => (
            <button
              key={r.path}
              onMouseDown={() => {
                navigate({ to: r.path });
                setQuery("");
              }}
              className="flex w-full items-center gap-2 px-4 py-2.5 text-left font-mono text-xs text-foreground transition-colors hover:bg-primary/10 hover:text-primary"
            >
              <ChevronRight size={12} className="text-primary" />
              {r.title}
              <span className="ml-auto text-[10px] text-muted-foreground">{r.section}</span>
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
