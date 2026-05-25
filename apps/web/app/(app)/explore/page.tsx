'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { RECIPES, TYPES, totalDuration, sessionCount, tagVariant } from '@/lib/mock-data'
import { IconSearch, IconArrow, IconGlobe } from '@/components/icons'

const PUBLIC_RECIPES = RECIPES.filter((r) => r.visibility === 'public' && !r.isArchived)

export default function ExplorePage() {
  const router = useRouter()
  const [activeType, setActiveType] = useState('Semua')
  const [query, setQuery] = useState('')

  const filtered = PUBLIC_RECIPES.filter((r) => {
    const matchType = activeType === 'Semua' || r.type === activeType
    const q = query.toLowerCase()
    const matchQuery = !q || r.name.toLowerCase().includes(q) || r.type.toLowerCase().includes(q) || r.subtype.toLowerCase().includes(q)
    return matchType && matchQuery
  })

  return (
    <main className="container col gap-4" style={{ paddingTop: 16, gap: 48, paddingBottom: 64 }}>
      {/* Header */}
      <section>
        <div className="row between" style={{ alignItems: 'flex-end', marginBottom: 28 }}>
          <div>
            <div className="t-label" style={{ color: 'var(--muted)', marginBottom: 8, display: 'flex', alignItems: 'center', gap: 6 }}>
              <IconGlobe size={12} /> Komunitas
            </div>
            <div className="t-display" style={{ fontSize: 48, letterSpacing: '-.025em' }}>
              Explore resep<span style={{ color: 'var(--coral-red)' }}>.</span>
            </div>
            <div className="t-body muted" style={{ marginTop: 6 }}>
              {PUBLIC_RECIPES.length} resep publik dari komunitas YAVA.
            </div>
          </div>
        </div>

        {/* Search + filter bar */}
        <div className="row gap-3" style={{ flexWrap: 'wrap', alignItems: 'center' }}>
          <div style={{ position: 'relative', flex: '0 0 320px' }}>
            <IconSearch size={16} style={{ position: 'absolute', left: 14, top: '50%', transform: 'translateY(-50%)', color: 'var(--muted)' }} />
            <input
              className="input input--search"
              placeholder="Cari nama, jenis, atau subtype…"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              style={{ width: '100%', paddingLeft: 42 }}
            />
          </div>
          <div className="tabs">
            {TYPES.map((t) => (
              <button
                key={t}
                className={`tab${t === activeType ? ' tab--active' : ''}`}
                onClick={() => setActiveType(t)}
              >{t}</button>
            ))}
          </div>
        </div>
      </section>

      {/* Results */}
      <section>
        {filtered.length === 0 ? (
          <div className="card" style={{ padding: 48, textAlign: 'center' }}>
            <div className="t-h2" style={{ marginBottom: 8 }}>Tidak ada resep</div>
            <div className="t-body muted">Coba kata kunci atau filter yang berbeda.</div>
          </div>
        ) : (
          <>
            <div className="t-label muted" style={{ marginBottom: 16 }}>
              {filtered.length} resep ditemukan
            </div>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 16 }}>
              {filtered.map((r) => {
                const total = totalDuration(r)
                const sessions = sessionCount(r)
                return (
                  <button
                    key={r.id}
                    className="recipe-card"
                    onClick={() => router.push(`/recipes/${r.id}`)}
                  >
                    <div className="row gap-2 wrap" style={{ marginBottom: 2 }}>
                      {r.tags.slice(0, 2).map((t) => (
                        <span key={t} className={`tag tag--${tagVariant(t)}`}>{t}</span>
                      ))}
                      {r.isDefault && (
                        <span className="tag tag--default">DEFAULT</span>
                      )}
                    </div>

                    <div className="recipe-card__title">{r.name}</div>
                    <div className="t-small muted" style={{ lineHeight: 1.45, WebkitLineClamp: 2, display: '-webkit-box', WebkitBoxOrient: 'vertical', overflow: 'hidden' }}>
                      {r.description}
                    </div>

                    <div className="recipe-card__meta" style={{ marginTop: 4 }}>
                      <span>{sessions} sesi</span>
                      <span className="dot" />
                      <span className="t-mono-num">{Math.floor(total / 60)}m {total % 60}s</span>
                      <span className="dot" />
                      <span>{r.params.ratio}</span>
                    </div>

                    <div className="row between" style={{ marginTop: 'auto', paddingTop: 12 }}>
                      <div className="row gap-2" style={{ alignItems: 'center' }}>
                        <span className="avatar avatar--sm" style={{ background: 'var(--coral-red)' }}>{r.author.initials}</span>
                        <span className="t-caption">{r.author.name}</span>
                        <span className="t-caption">·</span>
                        <span className="t-caption">{r.saves} simpan</span>
                      </div>
                      <span className="row gap-1" style={{ color: 'var(--deep-ink)', fontWeight: 700, fontSize: 13 }}>
                        Lihat <IconArrow size={14} />
                      </span>
                    </div>
                  </button>
                )
              })}
            </div>
          </>
        )}
      </section>
    </main>
  )
}
