import { useState } from "react";
import { Check, Copy } from "lucide-react";
import { highlight } from "@/lib/highlight";

interface CodeBlockProps {
  code: string;
  lang: string;
  title?: string;
  className?: string;
}

export function CodeBlock({ code, lang, title, className }: CodeBlockProps) {
  const [copied, setCopied] = useState(false);

  const copy = async () => {
    await navigator.clipboard.writeText(code);
    setCopied(true);
    setTimeout(() => setCopied(false), 1800);
  };

  return (
    <div
      className={`overflow-hidden rounded-lg border border-border bg-[oklch(0.11_0.003_145)] ${className ?? ""}`}
    >
      <div className="flex items-center justify-between border-b border-border px-4 py-2.5">
        <div className="flex items-center gap-3">
          <div className="flex gap-1.5">
            <span className="h-3 w-3 rounded-full bg-[oklch(0.62_0.2_25)]" />
            <span className="h-3 w-3 rounded-full bg-[oklch(0.8_0.16_85)]" />
            <span className="h-3 w-3 rounded-full bg-[oklch(0.72_0.19_145)]" />
          </div>
          {title && (
            <span className="font-mono text-xs text-muted-foreground">{title}</span>
          )}
        </div>
        <div className="flex items-center gap-3">
          <span className="font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
            {lang}
          </span>
          <button
            onClick={copy}
            aria-label="Copy code"
            className="rounded-md border border-border p-1.5 text-muted-foreground transition-all hover:border-primary/50 hover:text-primary hover:shadow-glow"
          >
            {copied ? <Check size={13} className="text-primary" /> : <Copy size={13} />}
          </button>
        </div>
      </div>
      <pre className="overflow-x-auto p-4 font-mono text-[13px] leading-relaxed text-foreground/90">
        <code>{highlight(code, lang)}</code>
      </pre>
    </div>
  );
}
