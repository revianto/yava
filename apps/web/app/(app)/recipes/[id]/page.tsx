'use client'

import { use, useState, useRef, useEffect } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { notFound } from 'next/navigation'
import { RECIPES, GROUPS, DISCUSSIONS, totalDuration, sessionCount, tagVariant } from '@/lib/mock-data'
import {
  IconPlay, IconArrowLeft, IconStar, IconShare, IconEdit, IconKebab,
  IconArchive, IconRestore, IconCopy, IconGroup,
} from '@/components/icons'
import { CozyMugSteam, CozyPlants } from '@/components/cozy-decorations'
import { DiscussionThread } from '@/components/discussion-thread'

function VisibilityTag({ visibility, isDefault }: { visibility: string; isDefault: boolean }) {
  if (isDefault) return <span className="tag tag--default">DEFAULT</span>
  if (visibility === 'public') return <span className="tag tag--public">PUBLIK</span>
  if (visibility === 'group') return <span className="tag tag--group">GRUP</span>
  return <span className="tag tag--private">PRIBADI</span>
}

function StepRow({ index, step }: { index: number; step: (typeof RECIPES)[0]['timeline'][0] }) {
  if (step.kind === 'note') {
    return (
      <div className="step">
        <div className="step__num step__num--note">N</div>
        <div>
          <div className="step__name muted">Catatan</div>
          <div className="step__note">{step.content}</div>
        </div>
        <div className="step__time muted">—</div>
      </div>
    )
  }
  return (
    <div className="step">
      <div className="step__num">{index}</div>
      <div>
        <div className="step__name">{step.name}</div>
        <div className="step__note">{step.note}</div>
      </div>
      <div className="step__time">{step.duration}<small>s</small></div>
    </div>
  )
}

