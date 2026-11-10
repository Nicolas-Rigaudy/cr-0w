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
