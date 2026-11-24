import type { Snapshot } from '@nest/data/useData'

export const label = 'Garden'

export default function Garden({ snap }: { snap: Snapshot }) {
  return (
    <div className="screen">
      <h1>Garden</h1>
      <p>{snap.inbox.length} inbox items</p>
    </div>
  )
}
