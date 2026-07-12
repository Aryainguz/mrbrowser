import { createFileRoute, Link } from "@tanstack/react-router";
import { ArrowRight, Terminal } from "lucide-react";
import { Navbar } from "@/components/Navbar";
import { Footer } from "@/components/Footer";
import { MatrixRain } from "@/components/MatrixRain";
import { TypingHeadline } from "@/components/landing/TypingHeadline";
import { HeroShowcase } from "@/components/landing/HeroShowcase";
import { FeaturesGrid } from "@/components/landing/FeaturesGrid";
import { UseCases } from "@/components/landing/UseCases";
import { Installation } from "@/components/landing/Installation";
import { FadeIn } from "@/components/FadeIn";

export const Route = createFileRoute("/")({
  head: () => ({
    meta: [
      { title: "Mr. Browser — Forget Selectors. Command with Intent." },
      {
        name: "description",
        content:
          "Open-source browser automation engine that resolves elements from plain-English intent. Self-healing scripts, local NLP, zero mandatory AI.",
      },
      { property: "og:title", content: "Mr. Browser — Forget Selectors. Command with Intent." },
      {
        property: "og:description",
        content:
          "Intent-driven, self-healing, local-first browser automation for QA, RPA, and scraping.",
      },
    ],
  }),
  component: Index,
});

function Index() {
  return (
    <div className="min-h-screen bg-background font-sans">
      <Navbar />

      {/* Hero */}
      <section className="scanlines relative overflow-hidden pb-20 pt-32 sm:pt-40">
        <MatrixRain className="absolute inset-0 h-full w-full opacity-40" />
        <div
          className="pointer-events-none absolute inset-0"
          style={{
            background:
              "radial-gradient(ellipse 80% 60% at 50% 0%, oklch(0.87 0.28 143 / 0.07), transparent 60%), linear-gradient(to bottom, transparent 40%, oklch(0.08 0 0) 95%)",
          }}
        />
        <div className="relative z-10 mx-auto max-w-6xl px-4 sm:px-6">
          <div className="text-center">
            <div className="mb-6 inline-flex items-center gap-2 rounded-full border border-primary/30 bg-primary/5 px-4 py-1.5 font-mono text-[11px] text-primary">
              <Terminal size={12} />
              open-source · v2.4.0 · MIT
            </div>
            <TypingHeadline />
            <p className="mx-auto mt-6 max-w-2xl text-sm leading-relaxed text-muted-foreground sm:text-base">
              Mr. Browser is an automation engine that abandons fragile CSS selectors.
              Describe what you want in plain English — a local NLP algorithm and the DOM
              Accessibility Tree do the rest. When the UI changes, scripts heal themselves.
            </p>
            <div className="mt-8 flex items-center justify-center gap-4">
              <Link
                to="/docs"
                className="group flex items-center gap-2 rounded-md border border-primary bg-primary px-6 py-3 font-mono text-sm font-semibold text-primary-foreground shadow-glow-strong transition-all hover:shadow-glow"
              >
                Get Started
                <ArrowRight size={15} className="transition-transform group-hover:translate-x-1" />
              </Link>
              <a
                href="#install"
                className="rounded-md border border-border px-6 py-3 font-mono text-sm text-muted-foreground transition-all hover:border-primary/50 hover:text-primary"
              >
                $ install
              </a>
            </div>
          </div>

          <FadeIn className="mt-20" delay={0.1}>
            <HeroShowcase />
          </FadeIn>
        </div>
      </section>

      <FeaturesGrid />
      <UseCases />
      <Installation />

      {/* CTA */}
      <section className="border-t border-border bg-[oklch(0.09_0_0)] py-20">
        <FadeIn className="mx-auto max-w-2xl px-4 text-center sm:px-6">
          <h2 className="font-mono text-2xl font-bold sm:text-3xl">
            <span className="text-primary text-glow">&gt;</span> Ready to stop fixing
            selectors?
          </h2>
          <p className="mt-3 text-sm text-muted-foreground">
            Read the docs, run your first flow in under five minutes.
          </p>
          <Link
            to="/docs"
            className="mt-8 inline-flex items-center gap-2 rounded-md border border-primary/60 bg-primary/10 px-8 py-3 font-mono text-sm font-semibold text-primary shadow-glow transition-all hover:bg-primary hover:text-primary-foreground hover:shadow-glow-strong"
          >
            Open Documentation <ArrowRight size={15} />
          </Link>
        </FadeIn>
      </section>

      <Footer />
    </div>
  );
}
