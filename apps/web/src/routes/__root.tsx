import type { QueryClient } from '@tanstack/react-query'
import { Link, Outlet, createRootRouteWithContext } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'

export interface RouterContext {
  queryClient: QueryClient
}

export const Route = createRootRouteWithContext<RouterContext>()({
  component: RootLayout,
})

const navigation = [
  { to: '/', label: 'Home' },
  { to: '/about', label: 'About' },
] as const

function RootLayout() {
  return (
    <div className="flex min-h-full flex-col bg-slate-950 text-slate-100">
      <header className="border-b border-slate-800/80 bg-slate-950/80 backdrop-blur">
        <nav className="mx-auto flex h-16 w-full max-w-5xl items-center justify-between px-6">
          <Link to="/" className="flex items-center gap-2 text-sm font-semibold tracking-tight">
            <span className="grid size-7 place-items-center rounded-md bg-indigo-500/15 text-indigo-300">
              ▲
            </span>
            monorepo
          </Link>
          <ul className="flex items-center gap-1 text-sm">
            {navigation.map((item) => (
              <li key={item.to}>
                <Link
                  to={item.to}
                  activeProps={{ className: 'bg-slate-800 text-white' }}
                  inactiveProps={{ className: 'text-slate-400 hover:text-slate-100' }}
                  className="rounded-md px-3 py-1.5 transition-colors"
                >
                  {item.label}
                </Link>
              </li>
            ))}
          </ul>
        </nav>
      </header>

      <main className="mx-auto w-full max-w-5xl flex-1 px-6 py-12">
        <Outlet />
      </main>

      <footer className="border-t border-slate-800/80 py-6 text-center text-xs text-slate-500">
        Turborepo · Bun · React · TanStack · Go · Gin
      </footer>

      <TanStackRouterDevtools position="bottom-right" />
    </div>
  )
}
