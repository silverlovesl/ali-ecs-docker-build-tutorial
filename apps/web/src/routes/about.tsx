import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/about')({
  component: AboutPage,
})

const commands = [
  { command: 'bun install', detail: 'install every workspace' },
  { command: 'bun run dev', detail: 'run web + api together' },
  { command: 'bun run build', detail: 'build all packages' },
  { command: 'bun run lint', detail: 'oxlint across the repo' },
  { command: 'bun run format', detail: 'oxfmt + gofmt' },
  { command: 'bun run typecheck', detail: 'tsc per workspace' },
]

function AboutPage() {
  return (
    <div className="space-y-8">
      <header className="space-y-3">
        <h1 className="text-3xl font-semibold tracking-tight">About this setup</h1>
        <p className="max-w-2xl text-slate-400">
          Shared configuration lives in <code className="text-slate-200">packages/</code>, so every
          workspace lints, formats, and type-checks the same way.
        </p>
      </header>

      <div className="overflow-hidden rounded-xl border border-slate-800">
        <table className="w-full text-left text-sm">
          <thead className="bg-slate-900/80 text-xs uppercase tracking-wider text-slate-400">
            <tr>
              <th className="px-4 py-3 font-medium">Command</th>
              <th className="px-4 py-3 font-medium">What it does</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-800 bg-slate-900/30">
            {commands.map((row) => (
              <tr key={row.command}>
                <td className="px-4 py-3 font-mono text-indigo-300">{row.command}</td>
                <td className="px-4 py-3 text-slate-400">{row.detail}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
