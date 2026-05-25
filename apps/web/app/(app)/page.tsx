'use client'

import { useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { RECIPES, TYPES, HERO_RECIPE_ID, totalDuration, sessionCount, tagVariant, fmt } from '@/lib/mock-data'
import { IconPlay, IconPlus, IconArrow, IconCheck, IconEdit, IconShare, IconGroup, IconBookmark } from '@/components/icons'
import { CozyFigureMug, CozyBranch } from '@/components/cozy-decorations'

function VisibilityTag({ visibility, isDefault }: { visibility: string; isDefault: boolean }) {
  if (isDefault) return <span className="tag tag--default">DEFAULT</span>
  if (visibility === 'public') return <span className="tag tag--public">PUBLIK</span>
  if (visibility === 'group') return <span className="tag tag--group">GRUP</span>
  return <span className="tag tag--private">PRIBADI</span>
}

function RecipeCard({ recipe, onClick }: { recipe: (typeof RECIPES)[0]; onClick: () => void }) {
  const total = totalDuration(recipe)
  const sessions = sessionCount(recipe)
  return (
    <button
      className="recipe-card"
      onClick={onClick}
      style={{ opacity: recipe.isArchived ? 0.6 : 1 }}
    >
      <div className="row gap-2 wrap" style={{ marginBottom: 2 }}>
        {recipe.tags.slice(0, 2).map((t) => (
          <span key={t} className={`tag tag--${tagVariant(t)}`}>{t}</span>
        ))}
        <span style={{ flex: 1 }} />
        {recipe.isArchived
          ? <span className="tag" style={{ background: 'var(--grid-paper)', color: 'var(--muted)', border: '1px solid var(--hairline)' }}>ARSIP</span>
          : <VisibilityTag visibility={recipe.visibility} isDefault={recipe.isDefault} />
        }
      </div>
      <div className="recipe-card__title">{recipe.name}</div>
      <div className="recipe-card__meta">
        <span>{sessions} sesi</span>
        <span className="dot" />
        <span className="t-mono-num">{Math.floor(total / 60)}m {total % 60}s</span>
        <span className="dot" />
        <span>{recipe.params.ratio}</span>
      </div>
      <div className="row between" style={{ marginTop: 'auto', paddingTop: 12 }}>
        <span className="t-caption">{recipe.subtype}</span>
        <span className="row gap-1" style={{ color: 'var(--deep-ink)', fontWeight: 700, fontSize: 13 }}>
          Lihat <IconArrow size={14} />
        </span>
      </div>
    </button>
  )
}

function HeroRecipeCard({ recipe, onStart, onOpen }: { recipe: (typeof RECIPES)[0]; onStart: () => void; onOpen: () => void }) {
  const total = totalDuration(recipe)
  const sessions = sessionCount(recipe)
  return (
    <div className="card card--dark card--hero" style={{ position: 'relative', overflow: 'hidden', padding: 40 }}>
      <CozyFigureMug />
      <div aria-hidden style={{
        position: 'absolute', right: -60, top: -60, width: 320, height: 320, borderRadius: '50%',
        background: 'radial-gradient(circle, rgba(61,43,255,.45), transparent 70%)',
      }} />
      <div className="row gap-2 wrap" style={{ marginBottom: 18, position: 'relative' }}>
        <span className="tag tag--espresso" style={{ fontSize: 13, padding: '5px 12px' }}>RESEP MINGGU INI</span>
        {recipe.tags.map((t) => (
          <span key={t} className={`tag tag--${tagVariant(t) === 'default' ? 'private' : tagVariant(t)}`} style={{ fontSize: 13, padding: '5px 12px' }}>{t}</span>
        ))}
      </div>
      <div style={{ position: 'relative', display: 'grid', gridTemplateColumns: '1.4fr 1fr', gap: 48, alignItems: 'end' }}>
        <div>
          <div className="t-display" style={{ marginBottom: 16, maxWidth: 540 }}>{recipe.name}</div>
          <div className="t-body muted-dark" style={{ maxWidth: 520, marginBottom: 28 }}>{recipe.description}</div>
          <div className="row gap-2">
            <button className="btn btn--primary btn--lg" onClick={onStart}>
              <IconPlay size={16} /> Mulai Brewing
            </button>
            <button className="btn btn--secondary-dark btn--lg" onClick={onOpen}>Lihat resep</button>
          </div>
        </div>
        <div className="col gap-3" style={{ alignItems: 'stretch' }}>
          {[
            { label: 'Total brew', val: `${Math.floor(total / 60)}:${String(total % 60).padStart(2, '0')}`, mono: true },
            { label: 'Sesi', val: String(sessions), mono: true },
            { label: 'Ratio', val: recipe.params.ratio, mono: false },
            { label: 'Suhu', val: recipe.params.temp, mono: true },
          ].map(({ label, val, mono }, i, arr) => (
            <div key={label} className="row between" style={{ paddingBottom: 14, borderBottom: i < arr.length - 1 ? '1px solid rgba(255,255,255,.10)' : undefined }}>
              <span className="t-label muted-dark">{label}</span>
              <span className={`t-h2${mono ? ' t-mono-num' : ''}`}>{val}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

const ACTIVITY = [
  { time: '08:14', action: 'Selesai brewing', what: 'V60 Light Roast 15g / 250ml', icon: <IconCheck size={16} /> },
  { time: 'Kemarin', action: 'Edit resep', what: 'Flair Espresso 18g / 36g', icon: <IconEdit size={16} /> },
  { time: 'Senin', action: 'Disetujui di grup', what: 'Kopi Susu Aren ke Komunitas Senayan', icon: <IconGroup size={16} /> },
  { time: 'Minggu', action: 'Duplikat', what: 'V60 Dark Roast 18g / 270ml dari Sistem YAVA', icon: <IconShare size={16} /> },
]

export default function Dashboard() {
  const router = useRouter()
  const [activeType, setActiveType] = useState('Semua')
  const [showArchived, setShowArchived] = useState(false)
  const hero = RECIPES.find((r) => r.id === HERO_RECIPE_ID)!

  const activeRecipes = RECIPES.filter((r) =>
    !r.isArchived && (activeType === 'Semua' || r.type === activeType)
  )
  const archivedRecipes = RECIPES.filter((r) =>
    r.isArchived && (activeType === 'Semua' || r.type === activeType)
  )

  return (
    <main className="container col gap-4" style={{ paddingTop: 16, gap: 48 }}>
      {/* Greeting */}
      <section className="col gap-3">
        <div className="row between" style={{ alignItems: 'flex-end' }}>
          <div>
            <div className="t-label" style={{ color: 'var(--muted)', marginBottom: 8 }}>Selasa · 23 Mei 2026</div>
            <div className="t-display" style={{ fontSize: 48, letterSpacing: '-.025em' }}>
              Pagi, Nadira.<span style={{ color: 'var(--coral-red)' }}>.</span>
            </div>
            <div className="t-body muted" style={{ marginTop: 6 }}>{activeRecipes.length} resep aktif · 12 sesi brewing minggu ini.</div>
          </div>
          <div className="row gap-2">
            <button className="btn btn--light-secondary"><IconBookmark size={16} /> Favorit</button>
            <Link href="/recipes/new" className="btn btn--light-primary"><IconPlus size={16} /> Buat resep</Link>
          </div>
        </div>
        <HeroRecipeCard
          recipe={hero}
          onStart={() => router.push(`/recipes/${hero.id}/brew`)}
          onOpen={() => router.push(`/recipes/${hero.id}`)}
        />
      </section>

      {/* Library */}
      <section className="col gap-3">
        <div className="row between" style={{ alignItems: 'flex-end' }}>
          <div>
            <div className="t-label" style={{ color: 'var(--muted)', marginBottom: 6 }}>Koleksi</div>
            <div className="t-h1">Resep saya</div>
          </div>
          <div className="tabs">
            {TYPES.map((t) => (
              <button
                key={t}
                className={`tab${t === activeType ? ' tab--active' : ''}`}
                onClick={() => setActiveType(t)}
              >
                {t}
              </button>
            ))}
          </div>
        </div>

        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 16 }}>
          {activeRecipes.map((r) => (
            <RecipeCard key={r.id} recipe={r} onClick={() => router.push(`/recipes/${r.id}`)} />
          ))}
          <Link
            href="/recipes/new"
            className="recipe-card"
            style={{
              background: 'transparent',
              border: '1.5px dashed var(--deep-ink)',
              alignItems: 'center',
              justifyContent: 'center',
              minHeight: 200,
              color: 'var(--deep-ink)',
              fontWeight: 700,
              textDecoration: 'none',
              display: 'flex',
              flexDirection: 'column',
              gap: 4,
            }}
          >
            <IconPlus size={28} />
            <span style={{ marginTop: 4 }}>Buat resep baru</span>
            <span className="t-caption" style={{ color: 'var(--muted)' }}>Tambahkan sesi timer atau notes</span>
          </Link>
        </div>

        {/* Arsip section */}
        {archivedRecipes.length > 0 && (
          <div style={{ marginTop: 8 }}>
            <button
              onClick={() => setShowArchived((v) => !v)}
              style={{ background: 'transparent', border: 0, cursor: 'pointer', display: 'inline-flex', alignItems: 'center', gap: 8, padding: 0, marginBottom: showArchived ? 16 : 0 }}
            >
              <span className="t-label" style={{ color: 'var(--muted)' }}>
                Arsip ({archivedRecipes.length})
              </span>
              <span style={{ fontSize: 11, color: 'var(--muted)', fontWeight: 700 }}>
                {showArchived ? '▲' : '▼'}
              </span>
            </button>
            {showArchived && (
              <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 16 }}>
                {archivedRecipes.map((r) => (
                  <RecipeCard key={r.id} recipe={r} onClick={() => router.push(`/recipes/${r.id}`)} />
                ))}
              </div>
            )}
          </div>
        )}
      </section>

      {/* Activity + Group */}
      <section style={{ display: 'grid', gridTemplateColumns: '1.4fr 1fr', gap: 16, paddingBottom: 64 }}>
        <div className="card" style={{ padding: 24 }}>
          <div className="row between" style={{ marginBottom: 16 }}>
            <div className="t-h2">Aktivitas terakhir</div>
            <a href="#" className="t-small" style={{ color: 'var(--deep-ink)', fontWeight: 700 }}>Lihat semua</a>
          </div>
          <div className="col">
            {ACTIVITY.map((row, i) => (
              <div key={i} className="row gap-3" style={{ padding: '14px 0', borderTop: i === 0 ? '0' : '1px solid var(--hairline)' }}>
                <div style={{
                  width: 32, height: 32, borderRadius: '50%',
                  display: 'inline-flex', alignItems: 'center', justifyContent: 'center',
                  background: 'var(--lavender-fog)', color: 'var(--deep-ink)',
                }}>{row.icon}</div>
                <div className="grow">
                  <div style={{ fontWeight: 700, fontSize: 14 }}>{row.action}</div>
                  <div className="t-small muted">{row.what}</div>
                </div>
                <div className="t-small muted t-mono-num">{row.time}</div>
              </div>
            ))}
          </div>
        </div>

        <div className="card card--electric card--hero" style={{ padding: 28, position: 'relative', overflow: 'hidden' }}>
          <CozyBranch />
          <span className="tag" style={{ background: 'rgba(255,255,255,.12)', color: 'var(--lilac)', border: '0' }}>GRUP AKTIF</span>
          <div className="t-h1" style={{ color: '#fff', marginTop: 16, lineHeight: 1.05 }}>Komunitas<br />Senayan</div>
          <div className="t-small" style={{ color: 'var(--lilac)', marginTop: 8, marginBottom: 24 }}>
            28 anggota · 42 resep aktif · 2 menunggu approval
          </div>
          <div className="row gap-1" style={{ marginBottom: 20 }}>
            {['RW', 'AD', 'MS', '+25'].map((n, i) => (
              <span key={i} className="avatar avatar--sm" style={{
                background: i === 3 ? 'rgba(255,255,255,.18)' : 'var(--lilac)',
                color: i === 3 ? 'var(--lilac)' : 'var(--electric)',
                border: '2px solid var(--electric)',
                marginLeft: i === 0 ? 0 : -8,
              }}>{n}</span>
            ))}
          </div>
          <button className="btn" style={{ background: 'var(--lilac)', color: 'var(--electric)' }}>
            Buka grup <IconArrow size={14} />
          </button>
        </div>
      </section>
    </main>
  )
}
