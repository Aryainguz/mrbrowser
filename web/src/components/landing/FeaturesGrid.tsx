import { useRef, type MouseEvent, type ReactNode } from "react";
import { BrainCircuit, Dna, Lock } from "lucide-react";
import { FadeIn } from "@/components/FadeIn";

const FEATURES = [
  {
    icon: BrainCircuit,
    title: "Intent-Driven Resolution",
    body: 'No CSS selectors. No XPath. Describe the element in plain English — "the submit button", "the email field" — and a local NLP algorithm resolves it against the DOM Accessibility Tree.',
  },
  {
    icon: Dna,
    title: "Memory Engine: Self-Healing",
    body: "Every resolved element leaves a structural fingerprint. When the UI changes, the Memory Engine compares historical fingerprints to find where the element moved — and your script keeps running.",
  },
  {
    icon: Lock,
    title: "Zero Mandatory AI",
    body: "The core resolution engine is 100% local and deterministic. No API keys, no cloud calls, no data leaving your machine. LLM assistance is strictly opt-in.",
  },
];

function GlowCard({ children }: { children: ReactNode }) {
  const ref = useRef<HTMLDivElement>(null);

  const onMove = (e: MouseEvent<HTMLDivElement>) => {
    const el = ref.current;
    if (!el) return;
    const r = el.getBoundingClientRect();
    el.style.setProperty("--mx", `${e.clientX - r.left}px`);
    el.style.setProperty("--my", `${e.clientY - r.top}px`);
  };

  return (
    <div
      ref={ref}
      onMouseMove={onMove}
      className="group relative rounded-xl border border-border p-px transition-colors"
      style={{
        background:
          "radial-gradient(240px circle at var(--mx, 50%) var(--my, 50%), oklch(0.87 0.28 143 / 0.5), transparent 70%)",
      }}
    >
      <div
        className="relative h-full rounded-[11px] p-6"
        style={{ background: "var(--gradient-card)" }}
      >
        <div
          className="pointer-events-none absolute inset-0 rounded-[11px] opacity-0 transition-opacity duration-300 group-hover:opacity-100"
          style={{
            background:
              "radial-gradient(300px circle at var(--mx, 50%) var(--my, 50%), oklch(0.87 0.28 143 / 0.08), transparent 70%)",
          }}
        />
        {children}
      </div>
    </div>
  );
}

export function FeaturesGrid() {
  return (
    <section id="features" className="mx-auto max-w-6xl scroll-mt-20 px-4 py-24 sm:px-6">
      <FadeIn>
        <p className="font-mono text-xs text-primary">$ mrbrowser --capabilities</p>
        <h2 className="mt-3 font-mono text-3xl font-bold tracking-tight sm:text-4xl">
          Built different. <span className="text-primary text-glow">On purpose.</span>
        </h2>
      </FadeIn>
      <div className="mt-12 grid gap-6 md:grid-cols-3">
        {FEATURES.map((f, i) => (
          <FadeIn key={f.title} delay={i * 0.12}>
            <GlowCard>
              <div className="flex h-11 w-11 items-center justify-center rounded-lg border border-primary/30 bg-primary/10">
                <f.icon size={22} className="text-primary" />
              </div>
              <h3 className="mt-5 font-mono text-base font-semibold text-foreground">
                {f.title}
              </h3>
              <p className="mt-3 text-sm leading-relaxed text-muted-foreground">{f.body}</p>
            </GlowCard>
          </FadeIn>
        ))}
      </div>
    </section>
  );
}
