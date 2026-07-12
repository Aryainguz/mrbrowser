import { ReactNode } from "react";

type Token = { text: string; cls?: string };

type Rule = { re: RegExp; cls: string };

const RULES: Record<string, Rule[]> = {
  yaml: [
    { re: /^#[^\n]*/, cls: "tok-c" },
    { re: /^"(?:[^"\\\n]|\\.)*"|^'[^'\n]*'/, cls: "tok-s" },
    { re: /^[A-Za-z_][\w-]*(?=\s*:)/, cls: "tok-key" },
    { re: /^\d+(\.\d+)?/, cls: "tok-n" },
    { re: /^- /, cls: "tok-k" },
  ],
  python: [
    { re: /^#[^\n]*/, cls: "tok-c" },
    { re: /^("""[\s\S]*?"""|"(?:[^"\\\n]|\\.)*"|'(?:[^'\\\n]|\\.)*')/, cls: "tok-s" },
    { re: /^@[\w.]+/, cls: "tok-f" },
    {
      re: /^(def|class|import|from|return|assert|with|as|for|in|if|elif|else|not|and|or|await|async|True|False|None|lambda|yield|raise|try|except)\b/,
      cls: "tok-k",
    },
    { re: /^[A-Za-z_]\w*(?=\()/, cls: "tok-f" },
    { re: /^\d+(\.\d+)?/, cls: "tok-n" },
  ],
  typescript: [
    { re: /^\/\/[^\n]*|^\/\*[\s\S]*?\*\//, cls: "tok-c" },
    { re: /^`(?:[^`\\]|\\.)*`|^"(?:[^"\\\n]|\\.)*"|^'(?:[^'\\\n]|\\.)*'/, cls: "tok-s" },
    {
      re: /^(import|from|export|const|let|var|await|async|function|return|new|type|interface|class|extends|implements|if|else|for|of|in|try|catch|throw|default|true|false|null|undefined)\b/,
      cls: "tok-k",
    },
    { re: /^[A-Za-z_$][\w$]*(?=\()/, cls: "tok-f" },
    { re: /^\d+(\.\d+)?/, cls: "tok-n" },
  ],
  bash: [
    { re: /^#[^\n]*/, cls: "tok-c" },
    { re: /^"(?:[^"\\\n]|\\.)*"|^'[^'\n]*'/, cls: "tok-s" },
    { re: /^(docker-compose|docker|pip|npm|npx|bun|mrbrowser|curl|git)\b/, cls: "tok-k" },
    { re: /^--?[\w-]+/, cls: "tok-key" },
    { re: /^\$ /, cls: "tok-c" },
  ],
};

export function tokenize(code: string, lang: string): Token[] {
  const rules = RULES[lang] ?? [];
  const out: Token[] = [];
  let i = 0;
  while (i < code.length) {
    const rest = code.slice(i);
    let matched = false;
    for (const rule of rules) {
      const m = rule.re.exec(rest);
      if (m && m.index === 0 && m[0].length > 0) {
        out.push({ text: m[0], cls: rule.cls });
        i += m[0].length;
        matched = true;
        break;
      }
    }
    if (!matched) {
      const w = /^[A-Za-z0-9_]+/.exec(rest);
      const text = w ? w[0] : code[i];
      out.push({ text });
      i += text.length;
    }
  }
  return out;
}

export function highlight(code: string, lang: string): ReactNode[] {
  return tokenize(code, lang).map((t, idx) =>
    t.cls ? (
      <span key={idx} className={t.cls}>
        {t.text}
      </span>
    ) : (
      <span key={idx}>{t.text}</span>
    ),
  );
}
