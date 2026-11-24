// What a brick's nest/ files export (item_267, docs/writing-a-brick.md).
import type { ComponentType, ReactNode } from 'react'
import type { Snapshot } from '../data/feed'

// Props every hero gets, built-in or from a brick (nest/hero/index.tsx, default export).
export interface HeroProps {
  snap: Snapshot
  // The butler line under the hero on phone Home.
  say?: string
}

// nest/pages/<name>.tsx: the page at /<name>.
export interface BrickPageModule {
  default: ComponentType<{ snap: Snapshot }>
  label?: string
  // SVG children drawn in a 22x22 box with a stroke, like the engine tab icons.
  icon?: ReactNode
}

// nest/sections/<name>.tsx: a panel after the Home or Desk panels.
export interface BrickSectionModule {
  default: ComponentType<{ snap: Snapshot }>
  page: 'home' | 'desk'
}
