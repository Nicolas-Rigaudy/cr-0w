import type { UiState } from '../data/state'
import { LABEL } from '../data/state'

const CLASS: Record<UiState, string> = { working: 'c-w', needs: 'c-n', failed: 'c-f', sleeping: 'c-s', queued: 'c-q', done: 'c-d' }
export const stateClass = (s: UiState) => CLASS[s]

// Shape carries the state; colour only repeats it. Pair it with the status word.
export function Glyph({ state, size = 14 }: { state: UiState; size?: number }) {
  const common = { width: size, height: size, 'aria-hidden': true, className: `glyph ${CLASS[state]}` } as const
  switch (state) {
    case 'working':
      return (
        <svg {...common} viewBox="0 0 22 22" fill="none">
          <circle cx="11" cy="11" r="4" fill="currentColor" />
          <path d="M11 2a9 9 0 0 1 9 9" stroke="currentColor" strokeWidth="2" />
        </svg>
      )
    case 'needs':
      return (
        <svg {...common} viewBox="0 0 16 16">
          <path d="M8 1 15 8 8 15 1 8Z" fill="currentColor" />
        </svg>
      )
    case 'failed':
      return (
        <svg {...common} viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.6">
          <path d="M8 2a6 6 0 0 1 6 6M8 14a6 6 0 0 1-6-6" />
          <path d="M5.5 5.5l5 5M10.5 5.5l-5 5" />
        </svg>
      )
    case 'sleeping':
      return (
        <svg {...common} viewBox="0 0 16 16">
          <path d="M10.5 2.5a6 6 0 1 0 3 9 5 5 0 0 1-3-9Z" fill="currentColor" />
        </svg>
      )
    case 'queued':
      return (
        <svg {...common} viewBox="0 0 18 18" fill="none">
          <circle cx="9" cy="9" r="6" stroke="currentColor" strokeWidth="1.6" strokeDasharray="2 3" />
        </svg>
      )
    case 'done':
      return (
        <svg {...common} viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round">
          <circle cx="8" cy="8" r="6.2" strokeWidth="1.4" />
          <path d="M5 8.2l2 2 4-4" />
        </svg>
      )
  }
}

export function StateWord({ state }: { state: UiState }) {
  return <span className={`tag ${CLASS[state]}`} style={{ borderColor: 'currentColor' }}>{LABEL[state]}</span>
}
