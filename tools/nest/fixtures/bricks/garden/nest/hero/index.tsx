import type { HeroProps } from '@nest/modules/types'

export default function GardenHero({ snap, say }: HeroProps) {
  return <div className="garden-hero" role="img" aria-label="Garden hero">{say ?? 'Garden'} · {snap.missions.length} missions</div>
}
