import { createFileRoute, Link, Outlet, useRouterState } from "@tanstack/react-router";
import { ChevronRight } from "lucide-react";
import { Footer } from "@/components/Footer";
import { DocsSearch } from "@/components/docs/DocsSearch";
import { DOC_PAGES, DOC_SECTIONS } from "@/components/docs/docsNav";

export const Route = createFileRoute("/docs")({
  component: DocsLayout,
});

function SidebarNav() {
  const pathname = useRouterState({ select: (s) => s.location.pathname });
  const isActive = (path: string) =>
    path === "/docs" ? pathname === "/docs" : pathname === path;

  return (
    <nav className="space-y-8">
      {DOC_SECTIONS.map((section) => (
        <div key={section}>
          <p className="mb-3 font-mono text-[11px] font-semibold uppercase tracking-wider text-muted-foreground">
            {section}
          </p>
          <ul className="space-y-1">
            {DOC_PAGES.filter((p) => p.section === section).map((p) => {
              const active = isActive(p.path);
              return (
                <li key={p.path}>
                  <Link
                    to={p.path}
                    className={`flex items-center gap-1.5 rounded-md px-3 py-2 font-mono text-[13px] transition-all ${
                      active
                        ? "bg-primary/10 text-primary text-glow"
                        : "text-muted-foreground hover:bg-muted hover:text-foreground"
                    }`}
                  >
                    <ChevronRight
                      size={12}
                      className={`transition-opacity ${active ? "opacity-100 text-primary" : "opacity-0"}`}
                    />
                    {p.title}
                  </Link>
                </li>
              );
            })}
          </ul>
        </div>
      ))}
    </nav>
  );
}

function DocsLayout() {
  return (
    <div className="min-h-screen bg-background font-sans">
      {/* Top bar */}
      <header className="fixed inset-x-0 top-0 z-50 glass border-x-0 border-t-0">
        <div className="mx-auto flex h-14 max-w-7xl items-center gap-6 px-4 sm:px-6">
          <Link
            to="/"
            className="shrink-0 font-mono text-sm font-bold text-primary text-glow transition-opacity hover:opacity-80"
          >
            [Mr. Browser]
          </Link>
          <span className="hidden font-mono text-xs text-muted-foreground sm:block">/docs</span>
          <div className="ml-auto flex w-full max-w-md items-center justify-end">
            <DocsSearch />
          </div>
        </div>
      </header>

      <div className="mx-auto flex max-w-7xl gap-10 px-4 pt-14 sm:px-6">
        {/* Sidebar */}
        <aside className="sticky top-14 hidden h-[calc(100vh-3.5rem)] w-60 shrink-0 overflow-y-auto border-r border-border py-10 pr-4 lg:block">
          <SidebarNav />
        </aside>

        {/* Content */}
        <main className="min-w-0 flex-1 py-12">
          {/* Mobile nav */}
          <div className="mb-8 flex gap-2 overflow-x-auto pb-2 lg:hidden">
            {DOC_PAGES.map((p) => (
              <Link
                key={p.path}
                to={p.path}
                activeOptions={{ exact: true }}
                activeProps={{
                  className: "border-primary/60 bg-primary/10 text-primary",
                }}
                className="shrink-0 rounded-full border border-border px-4 py-1.5 font-mono text-xs text-muted-foreground"
              >
                {p.title}
              </Link>
            ))}
          </div>
          <Outlet />
        </main>
      </div>

      <Footer />
    </div>
  );
}
