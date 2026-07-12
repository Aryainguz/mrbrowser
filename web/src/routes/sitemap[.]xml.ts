import { createFileRoute } from "@tanstack/react-router";
import type {} from "@tanstack/react-start";

// TODO: replace with your project URL once a project name or custom domain is set.
const BASE_URL = "";

interface SitemapEntry {
  path: string;
  lastmod?: string;
  changefreq?: "always" | "hourly" | "daily" | "weekly" | "monthly" | "yearly" | "never";
  priority?: string;
}

export const Route = createFileRoute("/sitemap.xml")({
  server: {
    handlers: {
      GET: async () => {
        const entries: SitemapEntry[] = [
          { path: "/", changefreq: "weekly", priority: "1.0" },
          { path: "/docs", changefreq: "weekly", priority: "0.9" },
          { path: "/docs/core-concepts", changefreq: "monthly", priority: "0.8" },
          { path: "/docs/use-cases", changefreq: "monthly", priority: "0.8" },
          { path: "/docs/yaml-workflows", changefreq: "monthly", priority: "0.8" },
          { path: "/docs/cli", changefreq: "monthly", priority: "0.8" },
          { path: "/docs/python", changefreq: "monthly", priority: "0.8" },
          { path: "/docs/typescript", changefreq: "monthly", priority: "0.8" },
          { path: "/docs/architecture", changefreq: "monthly", priority: "0.7" },
          { path: "/docs/contributing", changefreq: "monthly", priority: "0.7" },
          { path: "/docs/changelog", changefreq: "weekly", priority: "0.7" },
        ];

        const urls = entries.map((e) =>
          [
            `  <url>`,
            `    <loc>${BASE_URL}${e.path}</loc>`,
            e.lastmod ? `    <lastmod>${e.lastmod}</lastmod>` : null,
            e.changefreq ? `    <changefreq>${e.changefreq}</changefreq>` : null,
            e.priority ? `    <priority>${e.priority}</priority>` : null,
            `  </url>`,
          ]
            .filter(Boolean)
            .join("\n"),
        );

        const xml = [
          `<?xml version="1.0" encoding="UTF-8"?>`,
          `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`,
          ...urls,
          `</urlset>`,
        ].join("\n");

        return new Response(xml, {
          headers: {
            "Content-Type": "application/xml",
            "Cache-Control": "public, max-age=3600",
          },
        });
      },
    },
  },
});
