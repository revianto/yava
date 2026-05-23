'use client'

import { use, useState, useEffect, useRef, useMemo } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { notFound } from 'next/navigation'
import { RECIPES, fmt } from '@/lib/mock-data'
import type { RecipeSession } from '@/types'
import {
  IconPlay, IconPause, IconSkip, IconReset, IconClose, IconCheck,
} from '@/components/icons'

type Phase = 'countdown' | 'running' | 'paused' | 'complete'
type BrewVariant = 'focus' | 'ambient' | 'editorial'

interface EngineState {
  phase: Phase
  countdownN: number
  sessionIdx: number
  remaining: number
  progress: number
  totalProgress: number
  elapsedTotal: number
  totalDur: number
  current: RecipeSession | undefined
  isPaused: boolean
  isRunning: boolean
  pause: () => void
  resume: () => void
  skip: () => void
  reset: () => void
}

function useBrewingEngine(sessions: RecipeSession[]): EngineState {
  const [phase, setPhase] = useState<Phase>('countdown')
  const [countdownN, setCountdownN] = useState(3)
  const [sessionIdx, setSessionIdx] = useState(0)
  const [elapsedInSession, setElapsedInSession] = useState(0)

  const refStart = useRef<number | null>(null)
  const refPausedAt = useRef<number | null>(null)
  const refRAF = useRef<number | null>(null)
  const refPhase = useRef(phase)
  const refSessionIdx = useRef(sessionIdx)
  refPhase.current = phase
  refSessionIdx.current = sessionIdx

  // 3s countdown
  useEffect(() => {
    if (phase !== 'countdown') return
    setCountdownN(3)
    let n = 3
    const tick = () => {
      n -= 1
      if (n <= 0) {
        setCountdownN(0)
        setPhase('running')
        setSessionIdx(0)
        setElapsedInSession(0)
        refStart.current = performance.now()
      } else {
        setCountdownN(n)
        timer = setTimeout(tick, 800)
      }
    }
    let timer = setTimeout(tick, 800)
    return () => clearTimeout(timer)
  }, [phase === 'countdown']) // eslint-disable-line

  // rAF loop
  useEffect(() => {
    if (phase !== 'running') return
    const loop = () => {
      if (refStart.current === null) return
      const now = performance.now()
      const elapsedSec = (now - refStart.current) / 1000
      const idx = refSessionIdx.current
      const cur = sessions[idx]
      if (!cur) return
      if (elapsedSec >= cur.duration) {
        const next = idx + 1
        if (next >= sessions.length) {
          setPhase('complete')
          setElapsedInSession(cur.duration)
          return
        }
        const overflowMs = (elapsedSec - cur.duration) * 1000
        refStart.current = now - overflowMs
        setSessionIdx(next)
        setElapsedInSession(0)
      } else {
        setElapsedInSession(elapsedSec)
      }
      refRAF.current = requestAnimationFrame(loop)
    }
    refRAF.current = requestAnimationFrame(loop)
    return () => { if (refRAF.current) cancelAnimationFrame(refRAF.current) }
  }, [phase, sessionIdx, sessions])

  const cur = sessions[sessionIdx]
  const remaining = cur ? Math.max(0, cur.duration - elapsedInSession) : 0
  const progress = cur ? Math.min(1, elapsedInSession / cur.duration) : 0
  const totalDur = sessions.reduce((a, s) => a + s.duration, 0)
  const elapsedTotal =
    sessions.slice(0, sessionIdx).reduce((a, s) => a + s.duration, 0) + (cur ? elapsedInSession : 0)
  const totalProgress = totalDur ? Math.min(1, elapsedTotal / totalDur) : 0

  const pause = () => {
    if (phase !== 'running') return
    refPausedAt.current = performance.now()
    setPhase('paused')
    if (refRAF.current) cancelAnimationFrame(refRAF.current)
  }
  const resume = () => {
    if (phase !== 'paused' || refPausedAt.current === null || refStart.current === null) return
    const pauseDur = performance.now() - refPausedAt.current
    refStart.current += pauseDur
    refPausedAt.current = null
    setPhase('running')
  }
  const skip = () => {
    if (phase !== 'running' && phase !== 'paused') return
    const next = sessionIdx + 1
    if (next >= sessions.length) { setPhase('complete'); return }
    setSessionIdx(next)
    setElapsedInSession(0)
    refStart.current = performance.now()
    if (phase === 'paused') setPhase('running')
  }
  const reset = () => {
    if (refRAF.current) cancelAnimationFrame(refRAF.current)
    setSessionIdx(0)
    setElapsedInSession(0)
    setPhase('countdown')
  }

  return {
    phase, countdownN, sessionIdx, current: cur,
    remaining, progress, totalProgress, elapsedTotal, totalDur,
    pause, resume, skip, reset,
    isPaused: phase === 'paused',
    isRunning: phase === 'running',
  }
}

