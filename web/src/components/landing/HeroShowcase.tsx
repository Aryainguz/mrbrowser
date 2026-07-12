import { useEffect, useState } from "react";
import { Check, Lock, RotateCw } from "lucide-react";
import { highlight } from "@/lib/highlight";

type Action = "fill-email" | "fill-password" | "click-login" | null;

const SCRIPT: { line: string; action: Action }[] = [
  { line: "name: login_admin", action: null },
  { line: "steps:", action: null },
  { line: '  - open:', action: null },
  { line: '      url: "https://corp-portal.internal"', action: null },
  { line: '  - type:', action: null },
  { line: '      target: "Email"', action: null },
  { line: '      value: "admin@corp.com"', action: "fill-email" },
  { line: "  - type:", action: null },
  { line: '      target: "Password"', action: null },
  { line: '      value: "$SECRET_PASS"', action: "fill-password" },
  { line: '  - click:', action: null },
  { line: '      target: "Login"', action: "click-login" },
];

interface BrowserState {
  email: string;
  password: string;
  highlight: "email" | "password" | "login" | null;
  pressed: boolean;
  success: boolean;
}

const INITIAL: BrowserState = {
  email: "",
  password: "",
  highlight: null,
  pressed: false,
  success: false,
};

const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

export function HeroShowcase() {
  const [lines, setLines] = useState<string[]>([]);
  const [current, setCurrent] = useState("");
  const [browser, setBrowser] = useState<BrowserState>(INITIAL);

  useEffect(() => {
    let cancelled = false;

    const typeText = async (target: "email" | "password", value: string) => {
      for (let i = 1; i <= value.length; i++) {
        if (cancelled) return;
        setBrowser((b) => ({ ...b, [target]: value.slice(0, i) }));
        await sleep(35);
      }
    };

    (async () => {
      while (!cancelled) {
        setLines([]);
        setCurrent("");
        setBrowser(INITIAL);
        await sleep(600);

        for (const step of SCRIPT) {
          for (let c = 1; c <= step.line.length; c++) {
            if (cancelled) return;
            setCurrent(step.line.slice(0, c));
            await sleep(26);
          }
          setLines((p) => [...p, step.line]);
          setCurrent("");

          if (step.action === "fill-email") {
            setBrowser((b) => ({ ...b, highlight: "email" }));
            await sleep(280);
            await typeText("email", "admin@corp.com");
            await sleep(350);
            setBrowser((b) => ({ ...b, highlight: null }));
          } else if (step.action === "fill-password") {
            setBrowser((b) => ({ ...b, highlight: "password" }));
            await sleep(280);
            await typeText("password", "••••••••••");
            await sleep(350);
            setBrowser((b) => ({ ...b, highlight: null }));
          } else if (step.action === "click-login") {
            setBrowser((b) => ({ ...b, highlight: "login" }));
            await sleep(400);
            setBrowser((b) => ({ ...b, pressed: true }));
            await sleep(450);
            setBrowser((b) => ({ ...b, pressed: false, highlight: null, success: true }));
          }
          await sleep(120);
        }
        await sleep(3200);
      }
    })();

    return () => {
      cancelled = true;
    };
  }, []);

  return (
    <div className="grid gap-6 lg:grid-cols-2">
      {/* Left: macOS terminal */}
      <div className="scanlines overflow-hidden rounded-xl border border-primary/25 bg-[oklch(0.1_0.003_145)] shadow-glow">
        <div className="flex items-center gap-3 border-b border-border px-4 py-3">
          <div className="flex gap-1.5">
            <span className="h-3 w-3 rounded-full bg-[oklch(0.62_0.2_25)]" />
            <span className="h-3 w-3 rounded-full bg-[oklch(0.8_0.16_85)]" />
            <span className="h-3 w-3 rounded-full bg-[oklch(0.72_0.19_145)]" />
          </div>
          <span className="font-mono text-xs text-muted-foreground">
            ~/flows/login.yaml — mrbrowser run
          </span>
        </div>
        <div className="h-[300px] overflow-hidden p-4 font-mono text-[13px] leading-6 sm:h-[320px]">
          {lines.map((l, i) => (
            <div key={i} className="whitespace-pre">
              {highlight(l, "yaml")}
            </div>
          ))}
          <div className="whitespace-pre">
            {highlight(current, "yaml")}
            <span className="caret-blink inline-block h-4 w-[7px] translate-y-[3px] bg-primary" />
          </div>
        </div>
      </div>

      {/* Right: mock browser */}
      <div className="overflow-hidden rounded-xl border border-border bg-[oklch(0.13_0.003_145)]">
        <div className="flex items-center gap-3 border-b border-border px-4 py-3">
          <div className="flex gap-1.5">
            <span className="h-3 w-3 rounded-full bg-muted" />
            <span className="h-3 w-3 rounded-full bg-muted" />
            <span className="h-3 w-3 rounded-full bg-muted" />
          </div>
          <div className="flex flex-1 items-center gap-2 rounded-md bg-[oklch(0.09_0_0)] px-3 py-1.5">
            <Lock size={11} className="text-muted-foreground" />
            <span className="font-mono text-[11px] text-muted-foreground">
              corp-portal.internal/login
            </span>
            <RotateCw size={11} className="ml-auto text-muted-foreground" />
          </div>
        </div>

        <div className="flex h-[300px] items-center justify-center p-6 sm:h-[320px]">
          {browser.success ? (
            <div className="flex flex-col items-center gap-3 text-center">
              <div className="flex h-14 w-14 items-center justify-center rounded-full border border-primary/50 bg-primary/10 shadow-glow">
                <Check size={26} className="text-primary" />
              </div>
              <p className="font-mono text-sm font-semibold text-primary text-glow">
                AUTHENTICATED
              </p>
              <p className="font-mono text-[11px] text-muted-foreground">
                session established · 2 elements resolved by intent
              </p>
            </div>
          ) : (
            <div className="w-full max-w-xs space-y-4">
              <p className="text-center font-sans text-sm font-semibold text-foreground/80">
                Sign in to Corp Portal
              </p>
              <div className="space-y-1">
                <label className="font-sans text-[11px] font-medium text-muted-foreground">
                  Email
                </label>
                <div
                  className={`flex h-9 items-center rounded-md border bg-[oklch(0.09_0_0)] px-3 font-sans text-xs transition-all duration-300 ${
                    browser.highlight === "email"
                      ? "border-primary shadow-glow"
                      : "border-border"
                  }`}
                >
                  <span className="text-foreground/90">{browser.email}</span>
                  {browser.highlight === "email" && (
                    <span className="caret-blink ml-px inline-block h-3.5 w-px bg-primary" />
                  )}
                </div>
              </div>
              <div className="space-y-1">
                <label className="font-sans text-[11px] font-medium text-muted-foreground">
                  Password
                </label>
                <div
                  className={`flex h-9 items-center rounded-md border bg-[oklch(0.09_0_0)] px-3 font-sans text-xs transition-all duration-300 ${
                    browser.highlight === "password"
                      ? "border-primary shadow-glow"
                      : "border-border"
                  }`}
                >
                  <span className="tracking-widest text-foreground/90">{browser.password}</span>
                  {browser.highlight === "password" && (
                    <span className="caret-blink ml-px inline-block h-3.5 w-px bg-primary" />
                  )}
                </div>
              </div>
              <button
                className={`relative h-9 w-full overflow-hidden rounded-md font-sans text-xs font-semibold transition-all duration-150 ${
                  browser.highlight === "login"
                    ? "bg-primary text-primary-foreground shadow-glow-strong"
                    : "bg-primary/80 text-primary-foreground"
                } ${browser.pressed ? "scale-95" : "scale-100"}`}
              >
                Login
                {browser.pressed && (
                  <span className="animate-ripple absolute left-1/2 top-1/2 h-10 w-10 -translate-x-1/2 -translate-y-1/2 rounded-full bg-white/40" />
                )}
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
