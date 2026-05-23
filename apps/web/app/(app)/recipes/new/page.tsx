'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { TYPES } from '@/lib/mock-data'
import { IconArrowLeft, IconPlus } from '@/components/icons'

const STEPS = ['Info dasar', 'Alur brewing', 'Visibilitas']

interface SessionDraft {
  id: string
  name: string
  duration: number
  note: string
}

const BREW_TYPES = TYPES.filter((t) => t !== 'Semua')

export default function NewRecipePage() {
  const router = useRouter()
  const [step, setStep] = useState(0)

  const [name, setName] = useState('')
  const [type, setType] = useState('')
  const [description, setDescription] = useState('')
  const [dose, setDose] = useState('')
  const [yieldVal, setYieldVal] = useState('')
  const [temp, setTemp] = useState('')
  const [grind, setGrind] = useState('')
  const [ratio, setRatio] = useState('')

  const [sessions, setSessions] = useState<SessionDraft[]>([
    { id: '1', name: 'Bloom', duration: 45, note: '' },
  ])

  const [visibility, setVisibility] = useState<'private' | 'public' | 'group'>('private')

  const addSession = () => setSessions((prev) => [...prev, { id: String(Date.now()), name: '', duration: 60, note: '' }])
  const removeSession = (id: string) => setSessions((prev) => prev.filter((s) => s.id !== id))
  const updateSession = (id: string, field: keyof SessionDraft, val: string | number) =>
    setSessions((prev) => prev.map((s) => (s.id === id ? { ...s, [field]: val } : s)))

  const canNext = step === 0 ? name.trim() !== '' && type !== '' : true

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
          <div className="t-label" style={{ color: 'var(--muted)', marginBottom: 8 }}>Resep baru</div>
          <div className="t-h1">Buat resep kopi</div>
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
                fontWeight: 700, fontSize: 12, letterSpacing: '.04em', textTransform: 'uppercase',
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
                  <button
                    key={t}
                    className={`tab${type === t ? ' tab--active' : ''}`}
                    onClick={() => setType(t)}
                  >{t}</button>
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
            <div className="t-body muted" style={{ marginBottom: 8 }}>
              Tambahkan sesi timer. Timer otomatis lanjut ke sesi berikutnya tanpa jeda.
            </div>

            {sessions.map((s, i) => (
              <div key={s.id} className="card" style={{ padding: 20 }}>
                <div className="row between" style={{ marginBottom: 16 }}>
                  <div className="t-label">Sesi {i + 1}</div>
                  {sessions.length > 1 && (
                    <button
                      onClick={() => removeSession(s.id)}
                      style={{ background: 'transparent', border: 0, cursor: 'pointer', color: 'var(--muted)', fontSize: 13, fontWeight: 700 }}
                    >
                      Hapus
                    </button>
                  )}
                </div>
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12 }}>
                  <div>
                    <div className="t-label muted" style={{ marginBottom: 6 }}>Nama sesi</div>
                    <input className="input" value={s.name} onChange={(e) => updateSession(s.id, 'name', e.target.value)} placeholder="e.g. Bloom" style={{ width: '100%' }} />
                  </div>
                  <div>
                    <div className="t-label muted" style={{ marginBottom: 6 }}>Durasi (detik)</div>
                    <input className="input t-mono-num" type="number" value={s.duration} onChange={(e) => updateSession(s.id, 'duration', parseInt(e.target.value) || 0)} style={{ width: '100%' }} />
                  </div>
                  <div style={{ gridColumn: '1 / -1' }}>
                    <div className="t-label muted" style={{ marginBottom: 6 }}>Catatan (opsional)</div>
                    <input className="input" value={s.note} onChange={(e) => updateSession(s.id, 'note', e.target.value)} placeholder="Instruksi atau tips untuk sesi ini" style={{ width: '100%' }} />
                  </div>
                </div>
              </div>
            ))}

            <button className="btn btn--secondary btn--block" onClick={addSession} style={{ border: '1.5px dashed var(--deep-ink)' }}>
              <IconPlus size={16} /> Tambah sesi
            </button>
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
                  background: visibility === val ? 'var(--lavender-fog)' : undefined,
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
            <button className="btn btn--primary" onClick={() => router.push('/')}>
              Simpan resep
            </button>
          )}
        </div>
      </div>
    </main>
  )
}