// ─── Focus variant ───────────────────────────────────────────────────────────
function FocusView({ engine, recipe, sessions, onExit }: { engine: EngineState; recipe: (typeof RECIPES)[0]; sessions: RecipeSession[]; onExit: () => void }) {
  const { phase, countdownN, current, remaining, progress, totalProgress, sessionIdx, pause, resume, skip, reset, isPaused } = engine
  return (
    <div className="brewing" style={{ background: 'var(--abyss)', color: '#fff', display: 'flex', flexDirection: 'column' }}>
      <div className="row between" style={{ padding: '20px 32px' }}>
        <div className="row gap-3">
          <span className="logo" style={{ fontSize: 20, color: '#fff' }}>YAVA<span className="dot">.</span></span>
          <span className="t-small muted-dark" style={{ marginLeft: 4 }}>{recipe.name}</span>
        </div>
        <button className="icon-btn icon-btn--dark" onClick={onExit}><IconClose size={18} /></button>
      </div>

      {phase !== 'countdown' && phase !== 'complete' && (
        <div style={{ padding: '0 32px 24px' }}>
          <div className="row between" style={{ marginBottom: 8 }}>
            <span className="t-label muted-dark">Total brew</span>
            <span className="t-small muted-dark t-mono-num">{fmt(engine.elapsedTotal)} / {fmt(engine.totalDur)}</span>
          </div>
          <div className="progress"><i style={{ width: `${totalProgress * 100}%` }} /></div>
        </div>
      )}

      <div className="grow row center" style={{ padding: '0 32px' }}>
        {phase === 'countdown' && (
          <div className="col center" style={{ alignItems: 'center', textAlign: 'center' }}>
            <div className="t-label muted-dark" style={{ marginBottom: 24 }}>Siapkan peralatan Anda</div>
            <div key={countdownN} className="count-pop t-mono-num" style={{ fontSize: 280, fontWeight: 700, lineHeight: 1, color: 'var(--coral-red)', letterSpacing: '-.04em' }}>
              {countdownN}
            </div>
            <div className="t-body muted-dark" style={{ marginTop: 24 }}>Tarik napas. Tunggu hingga timer mulai.</div>
          </div>
        )}

        {(phase === 'running' || phase === 'paused') && current && (
          <div className="col center session-enter" key={sessionIdx} style={{ alignItems: 'center', textAlign: 'center', width: '100%', maxWidth: 880 }}>
            <div className="t-label" style={{ color: 'var(--lilac)', marginBottom: 18 }}>Sesi {sessionIdx + 1} dari {sessions.length}</div>
            <div className="t-display" style={{ fontSize: 72, color: '#fff', marginBottom: 28, letterSpacing: '-.025em' }}>{current.name}</div>
            <div className="t-mono-num" style={{ fontSize: 260, fontWeight: 700, lineHeight: .95, letterSpacing: '-.035em', color: 'var(--powder)', opacity: isPaused ? .5 : 1, transition: 'opacity 200ms' }}>
              {fmt(remaining)}
            </div>
            <div style={{ width: '100%', maxWidth: 520, marginTop: 28 }}>
              <div className="progress"><i style={{ width: `${progress * 100}%` }} /></div>
            </div>
            <div className="t-body muted-dark" style={{ marginTop: 22, maxWidth: 560 }}>{current.note || '—'}</div>
          </div>
        )}

        {phase === 'complete' && (
          <div className="col center session-enter" style={{ alignItems: 'center', textAlign: 'center' }}>
            <div className="t-label" style={{ color: 'var(--powder)', marginBottom: 18 }}>Brewing selesai</div>
            <div className="t-display" style={{ fontSize: 96, color: '#fff', letterSpacing: '-.03em', maxWidth: 920 }}>
              Selamat menikmati<span style={{ color: 'var(--coral-red)' }}>.</span>
            </div>
            <div className="t-body muted-dark" style={{ marginTop: 18, maxWidth: 520 }}>
              {sessions.length} sesi tuntas dalam {fmt(engine.totalDur)}. Diamkan satu menit sebelum sip pertama.
            </div>
            <div className="row gap-2" style={{ marginTop: 40 }}>
              <button className="btn btn--primary btn--lg" onClick={reset}><IconReset size={16} /> Ulangi</button>
              <button className="btn btn--secondary-dark btn--lg" onClick={onExit}>Kembali ke resep</button>
            </div>
          </div>
        )}
      </div>

      {(phase === 'running' || phase === 'paused') && (
        <div className="row center gap-2" style={{ padding: '32px 32px 48px' }}>
          <button className="icon-btn icon-btn--dark icon-btn--lg" onClick={reset}><IconReset size={20} /></button>
          {isPaused
            ? <button className="icon-btn icon-btn--lg" onClick={resume} style={{ background: 'var(--coral-red)', color: '#fff', borderColor: 'var(--coral-red)' }}><IconPlay size={24} /></button>
            : <button className="icon-btn icon-btn--lg" onClick={pause} style={{ background: '#fff', color: 'var(--deep-ink)', borderColor: '#fff' }}><IconPause size={24} /></button>
          }
          <button className="icon-btn icon-btn--dark icon-btn--lg" onClick={skip}><IconSkip size={20} /></button>
        </div>
      )}

      {phase === 'running' && sessions[sessionIdx + 1] && (
        <div className="row between" style={{ padding: '0 32px 32px', alignItems: 'center' }}>
          <span className="t-label muted-dark">Berikutnya</span>
          <span className="t-small muted-dark">
            <strong style={{ color: '#fff', fontWeight: 700 }}>{sessions[sessionIdx + 1].name}</strong> · {sessions[sessionIdx + 1].duration}s
          </span>
        </div>
      )}
    </div>
  )
}

