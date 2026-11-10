import { Link } from '../router'

export function Placeholder({ name, note }: { name: string; note: string }) {
  return (
    <div className="screen">
      <section className="pn card" aria-labelledby="ph" style={{ maxWidth: 520 }}>
        <div className="ph"><h1 id="ph" style={{ font: 'inherit' }}>{name}</h1></div>
        <div className="dim">{note}</div>
        <div><Link href="/" className="btn">Back to Home</Link></div>
      </section>
    </div>
  )
}
