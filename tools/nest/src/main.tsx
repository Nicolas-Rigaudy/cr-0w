import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { App } from './App'
import { fixtureSource } from './data/fixtures'
import { DataContext } from './data/useData'
import { Router } from './router'
import './theme.css'

// The adapter is picked here and nowhere else (AC3).
const source = fixtureSource()

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <DataContext.Provider value={source}>
      <Router>
        <App />
      </Router>
    </DataContext.Provider>
  </StrictMode>,
)

if ('serviceWorker' in navigator && import.meta.env.PROD) {
  addEventListener('load', () => navigator.serviceWorker.register('/sw.js'))
}
