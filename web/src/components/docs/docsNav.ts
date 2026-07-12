export interface DocPage {
  title: string;
  path: string;
  section: string;
}

export const DOC_PAGES: DocPage[] = [
  { title: "Getting Started", path: "/docs", section: "Introduction" },
  { title: "Core Concepts", path: "/docs/core-concepts", section: "Introduction" },
  { title: "Use Cases", path: "/docs/use-cases", section: "Guides" },
  { title: "YAML Workflow Reference", path: "/docs/yaml-workflows", section: "Guides" },
  { title: "CLI Reference", path: "/docs/cli", section: "Guides" },
  { title: "Python SDK", path: "/docs/python", section: "SDKs" },
  { title: "TypeScript SDK", path: "/docs/typescript", section: "SDKs" },
  { title: "Architecture", path: "/docs/architecture", section: "Reference" },
  { title: "Contributing", path: "/docs/contributing", section: "Project" },
  { title: "Changelog", path: "/docs/changelog", section: "Project" },
];

export const DOC_SECTIONS = ["Introduction", "Guides", "SDKs", "Reference", "Project"] as const;

export function getPagerLinks(path: string) {
  const idx = DOC_PAGES.findIndex((p) => p.path === path);
  return {
    prev: idx > 0 ? DOC_PAGES[idx - 1] : null,
    next: idx >= 0 && idx < DOC_PAGES.length - 1 ? DOC_PAGES[idx + 1] : null,
  };
}
