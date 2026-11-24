import type { Snapshot } from '@nest/data/useData'

export const page = 'home'

export default function Beds({ snap }: { snap: Snapshot }) {
  return <p className="garden-beds">Beds: {snap.missions.length}</p>
}
