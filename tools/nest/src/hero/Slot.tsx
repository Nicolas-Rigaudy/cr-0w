// The Home hero slot (item_267): crow.yaml modules.nest.hero picks holotable (built in), none, or a brick hero.
import type { ComponentProps } from 'react'
import { config, heroes } from 'virtual:nest-modules'
import type { Snapshot } from '../data/useData'
import { Hero } from './Hero'

export const HERO = config.hero

export function HeroSlot({ snap, ...props }: ComponentProps<typeof Hero> & { snap: Snapshot }) {
  if (HERO === 'holotable') return <Hero {...props} />
  const Brick = heroes[HERO]
  return Brick ? <Brick snap={snap} say={props.say} /> : null
}
