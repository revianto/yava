'use client'

import { use, useState } from 'react'
import { useRouter } from 'next/navigation'
import { notFound } from 'next/navigation'
import { RECIPES, TYPES } from '@/lib/mock-data'
import { IconArrowLeft, IconPlus } from '@/components/icons'

const STEPS = ['Info dasar', 'Alur brewing', 'Visibilitas']
const BREW_TYPES = TYPES.filter((t) => t !== 'Semua')

interface SessionDraft {
  id: string
  kind: 'session'
  name: string
  duration: number
  note: string
}

interface NoteDraft {
  id: string
  kind: 'note'
  content: string
}

type TimelineDraft = SessionDraft | NoteDraft

export default function EditRecipePage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params)
  const router = useRouter()

  const recipe = RECIPES.find((r) => r.id === id)
  if (!recipe) notFound()
  if (recipe.isDefault) notFound()

  const [step, setStep] = useState(0)

  const [name, setName] = useState(recipe.name)
  const [type, setType] = useState(recipe.type)
  const [description, setDescription] = useState(recipe.description)
  const [dose, setDose] = useState(recipe.params.dose)
  const [yieldVal, setYieldVal] = useState(recipe.params.yield)
  const [temp, setTemp] = useState(recipe.params.temp)
  const [grind, setGrind] = useState(recipe.params.grind)
  const [ratio, setRatio] = useState(recipe.params.ratio)

  const [timeline, setTimeline] = useState<TimelineDraft[]>(
    recipe.timeline.map((item, i) =>
      item.kind === 'session'
        ? { id: String(i), kind: 'session', name: item.name, duration: item.duration, note: item.note ?? '' }
        : { id: String(i), kind: 'note', content: item.content }
    )
  )

  const [visibility, setVisibility] = useState<'private' | 'public' | 'group'>(recipe.visibility)
  const [saved, setSaved] = useState(false)

  const addSession = () =>
    setTimeline((prev) => [...prev, { id: String(Date.now()), kind: 'session', name: '', duration: 60, note: '' }])

  const addNote = () =>
    setTimeline((prev) => [...prev, { id: String(Date.now()), kind: 'note', content: '' }])

  const removeItem = (id: string) => setTimeline((prev) => prev.filter((t) => t.id !== id))

  const updateSession = (id: string, field: keyof Omit<SessionDraft, 'id' | 'kind'>, val: string | number) =>
    setTimeline((prev) => prev.map((t) => (t.id === id && t.kind === 'session' ? { ...t, [field]: val } : t)))

  const updateNote = (id: string, content: string) =>
    setTimeline((prev) => prev.map((t) => (t.id === id && t.kind === 'note' ? { ...t, content } : t)))

  const sessionCount = timeline.filter((t) => t.kind === 'session').length
  let sessionIdx = 0

  const canNext = step === 0 ? name.trim() !== '' && type !== '' : true

  const handleSave = () => {
    setSaved(true)
    setTimeout(() => router.push(`/recipes/${id}`), 800)
  }

  return (
    <main className="container" style={{ paddingTop: 0, paddingBottom: 64 }}>
      <div style={{ maxWidth: 640, margin: '0 auto', paddingTop: 32 }}>
        {/* Back */}
        <button
          onClick={() => router.back()}
          style={{ background: 'transparent', border: 0, cursor: 'pointer', display: 'inline-flex', alignItems: 'center', gap: 6, color: 'var(--muted)', fontWeight: 500, fontSize: 13, padding: 0, marginBottom: 32 }}
        >
          <IconArrowLeft size={14} /> Batal
        </button>

        {/* Header */}
        <div style={{ marginBottom: 32 }}>
          <div className="t-label" style={{ color: 'var(--muted)', marginBottom: 8 }}>Edit resep</div>
          <div className="t-h1">{recipe.name}</div>
        </div>

        {/* Step pills */}
        <div className="row gap-2" style={{ marginBottom: 40, flexWrap: 'wrap' }}>
          {STEPS.map((label, i) => (
            <div key={i} className="row gap-2">
              <div style={{
                display: 'inline-flex', alignItems: 'center', gap: 8,
                padding: '6px 14px', borderRadius: 99,
                background: i === step ? 'var(--deep-ink)' : i < step ? 'var(--hairline)' : 'transparent',
                border: i > step ? '1.5px solid var(--hairline)' : 'none',
                color: i === step ? '#fff' : 'var(--muted)',
                fontWeight: 700, fontSize: 12, letterSpacing: '.04em', textTransform: 'uppercase' as const,
              }}>
                <span style={{
                  width: 20, height: 20, borderRadius: '50%',
                  background: i === step ? 'rgba(255,255,255,.2)' : i < step ? 'var(--coral-red)' : 'var(--hairline)',
                  color: i < step ? '#fff' : 'inherit',
                  display: 'inline-flex', alignItems: 'center', justifyContent: 'center',
                  fontSize: 11, fontWeight: 700,
                }}>{i + 1}</span>
                {label}
              </div>
              {i < STEPS.length - 1 && <div style={{ width: 16, height: 1, background: 'var(--hairline)' }} />}
            </div>
          ))}
        </div>

        {/* Step 0 — Info dasar */}
        {step === 0 && (
          <div className="col gap-3">
            <div>
              <div className="t-label muted" style={{ marginBottom: 8 }}>Nama resep *</div>
              <input className="input" value={name} onChange={(e) => setName(e.target.value)} placeholder="e.g. V60 Kenya AA" style={{ width: '100%' }} />
            </div>

            <div>
              <div className="t-label muted" style={{ marginBottom: 8 }}>Jenis kopi *</div>
              <div className="row gap-2 wrap">
                {BREW_TYPES.map((t) => (
                  <button key={t} className={`tab${type === t ? ' tab--active' : ''}`} onClick={() => setType(t)}>{t}</button>
                ))}
              </div>
            </div>

            <div>
              <div className="t-label muted" style={{ marginBottom: 8 }}>Deskripsi</div>
              <textarea
                className="input"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="Ceritakan profil rasa dan karakteristik resep ini..."
                rows={3}
                style={{ width: '100%', height: 'auto', resize: 'vertical', paddingTop: 10, paddingBottom: 10 }}
              />
            </div>

            <div className="card" style={{ padding: 20 }}>
              <div className="t-label muted" style={{ marginBottom: 16 }}>Parameter brew</div>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 16 }}>
                {[
                  { label: 'Dose', val: dose, set: setDose, placeholder: '18g' },
                  { label: 'Yield', val: yieldVal, set: setYieldVal, placeholder: '36g' },
                  { label: 'Suhu', val: temp, set: setTemp, placeholder: '93°C' },
                  { label: 'Grind', val: grind, set: setGrind, placeholder: 'Medium-Fine' },
                  { label: 'Ratio', val: ratio, set: setRatio, placeholder: '1:16' },
                ].map(({ label, val, set, placeholder }) => (
                  <div key={label}>
                    <div className="t-label muted" style={{ marginBottom: 6 }}>{label}</div>
                    <input className="input" value={val} onChange={(e) => set(e.target.value)} placeholder={placeholder} style={{ width: '100%' }} />
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}

        {/* Step 1 — Alur brewing */}
        {step === 1 && (
          <div className="col gap-3">
            <div className="t-body muted" style={{ marginBottom: 4 }}>
              {sessionCount} sesi · drag untuk reorder (segera hadir)
            </div>

            {timeline.map((item) => {
              if (item.kind === 'session') {
                sessionIdx++
                const idx = sessionIdx
                return (
                  <div key={item.id} className="card" style={{ padding: 20 }}>
                    <div className="row between" style={{ marginBottom: 16 }}>
                      <div className="t-label">Sesi {idx}</div>
                      <button
                        onClick={() => removeItem(item.id)}
                        style={{ background: 'transparent', border: 0, cursor: 'pointer', color: 'var(--muted)', fontSize: 13, fontWeight: 700 }}
                      >
                        Hapus
                      </button>
                    </div>
                    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12 }}>
                      <div>
                        <div className="t-label muted" style={{ marginBottom: 6 }}>Nama sesi</div>
                        <input className="input" value={item.name} onChange={(e) => updateSession(item.id, 'name', e.target.value)} placeholder="e.g. Bloom" style={{ width: '100%' }} />
                      </div>
                      <div>
                        <div className="t-label muted" style={{ marginBottom: 6 }}>Durasi (detik)</div>
                        <input className="input t-mono-num" type="number" value={item.duration} onChange={(e) => updateSession(item.id, 'duration', parseInt(e.target.value) || 0)} style={{ width: '100%' }} />
                      </div>
                      <div style={{ gridColumn: '1 / -1' }}>
                        <div className="t-label muted" style={{ marginBottom: 6 }}>Catatan (opsional)</div>
                        <input className="input" value={item.note} onChange={(e) => updateSession(item.id, 'note', e.target.value)} placeholder="Instruksi atau tips untuk sesi ini" style={{ width: '100%' }} />
                      </div>
                    </div>
                  </div>
                )
              }

              return (
                <div key={item.id} className="card" style={{ padding: 20, background: 'rgba(45,82,74,.04)' }}>
                  <div className="row between" style={{ marginBottom: 12 }}>
                    <div className="t-label muted">Catatan</div>
                    <button
                      onClick={() => removeItem(item.id)}
                      style={{ background: 'transparent', border: 0, cursor: 'pointer', color: 'var(--muted)', fontSize: 13, fontWeight: 700 }}
                    >
                      Hapus
                    </button>
                  </div>
                  <textarea
                    className="input"
                    value={item.content}
                    onChange={(e) => updateNote(item.id, e.target.value)}
                    placeholder="Catatan atau instruksi tambahan..."
                    rows={2}
                    style={{ width: '100%', height: 'auto', resize: 'vertical', paddingTop: 8, paddingBottom: 8, fontSize: 13 }}
                  />
                </div>
              )
            })}

            <div className="row gap-2">
              <button className="btn btn--secondary" style={{ flex: 1, border: '1.5px dashed var(--deep-ink)' }} onClick={addSession}>
                <IconPlus size={16} /> Tambah sesi
              </button>
              <button className="btn btn--secondary" style={{ flex: 1, border: '1.5px dashed var(--hairline)', color: 'var(--muted)' }} onClick={addNote}>
                <IconPlus size={16} /> Tambah catatan
              </button>
            </div>
          </div>
        )}

        {/* Step 2 — Visibilitas */}
        {step === 2 && (
          <div className="col gap-3">
            <div className="t-body muted" style={{ marginBottom: 8 }}>Siapa yang bisa melihat resep ini?</div>
            {[
              { val: 'private', label: 'Pribadi', desc: 'Hanya kamu yang bisa melihat.' },
              { val: 'public', label: 'Publik', desc: 'Semua orang bisa melihat di halaman Explore.' },
              { val: 'group', label: 'Grup', desc: 'Hanya anggota grup yang kamu pilih.' },
            ].map(({ val, label, desc }) => (
              <button
                key={val}
                onClick={() => setVisibility(val as typeof visibility)}
                className="card"
                style={{
                  padding: 20, textAlign: 'left', cursor: 'pointer', width: '100%',
                  border: visibility === val ? '2px solid var(--deep-ink)' : undefined,
                  background: visibility === val ? 'var(--lavender-fog)' : 'var(--white)',
                }}
              >
                <div style={{ fontWeight: 700, fontSize: 16, marginBottom: 4 }}>{label}</div>
                <div className="t-small muted">{desc}</div>
              </button>
            ))}
          </div>
        )}

        {/* Navigation */}
        <div className="row between" style={{ marginTop: 40 }}>
          <button
            className="btn btn--secondary"
            onClick={() => setStep((s) => s - 1)}
            disabled={step === 0}
            style={{ opacity: step === 0 ? .4 : 1 }}
          >
            Sebelumnya
          </button>

          {step < STEPS.length - 1 ? (
            <button className="btn btn--primary" onClick={() => setStep((s) => s + 1)} disabled={!canNext} style={{ opacity: canNext ? 1 : .4 }}>
              Selanjutnya
            </button>
          ) : (
            <button className="btn btn--primary" onClick={handleSave} disabled={saved} style={{ opacity: saved ? .6 : 1 }}>
              {saved ? 'Menyimpan…' : 'Simpan perubahan'}
            </button>
          )}
        </div>
      </div>
    </main>
  )
}
