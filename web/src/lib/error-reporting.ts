// Error reporting — logs errors to the console.
export function reportError(error: unknown, _context: Record<string, unknown> = {}) {
  console.error("[mr-browser]", error);
}
