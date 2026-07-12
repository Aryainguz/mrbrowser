import { AlertTriangle, Lightbulb } from "lucide-react";
import type { ReactNode } from "react";

export function Callout({
  type,
  title,
  children,
}: {
  type: "tip" | "warning";
  title?: string;
  children: ReactNode;
}) {
  const isTip = type === "tip";
  return (
    <div
      className={`my-6 rounded-lg border p-4 ${
        isTip
          ? "border-primary/40 bg-primary/5"
          : "border-destructive/50 bg-destructive/5"
      }`}
    >
      <div className="flex items-center gap-2">
        {isTip ? (
          <Lightbulb size={15} className="text-primary" />
        ) : (
          <AlertTriangle size={15} className="text-destructive" />
        )}
        <span
          className={`font-mono text-xs font-semibold uppercase tracking-wider ${
            isTip ? "text-primary" : "text-destructive"
          }`}
        >
          {title ?? (isTip ? "Tip" : "Warning")}
        </span>
      </div>
      <div className="mt-2 text-sm leading-relaxed text-foreground/85">{children}</div>
    </div>
  );
}