// ─── Ambient variant ─────────────────────────────────────────────────────────
function AmbientView({ engine, recipe, sessions, onExit }: { engine: EngineState; recipe: (typeof RECIPES)[0]; sessions: RecipeSession[]; onExit: () => void }) {
  const { phase, countdownN, current, remaining, progress, totalProgress, sessionIdx, pause, resume, skip, reset, isPaused } = engine
  return (
    <div className="brewing" style={{ background: 'var(--lavender-fog)', color: 'var(--deep-ink)', display: 'flex', flexDirection: 'column' }}>
      <div className="row between" style={{ padding: '20px 32px' }}>
        <div className="row gap-3">
          <span className="logo" style={{ fontSize: 20 }}>YAVA<span className="dot">.</span></span>
          <span className="t-small muted" style={{ marginLeft: 4 }}>{recipe.name}</span>
        </div>
        <button className="icon-btn" onClick={onExit}><IconClose size={18} /></button>
      </div>

      <div className="grow row center" style={{ padding: '0 48px' }}>
        {phase === 'countdown' && (
          <div className="col center" style={{ alignItems: 'center', textAlign: 'center' }}>
            <div className="t-label muted" style={{ marginBottom: 16 }}>Siapkan peralatan</div>
            <div key={countdownN} className="count-pop t-mono-num" style={{ fontSize: 320, fontWeight: 700, lineHeight: 1, color: 'var(--deep-ink)', letterSpacing: '-.04em' }}>{countdownN}</div>
            <div className="t-body muted" style={{ marginTop: 20 }}>Timer akan mulai dalam beberapa detik.</div>
          </div>
        )}

        {(phase === 'running' || phase === 'paused') && current && (
          <div style={{ width: '100%', maxWidth: 1080, display: 'grid', gridTemplateColumns: '1.4fr 1fr', gap: 32, alignItems: 'stretch' }}>
            <div className="card card--abyss card--hero session-enter" key={sessionIdx} style={{ padding: 40, color: '#fff', display: 'flex', flexDirection: 'column', justifyContent: 'space-between' }}>
              <div className="row between">
                <span className="t-label" style={{ color: 'var(--lilac)' }}>Sesi {sessionIdx + 1} / {sessions.length}</span>
                <span className="t-label muted-dark t-mono-num">{fmt(engine.elapsedTotal)} total</span>
              </div>
              <div style={{ paddingTop: 16 }}>
                <div className="t-h1" style={{ color: '#fff', marginBottom: 8 }}>{current.name}</div>
                <div className="t-body muted-dark">{current.note}</div>
              </div>
              <div className="col gap-3" style={{ marginTop: 20 }}>
                <div className="t-mono-num" style={{ fontSize: 200, fontWeight: 700, lineHeight: 1, color: 'var(--powder)', letterSpacing: '-.035em', opacity: isPaused ? .5 : 1, transition: 'opacity 200ms' }}>{fmt(remaining)}</div>
                <div className="progress"><i style={{ width: `${progress * 100}%` }} /></div>
              </div>
            </div>

            <div className="col gap-3">
              <div className="card" style={{ padding: 20 }}>
                <div className="t-label muted" style={{ marginBottom: 12 }}>Alur</div>
                <div className="col">
                  {sessions.map((s, i) => {
                    const done = i < sessionIdx
                    const active = i === sessionIdx
                    return (
                      <div key={i} className="row gap-3" style={{ padding: '12px 0', borderTop: i === 0 ? 0 : '1px solid var(--hairline)', opacity: done ? .45 : 1 }}>
                        <span style={{
                          width: 28, height: 28, borderRadius: '50%',
                          background: active ? 'var(--coral-red)' : done ? 'var(--deep-ink)' : 'transparent',
                          border: active || done ? '0' : '1.5px solid var(--deep-ink)',
                          color: active || done ? '#fff' : 'var(--deep-ink)',
                          display: 'inline-flex', alignItems: 'center', justifyContent: 'center',
                          fontWeight: 700, fontSize: 12,
                        }}>{done ? <IconCheck size={14} /> : i + 1}</span>
                        <div className="grow">
                          <div style={{ fontWeight: 700, fontSize: 14 }}>{s.name}</div>
                          <div className="t-caption">{s.duration}s</div>
                        </div>
                        {active && <span className="live-dot" />}
                      </div>
                    )
                  })}
                </div>
              </div>

              <div className="card" style={{ padding: 20 }}>
                <div className="t-label muted" style={{ marginBottom: 12 }}>Kontrol</div>
                <div className="row gap-2">
                  {isPaused
                    ? <button className="btn btn--primary btn--block" onClick={resume}><IconPlay size={16} /> Lanjutkan</button>
                    : <button className="btn btn--light-secondary btn--block" onClick={pause}><IconPause size={16} /> Jeda</button>
                  }
                </div>
                <div className="row gap-2" style={{ marginTop: 8 }}>
                  <button className="btn btn--secondary" style={{ flex: 1 }} onClick={skip}><IconSkip size={14} /> Skip sesi</button>
                  <button className="icon-btn" onClick={reset}><IconReset size={18} /></button>
                </div>
              </div>

              <div className="card card--electric" style={{ padding: 18, color: 'var(--lilac)' }}>
                <div className="t-label" style={{ color: 'var(--lilac)', marginBottom: 4 }}>Berikutnya</div>
                <div className="t-h3" style={{ color: '#fff' }}>{sessions[sessionIdx + 1] ? sessions[sessionIdx + 1].name : 'Selesai brewing'}</div>
                <div className="t-small" style={{ color: 'var(--lilac)' }}>{sessions[sessionIdx + 1] ? `${sessions[sessionIdx + 1].duration}s` : 'Nikmati seteguk pertama'}</div>
              </div>
            </div>
          </div>
        )}

        {phase === 'complete' && (
          <div className="col center" style={{ alignItems: 'center', textAlign: 'center', maxWidth: 760 }}>
            <div className="t-label muted" style={{ marginBottom: 18 }}>Brewing selesai</div>
            <div className="t-display" style={{ fontSize: 88, letterSpacing: '-.03em' }}>Selamat menikmati<span style={{ color: 'var(--coral-red)' }}>.</span></div>
            <div className="t-body muted" style={{ marginTop: 16, marginBottom: 36 }}>{sessions.length} sesi tuntas dalam {fmt(engine.totalDur)}.</div>
            <div className="row gap-2">
              <button className="btn btn--primary btn--lg" onClick={reset}><IconReset size={16} /> Ulangi</button>
              <button className="btn btn--secondary btn--lg" onClick={onExit}>Kembali ke resep</button>
            </div>
          </div>
        )}
      </div>

      {(phase === 'running' || phase === 'paused') && (
        <div style={{ padding: '0 32px 32px' }}>
          <div className="row between" style={{ marginBottom: 8 }}>
            <span className="t-label muted">Total brew</span>
            <span className="t-small muted t-mono-num">{fmt(engine.elapsedTotal)} / {fmt(engine.totalDur)}</span>
          </div>
          <div className="progress progress--light"><i style={{ width: `${totalProgress * 100}%` }} /></div>
        </div>
      )}
    </div>
  )
}

