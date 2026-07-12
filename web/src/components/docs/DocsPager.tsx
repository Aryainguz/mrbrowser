import { Link } from "@tanstack/react-router";
import { ArrowLeft, ArrowRight } from "lucide-react";
import { getPagerLinks } from "./docsNav";

export function DocsPager({ path }: { path: string }) {
  const { prev, next } = getPagerLinks(path);

  return (
    <nav className="mt-16 grid gap-4 border-t border-border pt-8 sm:grid-cols-2">
      {prev ? (
        <Link
          to={prev.path}
          className="group flex flex-col gap-2 rounded-xl border border-border p-6 transition-all hover:-translate-y-0.5 hover:border-primary/50 hover:shadow-glow"
        >
          <span className="flex items-center gap-2 font-mono text-[11px] text-muted-foreground">
            <ArrowLeft size={13} className="transition-transform group-hover:-translate-x-1" />
            PREVIOUS
          </span>
          <span className="font-mono text-base font-semibold text-foreground transition-colors group-hover:text-primary">
            {prev.title}
          </span>
        </Link>
      ) : (
        <div />
      )}
      {next ? (
        <Link
          to={next.path}
          className="group flex flex-col items-end gap-2 rounded-xl border border-border p-6 text-right transition-all hover:-translate-y-0.5 hover:border-primary/50 hover:shadow-glow"
        >
          <span className="flex items-center gap-2 font-mono text-[11px] text-muted-foreground">
            NEXT
            <ArrowRight size={13} className="transition-transform group-hover:translate-x-1" />
          </span>
          <span className="font-mono text-base font-semibold text-foreground transition-colors group-hover:text-primary">
            {next.title}
          </span>
        </Link>
      ) : (
        <div />
      )}
    </nav>
  );
}