export default function RecipeDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params)
  const router = useRouter()
  const recipe = RECIPES.find((r) => r.id === id)
  if (!recipe) notFound()

  const total = totalDuration(recipe)
  const sessions = sessionCount(recipe)
  let sessionIdx = 0

  const [isArchived, setIsArchived] = useState(recipe.isArchived ?? false)
  const [menuOpen, setMenuOpen] = useState(false)
  const [duplicated, setDuplicated] = useState(false)
  const [submitGroupOpen, setSubmitGroupOpen] = useState(false)
  const [submittedGroup, setSubmittedGroup] = useState<string | null>(null)
  const menuRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!menuOpen) return
    const handler = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) setMenuOpen(false)
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [menuOpen])

  const handleDuplicate = () => {
    setMenuOpen(false)
    setDuplicated(true)
    setTimeout(() => setDuplicated(false), 2500)
  }

  return (
    <main className="container">
      {/* Breadcrumb */}
      <button
        onClick={() => router.back()}
        style={{
          background: 'transparent', border: 0, cursor: 'pointer',
          display: 'inline-flex', alignItems: 'center', gap: 6,
          color: 'var(--muted)', fontWeight: 500, fontSize: 13, padding: 0, marginBottom: 24,
        }}
      >
        <IconArrowLeft size={14} /> Kembali ke {recipe.type} · {recipe.subtype}
      </button>

      {/* Toast notifications */}
      {(duplicated || submittedGroup) && (
        <div style={{
          position: 'fixed', bottom: 32, left: '50%', transform: 'translateX(-50%)',
          background: 'var(--deep-ink)', color: '#FFF8E7', padding: '12px 24px',
          borderRadius: 99, fontWeight: 700, fontSize: 14, zIndex: 100,
          boxShadow: '0 8px 32px rgba(0,0,0,.24)',
          animation: 'fadeIn 200ms ease-out',
        }}>
          {duplicated
            ? 'Resep berhasil diduplikat ke koleksi kamu.'
            : `Resep dikirim ke ${submittedGroup} untuk ditinjau admin.`}
        </div>
      )}

      {/* Title block */}
      <div className="row between" style={{ alignItems: 'flex-start', gap: 32, marginBottom: 32, position: 'relative' }}>
        <CozyPlants />
        <div style={{ maxWidth: 720 }}>
          <div className="row gap-2 wrap" style={{ marginBottom: 14 }}>
            {recipe.tags.map((t) => (
              <span key={t} className={`tag tag--${tagVariant(t)}`}>{t}</span>
            ))}
            <VisibilityTag visibility={recipe.visibility} isDefault={recipe.isDefault} />
            {isArchived && (
              <span className="tag" style={{ background: 'var(--grid-paper)', color: 'var(--muted)', border: '1px solid var(--hairline)' }}>ARSIP</span>
            )}
          </div>
          <div className="t-display" style={{ fontSize: 52, marginBottom: 12, opacity: isArchived ? .55 : 1 }}>{recipe.name}</div>
          <div className="t-body muted" style={{ maxWidth: 600 }}>{recipe.description}</div>
          <div className="row gap-2" style={{ marginTop: 20, alignItems: 'center' }}>
            <span className="avatar avatar--sm" style={{ background: 'var(--coral-red)' }}>{recipe.author.initials}</span>
            <span className="t-small">oleh <strong>{recipe.author.name}</strong></span>
            <span className="t-small muted">·</span>
            <span className="t-small muted">{recipe.saves} simpan</span>
            <span className="t-small muted">·</span>
            <span className="t-small muted">terakhir brew {recipe.lastBrewed}</span>
          </div>
        </div>
        <div className="row gap-1">
          <button className="icon-btn" title="Favorit"><IconStar size={18} /></button>
          <button className="icon-btn" title="Bagikan"><IconShare size={18} /></button>
          {!recipe.isDefault && (
            <Link href={`/recipes/${id}/edit`} className="icon-btn" title="Edit" style={{ display: 'inline-flex', alignItems: 'center', justifyContent: 'center', textDecoration: 'none' }}>
              <IconEdit size={18} />
            </Link>
          )}
          {/* Kebab dropdown */}
          <div ref={menuRef} style={{ position: 'relative' }}>
            <button className="icon-btn" title="Lainnya" onClick={() => { setMenuOpen((v) => !v); setSubmitGroupOpen(false) }}>
              <IconKebab size={18} />
            </button>
            {menuOpen && (
              <div style={{
                position: 'absolute', top: 'calc(100% + 8px)', right: 0, zIndex: 50,
                background: '#FBF8F1', border: '1px solid var(--hairline)', borderRadius: 14,
                boxShadow: '0 8px 32px -8px rgba(45,82,74,.20)',
                minWidth: 220, overflow: 'hidden',
              }}>
                {submitGroupOpen ? (
                  <>
                    <div style={{ padding: '10px 16px 8px', display: 'flex', alignItems: 'center', gap: 8 }}>
                      <button
                        onClick={() => setSubmitGroupOpen(false)}
                        style={{ background: 'none', border: 0, cursor: 'pointer', padding: 0, color: 'var(--muted)', display: 'flex', alignItems: 'center' }}
                      >
                        <IconArrowLeft size={14} />
                      </button>
                      <span className="t-label" style={{ color: 'var(--muted)' }}>Pilih grup tujuan</span>
                    </div>
                    <div style={{ height: 1, background: 'var(--hairline)' }} />
                    {GROUPS.map((g) => (
                      <button
                        key={g.id}
                        onClick={() => {
                          setSubmittedGroup(g.name)
                          setMenuOpen(false)
                          setSubmitGroupOpen(false)
                          setTimeout(() => setSubmittedGroup(null), 3000)
                        }}
                        style={{ width: '100%', padding: '12px 16px', background: 'transparent', border: 0, cursor: 'pointer', display: 'flex', alignItems: 'center', gap: 10, textAlign: 'left' }}
                        onMouseEnter={(e) => (e.currentTarget.style.background = 'var(--lavender-fog)')}
                        onMouseLeave={(e) => (e.currentTarget.style.background = 'transparent')}
                      >
                        <span style={{
                          width: 32, height: 32, borderRadius: 8, background: 'var(--deep-ink)',
                          display: 'inline-flex', alignItems: 'center', justifyContent: 'center', flex: 'none',
                        }}>
                          <IconGroup size={14} style={{ color: 'var(--lilac)' }} />
                        </span>
                        <div>
                          <div style={{ fontSize: 13, fontWeight: 700, color: 'var(--deep-ink)' }}>{g.name}</div>
                          <div style={{ fontSize: 11, color: 'var(--muted)', marginTop: 1 }}>{g.members.length} anggota</div>
                        </div>
                      </button>
                    ))}
                  </>
                ) : (
                  <>
                    <button
                      onClick={handleDuplicate}
                      style={{ width: '100%', padding: '12px 16px', background: 'transparent', border: 0, cursor: 'pointer', display: 'flex', alignItems: 'center', gap: 10, fontSize: 14, fontWeight: 600, color: 'var(--deep-ink)', textAlign: 'left' }}
                      onMouseEnter={(e) => (e.currentTarget.style.background = 'var(--lavender-fog)')}
                      onMouseLeave={(e) => (e.currentTarget.style.background = 'transparent')}
                    >
                      <IconCopy size={16} /> Duplikat resep
                    </button>
                    <button
                      onClick={() => setSubmitGroupOpen(true)}
                      style={{ width: '100%', padding: '12px 16px', background: 'transparent', border: 0, cursor: 'pointer', display: 'flex', alignItems: 'center', gap: 10, fontSize: 14, fontWeight: 600, color: 'var(--deep-ink)', textAlign: 'left' }}
                      onMouseEnter={(e) => (e.currentTarget.style.background = 'var(--lavender-fog)')}
                      onMouseLeave={(e) => (e.currentTarget.style.background = 'transparent')}
                    >
                      <IconGroup size={16} /> Submit ke grup <span style={{ marginLeft: 'auto', color: 'var(--muted)', fontSize: 12 }}>›</span>
                    </button>
                    <div style={{ height: 1, background: 'var(--hairline)', margin: '0 12px' }} />
                    {isArchived ? (
                      <button
                        onClick={() => { setIsArchived(false); setMenuOpen(false) }}
                        style={{ width: '100%', padding: '12px 16px', background: 'transparent', border: 0, cursor: 'pointer', display: 'flex', alignItems: 'center', gap: 10, fontSize: 14, fontWeight: 600, color: 'var(--deep-ink)', textAlign: 'left' }}
                        onMouseEnter={(e) => (e.currentTarget.style.background = 'var(--lavender-fog)')}
                        onMouseLeave={(e) => (e.currentTarget.style.background = 'transparent')}
                      >
                        <IconRestore size={16} /> Kembalikan dari arsip
                      </button>
                    ) : (
                      <button
                        onClick={() => { setIsArchived(true); setMenuOpen(false) }}
                        style={{ width: '100%', padding: '12px 16px', background: 'transparent', border: 0, cursor: 'pointer', display: 'flex', alignItems: 'center', gap: 10, fontSize: 14, fontWeight: 600, color: 'var(--deep-ink)', textAlign: 'left' }}
                        onMouseEnter={(e) => (e.currentTarget.style.background = 'var(--lavender-fog)')}
                        onMouseLeave={(e) => (e.currentTarget.style.background = 'transparent')}
                      >
                        <IconArchive size={16} /> Arsipkan resep
                      </button>
                    )}
                  </>
                )}
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Params */}
      <div className="params">
        {[
          { label: 'Dose', val: recipe.params.dose },
          { label: 'Yield', val: recipe.params.yield },
          { label: 'Suhu', val: recipe.params.temp },
          { label: 'Grind', val: recipe.params.grind },
          { label: 'Ratio', val: recipe.params.ratio },
        ].map(({ label, val }) => (
          <div key={label} className="params__cell">
            <span className="lbl">{label}</span>
            <span className="val t-mono-num">{val}</span>
          </div>
        ))}
      </div>

      {/* Two-col: timeline + CTA */}
      <div style={{ display: 'grid', gridTemplateColumns: '1.6fr 1fr', gap: 32, marginTop: 32, paddingBottom: 64 }}>
        {/* Timeline */}
        <div>
          <div className="row between" style={{ marginBottom: 12, alignItems: 'flex-end' }}>
            <div>
              <div className="t-label" style={{ color: 'var(--muted)', marginBottom: 6 }}>Alur brewing</div>
              <div className="t-h2">{sessions} sesi · {Math.floor(total / 60)} menit {total % 60} detik total</div>
            </div>
            <button className="btn btn--secondary"><IconEdit size={14} /> Edit alur</button>
          </div>
          <div className="card" style={{ padding: '8px 24px' }}>
            {recipe.timeline.map((step, i) => {
              if (step.kind === 'session') sessionIdx++
              return <StepRow key={i} index={sessionIdx} step={step} />
            })}
          </div>

          {recipe.visibility === 'group' && (
            <DiscussionThread
              recipeId={id}
              initialDiscussions={DISCUSSIONS.filter((d) => d.recipeId === id)}
              isAdmin={false}
            />
          )}
        </div>

        {/* Start CTA */}
        <aside>
          <div className="card card--dark card--hero" style={{ padding: 28, position: 'sticky', top: 24, overflow: 'hidden' }}>
            <CozyMugSteam />
            <div className="t-label muted-dark" style={{ marginBottom: 12 }}>Siap brewing</div>
            <div className="t-h1" style={{ color: '#fff', marginBottom: 6 }}>
              {Math.floor(total / 60)}<span className="muted-dark" style={{ fontSize: 22 }}>m</span>{' '}
              {total % 60}<span className="muted-dark" style={{ fontSize: 22 }}>s</span>
            </div>
            <div className="t-small muted-dark" style={{ marginBottom: 20 }}>
              Timer akan berjalan otomatis tanpa jeda antar sesi.
            </div>

            <Link
              href={`/recipes/${id}/brew`}
              className="btn btn--primary btn--xl btn--block"
              style={{ display: 'flex', justifyContent: 'center', textDecoration: 'none' }}
            >
              <IconPlay size={18} /> Mulai Brewing
            </Link>

            <div className="row gap-2" style={{ marginTop: 12 }}>
              <button className="btn btn--secondary-dark" style={{ flex: 1 }}>Praktek silent</button>
            </div>

            <hr className="divider--dark" style={{ margin: '24px 0' }} />

            <div className="t-label muted-dark" style={{ marginBottom: 10 }}>Persiapan</div>
            <ul style={{ margin: 0, padding: 0, listStyle: 'none' }}>
              {[
                `${recipe.params.dose} biji kopi`,
                `Air ${recipe.params.temp}, ${recipe.params.yield}`,
                `Grinder set: ${recipe.params.grind}`,
                `Timer YAVA (otomatis)`,
              ].map((line, i) => (
                <li key={i} className="row gap-2" style={{ padding: '6px 0', fontSize: 13 }}>
                  <span style={{
                    width: 18, height: 18, borderRadius: '50%',
                    border: '1.5px solid rgba(255,255,255,.30)',
                    display: 'inline-flex', alignItems: 'center', justifyContent: 'center', flex: 'none',
                  }}>
                    <span style={{ width: 4, height: 4, borderRadius: '50%', background: '#fff', opacity: .5 }} />
                  </span>
                  <span style={{ color: 'rgba(255,255,255,.85)' }}>{line}</span>
                </li>
              ))}
            </ul>
          </div>
        </aside>
      </div>
    </main>
  )
}
