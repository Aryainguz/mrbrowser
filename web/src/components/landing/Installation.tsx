import { useState } from "react";
import { Check, Copy } from "lucide-react";
import { FadeIn } from "@/components/FadeIn";
import { highlight } from "@/lib/highlight";

const TABS = [
  { id: "docker", label: "Docker", command: "docker compose -f docker/docker-compose.yml up -d", lang: "bash" },
  { id: "source", label: "From Source", command: "make build && make install", lang: "bash" },
  { id: "python", label: "Python SDK", command: "pip install mrbrowser", lang: "bash" },
  {
    id: "typescript",
    label: "TypeScript SDK",
    command: "npm install @mrbrowser/sdk",
    lang: "bash",
  },
] as const;

export function Installation() {
  const [active, setActive] = useState<string>("docker");
  const [copied, setCopied] = useState(false);
  const tab = TABS.find((t) => t.id === active) ?? TABS[0];

  const copy = async () => {
    await navigator.clipboard.writeText(tab.command);
    setCopied(true);
    setTimeout(() => setCopied(false), 1800);
  };

  return (
    <section id="install" className="mx-auto max-w-3xl scroll-mt-20 px-4 py-24 sm:px-6">
      <FadeIn>
        <div className="text-center">
          <p className="font-mono text-xs text-primary">$ init</p>
          <h2 className="mt-3 font-mono text-3xl font-bold tracking-tight sm:text-4xl">
            Deploy in <span className="text-primary text-glow">seconds</span>
          </h2>
          <p className="mt-3 text-sm text-muted-foreground">
            Pick your runtime. The engine ships everywhere.
          </p>
        </div>

        <div className="glass mt-10 overflow-hidden rounded-xl">
          <div className="flex border-b border-border">
            {TABS.map((t) => (
              <button
                key={t.id}
                onClick={() => setActive(t.id)}
                className={`flex-1 px-4 py-3 font-mono text-xs transition-all ${
                  active === t.id
                    ? "border-b-2 border-primary bg-primary/5 text-primary"
                    : "text-muted-foreground hover:text-foreground"
                }`}
              >
                {t.label}
              </button>
            ))}
          </div>
          <div className="flex items-center justify-between gap-4 p-5">
            <code className="font-mono text-sm text-foreground/90">
              <span className="mr-2 text-primary">$</span>
              {highlight(tab.command, tab.lang)}
            </code>
            <button
              onClick={copy}
              aria-label="Copy install command"
              className="flex items-center gap-2 rounded-md border border-primary/40 bg-primary/10 px-3 py-2 font-mono text-xs text-primary shadow-glow transition-all hover:bg-primary hover:text-primary-foreground hover:shadow-glow-strong"
            >
              {copied ? <Check size={13} /> : <Copy size={13} />}
              {copied ? "copied" : "copy"}
            </button>
          </div>
        </div>
      </FadeIn>
    </section>
  );
}
