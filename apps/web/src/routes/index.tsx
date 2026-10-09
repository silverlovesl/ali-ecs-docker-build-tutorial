import { useQuery } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'

import { getHealth } from '@/lib/api'

export const Route = createFileRoute('/')({
  component: HomePage,
})

const stack = [
  { name: 'Turborepo', detail: 'task graph, caching, remote-ready' },
  { name: 'Bun', detail: 'workspaces + package manager' },
  { name: 'React 19 + TypeScript', detail: 'strict, bundler resolution' },
  { name: 'Tailwind CSS v4', detail: 'CSS-first config via @theme' },
  { name: 'TanStack Router + Query', detail: 'type-safe routing & data' },
  { name: 'Go + Gin', detail: 'graceful shutdown, air hot reload' },
]

function HomePage() {
  const health = useQuery({
    queryKey: ['health'],
    queryFn: ({ signal }) => getHealth(signal),
    retry: 1,
  })

  return (
    <div className="space-y-10">
      <section className="space-y-4">
        <p className="text-xs font-medium uppercase tracking-[0.2em] text-indigo-400">
          apps/web · apps/api
        </p>
        <h1 className="text-balance text-4xl font-semibold tracking-tight sm:text-5xl">
          Full-stack monorepo, ready to build on.
        </h1>
        <p className="max-w-2xl text-pretty text-slate-400">
          A Bun workspace driven by Turborepo. The React app talks to a Gin API through a Vite dev
          proxy, and formatting plus linting run on every commit.
        </p>
      </section>

      <section className="rounded-xl border border-slate-800 bg-slate-900/60 p-6">
        <div className="flex flex-wrap items-center justify-between gap-3">
          <h2 className="text-sm font-semibold uppercase tracking-wider text-slate-300">
            API status
          </h2>
          <span className="font-mono text-xs text-slate-500">GET /api/health</span>
        </div>

        <div className="mt-4 font-mono text-sm">
          {health.isPending && <p className="text-slate-400">Checking…</p>}
          {health.isError && (
            <p className="text-amber-400">
              API unreachable — start it with <code>bun run dev --filter=@repo/api</code>
            </p>
          )}
          {health.data && (
            <pre className="overflow-x-auto rounded-lg bg-slate-950 p-4 text-emerald-300">
              {JSON.stringify(health.data, null, 2)}
            </pre>
          )}
        </div>
      </section>

      <section className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {stack.map((item) => (
          <div
            key={item.name}
            className="rounded-xl border border-slate-800 bg-slate-900/40 p-5 transition-colors hover:border-slate-700"
          >
            <h3 className="text-sm font-semibold text-slate-100">{item.name}</h3>
            <p className="mt-1 text-sm text-slate-400">{item.detail}</p>
          </div>
        ))}
      </section>
    </div>
  )
}