// ─── Editorial variant ───────────────────────────────────────────────────────
function EditorialView({ engine, recipe, sessions, onExit }: { engine: EngineState; recipe: (typeof RECIPES)[0]; sessions: RecipeSession[]; onExit: () => void }) {
  const { phase, countdownN, current, remaining, progress, totalProgress, sessionIdx, pause, resume, skip, reset, isPaused } = engine
  return (
    <div className="brewing grid-paper-bg" style={{ color: 'var(--deep-ink)', display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
      <div className="row between" style={{ padding: '20px 40px' }}>
        <div className="row gap-3">
          <span className="logo" style={{ fontSize: 20 }}>YAVA<span className="dot">.</span></span>
          <span className="t-small muted">— EDISI BREWING · {recipe.name}</span>
        </div>
        <button className="icon-btn" onClick={onExit}><IconClose size={18} /></button>
      </div>

      <div className="grow" style={{ position: 'relative', padding: '0 40px' }}>
        {phase === 'countdown' && (
          <div style={{ height: '100%', display: 'grid', gridTemplateColumns: '1fr 1fr', alignItems: 'center', gap: 48 }}>
            <div>
              <div className="t-label" style={{ color: 'var(--coral-red)', marginBottom: 8 }}>NO. 03 — PREPARATION</div>
              <div className="t-display" style={{ fontSize: 96, letterSpacing: '-.035em', lineHeight: .95 }}>
                Tarik napas.<br />Mulai dari <span style={{ color: 'var(--coral-red)' }}>aroma</span>.
              </div>
              <div className="t-body muted" style={{ marginTop: 16, maxWidth: 480 }}>Tiga detik untuk hadir. Tidak lebih, tidak kurang.</div>
            </div>
            <div className="col center" style={{ alignItems: 'center', padding: 40 }}>
              <div key={countdownN} className="count-pop t-mono-num" style={{ fontSize: 460, fontWeight: 700, lineHeight: 1, color: 'var(--coral-red)', letterSpacing: '-.05em' }}>{countdownN}</div>
            </div>
          </div>
        )}

        {(phase === 'running' || phase === 'paused') && current && (
          <div key={sessionIdx} className="session-enter" style={{ display: 'grid', gridTemplateColumns: '1fr 1.2fr', height: '100%', alignItems: 'stretch', gap: 0, position: 'relative' }}>
            <div className="col" style={{ padding: '24px 0', justifyContent: 'space-between' }}>
              <div>
                <div className="row gap-2" style={{ marginBottom: 12 }}>
                  <span className="tag tag--espresso">SESI {String(sessionIdx + 1).padStart(2, '0')}</span>
                  <span className="tag tag--default">{recipe.type.toUpperCase()}</span>
                </div>
                <div className="t-display" style={{ fontSize: 112, letterSpacing: '-.035em', lineHeight: .9 }}>
                  {current.name}<span style={{ color: 'var(--coral-red)' }}>.</span>
                </div>
                <div className="t-body" style={{ marginTop: 18, maxWidth: 460, fontSize: 17 }}>{current.note}</div>
              </div>
              <div>
                <div className="row gap-1" style={{ marginBottom: 20, flexWrap: 'wrap' }}>
                  {sessions.map((s, i) => (
                    <div key={i} style={{ height: 6, flex: 1, minWidth: 40, borderRadius: 99, background: i < sessionIdx ? 'var(--deep-ink)' : i === sessionIdx ? 'var(--coral-red)' : 'rgba(26,21,48,.15)', position: 'relative', overflow: 'hidden' }}>
                      {i === sessionIdx && (
                        <div style={{ position: 'absolute', inset: 0, background: 'var(--deep-ink)', width: `${progress * 100}%`, transition: 'width 100ms linear' }} />
                      )}
                    </div>
                  ))}
                </div>
                <div className="row gap-2">
                  {isPaused
                    ? <button className="btn btn--primary btn--lg" onClick={resume}><IconPlay size={16} /> Lanjutkan</button>
                    : <button className="btn btn--light-primary btn--lg" onClick={pause}><IconPause size={16} /> Jeda</button>
                  }
                  <button className="btn btn--secondary btn--lg" onClick={skip}><IconSkip size={14} /> Sesi berikut</button>
                  <button className="icon-btn icon-btn--lg" onClick={reset}><IconReset size={20} /></button>
                </div>
                <div className="t-small muted" style={{ marginTop: 16 }}>
                  Berikutnya: <strong style={{ color: 'var(--deep-ink)' }}>{sessions[sessionIdx + 1] ? sessions[sessionIdx + 1].name : 'selesai — diamkan satu menit'}</strong>
                </div>
              </div>
            </div>

            <div style={{ position: 'relative', background: 'var(--abyss)', color: 'var(--powder)', borderRadius: 16, margin: '16px 0 16px 32px', padding: 40, display: 'flex', flexDirection: 'column', justifyContent: 'space-between', overflow: 'hidden' }}>
              <div className="row between">
                <span className="t-label" style={{ color: 'var(--lilac)' }}>REMAINING</span>
                <span className="t-label" style={{ color: 'var(--lilac)' }}>{Math.round(progress * 100)}%</span>
              </div>
              <div className="t-mono-num" style={{ fontSize: 360, fontWeight: 700, lineHeight: .85, letterSpacing: '-.04em', textAlign: 'center', opacity: isPaused ? .5 : 1, transition: 'opacity 200ms' }}>
                {fmt(remaining)}
              </div>
              <div>
                <div className="row between" style={{ marginBottom: 8 }}>
                  <span className="t-label" style={{ color: 'var(--lilac)' }}>TOTAL</span>
                  <span className="t-small t-mono-num" style={{ color: 'var(--powder)' }}>{fmt(engine.elapsedTotal)} / {fmt(engine.totalDur)}</span>
                </div>
                <div className="progress"><i style={{ width: `${totalProgress * 100}%` }} /></div>
              </div>
            </div>
          </div>
        )}

        {phase === 'complete' && (
          <div style={{ height: '100%', display: 'grid', gridTemplateColumns: '1fr 1fr', alignItems: 'center', gap: 48 }}>
            <div>
              <div className="t-label" style={{ color: 'var(--coral-red)', marginBottom: 8 }}>NO. 09 — FINIS</div>
              <div className="t-display" style={{ fontSize: 120, letterSpacing: '-.035em', lineHeight: .9 }}>Selamat<br />menikmati<span style={{ color: 'var(--coral-red)' }}>.</span></div>
              <div className="t-body muted" style={{ marginTop: 20, maxWidth: 520, fontSize: 17 }}>{sessions.length} sesi tuntas dalam {fmt(engine.totalDur)}. Diamkan satu menit.</div>
              <div className="row gap-2" style={{ marginTop: 32 }}>
                <button className="btn btn--primary btn--lg" onClick={reset}><IconReset size={16} /> Ulangi</button>
                <button className="btn btn--secondary btn--lg" onClick={onExit}>Kembali ke resep</button>
              </div>
            </div>
            <div style={{ background: 'var(--coral-red)', color: '#fff', borderRadius: 16, padding: 40, display: 'flex', flexDirection: 'column', justifyContent: 'space-between', minHeight: 480 }}>
              <div className="t-label">RINGKASAN</div>
              <div className="col gap-3">
                <div>
                  <div className="t-label" style={{ opacity: .7, marginBottom: 4 }}>Total waktu</div>
                  <div className="t-mono-num" style={{ fontSize: 96, fontWeight: 700, lineHeight: 1, letterSpacing: '-.03em' }}>{fmt(engine.totalDur)}</div>
                </div>
                <hr style={{ border: 0, borderTop: '1.5px solid rgba(255,255,255,.30)' }} />
                <div className="col gap-1">
                  {sessions.map((s, i) => (
                    <div key={i} className="row between" style={{ padding: '6px 0', fontSize: 14 }}>
                      <span style={{ fontWeight: 700 }}>{i + 1}. {s.name}</span>
                      <span className="t-mono-num" style={{ fontWeight: 700 }}>{fmt(s.duration)}</span>
                    </div>
                  ))}
                </div>
              </div>
              <div className="t-small">{recipe.params.dose} · {recipe.params.yield} · {recipe.params.temp}</div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

// ─── Root ────────────────────────────────────────────────────────────────────
export default function BrewPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params)
  const router = useRouter()
  const [variant] = useState<BrewVariant>('focus')

  const recipe = RECIPES.find((r) => r.id === id)
  if (!recipe) notFound()

  const sessions = useMemo(
    () => recipe.timeline.filter((s): s is RecipeSession => s.kind === 'session'),
    [recipe]
  )
  if (sessions.length === 0) notFound()

  const engine = useBrewingEngine(sessions)

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') router.push(`/recipes/${id}`)
      if (e.code === 'Space') { e.preventDefault(); engine.isPaused ? engine.resume() : engine.pause() }
      if (e.key === 'ArrowRight') engine.skip()
      if (e.key === 'r' || e.key === 'R') engine.reset()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [engine.isPaused]) // eslint-disable-line

  const onExit = () => router.push(`/recipes/${id}`)
  const props = { engine, recipe, sessions, onExit }

  if (variant === 'ambient') return <AmbientView {...props} />
  if (variant === 'editorial') return <EditorialView {...props} />
  return <FocusView {...props} />
}
