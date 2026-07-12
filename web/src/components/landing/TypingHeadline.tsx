import { useEffect, useState } from "react";

const PHRASE_1 = "Forget Selectors.";
const PHRASE_2 = "Automate in plain English.";
const GLITCH_CHARS = "!<>-_\\/[]{}—=+*^?#$%01";

const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

export function TypingHeadline() {
  const [text, setText] = useState("");
  const [glitching, setGlitching] = useState(false);
  const [phase, setPhase] = useState<1 | 2>(1);

  useEffect(() => {
    let cancelled = false;

    const type = async (phrase: string) => {
      for (let i = 1; i <= phrase.length; i++) {
        if (cancelled) return;
        setText(phrase.slice(0, i));
        await sleep(65);
      }
    };

    const glitch = async (phrase: string) => {
      setGlitching(true);
      for (let f = 0; f < 8; f++) {
        if (cancelled) return;
        setText(
          phrase
            .split("")
            .map((c) =>
              c === " " || Math.random() > 0.35
                ? c
                : GLITCH_CHARS[Math.floor(Math.random() * GLITCH_CHARS.length)],
            )
            .join(""),
        );
        await sleep(55);
      }
      setGlitching(false);
    };

    const erase = async () => {
      let current = "";
      setText((t) => {
        current = t;
        return t;
      });
      for (let i = current.length; i >= 0; i--) {
        if (cancelled) return;
        setText(current.slice(0, i));
        await sleep(22);
      }
    };

    (async () => {
      while (!cancelled) {
        setPhase(1);
        await type(PHRASE_1);
        await sleep(1400);
        await glitch(PHRASE_1);
        await erase();
        await sleep(300);
        setPhase(2);
        await type(PHRASE_2);
        await sleep(4200);
        await glitch(PHRASE_2);
        await erase();
        await sleep(400);
      }
    })();

    return () => {
      cancelled = true;
    };
  }, []);

  return (
    <h1 className="font-mono text-4xl font-bold leading-tight tracking-tight sm:text-5xl lg:text-6xl">
      <span
        className={
          glitching
            ? "text-destructive text-glow-red"
            : phase === 2
              ? "text-primary text-glow"
              : "text-foreground"
        }
      >
        {text}
      </span>
      <span className="caret-blink ml-1 inline-block h-[0.9em] w-[0.5ch] translate-y-[0.12em] bg-primary align-baseline" />
    </h1>
  );
}
