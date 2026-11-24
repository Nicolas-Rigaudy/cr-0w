declare module 'virtual:core-records' {
  export const missions: unknown[]
  export const inbox: unknown[]
}

declare module 'virtual:core-api' {
  const files: Record<string, string>
  export default files
}

interface ImportMetaEnv {
  // 'api' makes npm run dev read the core API instead of fixtures.
  readonly VITE_NEST_SOURCE?: string
}

// Bricks' Nest files, compiled in at build time (item_267, nest-modules.ts).
declare module 'virtual:nest-modules' {
  import type { ComponentType } from 'react'
  import type { BrickPageModule, BrickSectionModule, HeroProps } from './modules/types'
  export const config: { theme: string; hero: string; pages: string[] }
  export const themes: string[]
  export const heroes: Record<string, ComponentType<HeroProps>>
  export const pages: Array<{ key: string; brick: string; mod: BrickPageModule }>
  export const sections: Array<{ id: string; brick: string; mod: BrickSectionModule }>
}

declare module 'virtual:nest-themes' {}
