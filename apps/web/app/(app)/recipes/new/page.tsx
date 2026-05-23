'use client'

import { useState } from 'react'
import Link from 'next/link'
import { ArrowLeft, Plus, Trash2, ChevronRight, ChevronLeft } from 'lucide-react'
import { Button, buttonVariants } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { recipeTypes } from '@/lib/mock-data'
import { cn } from '@/lib/utils'

const STEPS = ['Info Dasar', 'Sessions', 'Visibilitas']

interface SessionDraft {
  id: string
  label: string
  durationSeconds: number
  waterMl: string
  notes: string
}

export default function NewRecipePage() {
  const [step, setStep] = useState(0)

  const [title, setTitle] = useState('')
  const [typeId, setTypeId] = useState('')
  const [description, setDescription] = useState('')
  const [coffeeBeans, setCoffeeBeans] = useState('')
  const [coffeeGrams, setCoffeeGrams] = useState('')
  const [waterMl, setWaterMl] = useState('')
  const [grindSize, setGrindSize] = useState('')
  const [waterTempC, setWaterTempC] = useState('')

  const [sessions, setSessions] = useState<SessionDraft[]>([
    { id: '1', label: 'Bloom', durationSeconds: 45, waterMl: '', notes: '' },
  ])

  const [visibility, setVisibility] = useState<'private' | 'public'>('private')

  const addSession = () => {
    setSessions((prev) => [
      ...prev,
      { id: String(Date.now()), label: '', durationSeconds: 60, waterMl: '', notes: '' },
    ])
  }

  const removeSession = (id: string) => {
    setSessions((prev) => prev.filter((s) => s.id !== id))
  }

  const updateSession = (id: string, field: keyof SessionDraft, value: string | number) => {
    setSessions((prev) => prev.map((s) => (s.id === id ? { ...s, [field]: value } : s)))
  }

  const canNext = step === 0 ? title.trim() !== '' && typeId !== '' : true

  return (
    <div className="max-w-xl mx-auto p-4 md:p-6 space-y-6">
      <Link
        href="/recipes"
        className={cn(buttonVariants({ variant: 'ghost', size: 'sm' }), '-ml-2 gap-2 text-muted-foreground')}
      >
        <ArrowLeft className="w-4 h-4" />
        Batal
      </Link>

      <div>
        <h1 className="text-xl font-semibold">Resep Baru</h1>
        <p className="text-sm text-muted-foreground mt-0.5">Buat resep kopi kamu sendiri</p>
      </div>

      <div className="flex items-center gap-2">
        {STEPS.map((label, i) => (
          <div key={i} className="flex items-center gap-2">
            <div className={cn(
              'flex items-center justify-center w-6 h-6 rounded-full text-xs font-medium',
              i <= step ? 'bg-primary text-primary-foreground' : 'bg-muted text-muted-foreground'
            )}>
              {i + 1}
            </div>
            <span className={cn(
              'text-xs hidden sm:inline',
              i === step ? 'text-foreground font-medium' : 'text-muted-foreground'
            )}>
              {label}
            </span>
            {i < STEPS.length - 1 && <div className="w-4 h-px bg-border" />}
          </div>
        ))}
      </div>

      {step === 0 && (
        <div className="space-y-4">
          <div className="space-y-1.5">
            <Label htmlFor="title">Nama Resep <span className="text-destructive">*</span></Label>
            <Input id="title" value={title} onChange={(e) => setTitle(e.target.value)} placeholder="e.g. V60 My Favorite" />
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="type">Jenis Kopi <span className="text-destructive">*</span></Label>
            <Select value={typeId} onValueChange={(v) => setTypeId(v ?? '')}>
              <SelectTrigger id="type" className="w-full">
                <SelectValue placeholder="Pilih jenis..." />
              </SelectTrigger>
              <SelectContent>
                {recipeTypes.map((t) => (
                  <SelectItem key={t.id} value={t.id}>{t.name}</SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="desc">Deskripsi</Label>
            <Textarea id="desc" value={description} onChange={(e) => setDescription(e.target.value)} placeholder="Ceritakan sedikit tentang resep ini..." rows={3} />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <Label htmlFor="beans">Biji Kopi</Label>
              <Input id="beans" value={coffeeBeans} onChange={(e) => setCoffeeBeans(e.target.value)} placeholder="e.g. Ethiopia Yirgacheffe" />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="grind">Ukuran Gilingan</Label>
              <Input id="grind" value={grindSize} onChange={(e) => setGrindSize(e.target.value)} placeholder="e.g. Medium-coarse" />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="cgrams">Kopi (gram)</Label>
              <Input id="cgrams" type="number" value={coffeeGrams} onChange={(e) => setCoffeeGrams(e.target.value)} placeholder="20" />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="wml">Air (ml)</Label>
              <Input id="wml" type="number" value={waterMl} onChange={(e) => setWaterMl(e.target.value)} placeholder="300" />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="temp">Suhu Air (°C)</Label>
              <Input id="temp" type="number" value={waterTempC} onChange={(e) => setWaterTempC(e.target.value)} placeholder="93" />
            </div>
          </div>
        </div>
      )}

      {step === 1 && (
        <div className="space-y-4">
          <p className="text-sm text-muted-foreground">
            Tambahkan langkah-langkah brewing. Timer akan auto-advance ke session berikutnya.
          </p>

          <div className="space-y-3">
            {sessions.map((s, i) => (
              <div key={s.id} className="p-4 rounded-lg border border-border space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-xs font-medium text-muted-foreground">Session {i + 1}</span>
                  {sessions.length > 1 && (
                    <button onClick={() => removeSession(s.id)} className="text-muted-foreground hover:text-destructive">
                      <Trash2 className="w-3.5 h-3.5" />
                    </button>
                  )}
                </div>

                <div className="grid grid-cols-2 gap-3">
                  <div className="space-y-1">
                    <Label className="text-xs">Label</Label>
                    <Input
                      value={s.label}
                      onChange={(e) => updateSession(s.id, 'label', e.target.value)}
                      placeholder="e.g. Bloom"
                      className="h-8 text-sm"
                    />
                  </div>
                  <div className="space-y-1">
                    <Label className="text-xs">Durasi (detik)</Label>
                    <Input
                      type="number"
                      value={s.durationSeconds}
                      onChange={(e) => updateSession(s.id, 'durationSeconds', parseInt(e.target.value) || 0)}
                      className="h-8 text-sm"
                    />
                  </div>
                  <div className="space-y-1">
                    <Label className="text-xs">Air (ml)</Label>
                    <Input
                      type="number"
                      value={s.waterMl}
                      onChange={(e) => updateSession(s.id, 'waterMl', e.target.value)}
                      placeholder="Optional"
                      className="h-8 text-sm"
                    />
                  </div>
                  <div className="space-y-1">
                    <Label className="text-xs">Catatan</Label>
                    <Input
                      value={s.notes}
                      onChange={(e) => updateSession(s.id, 'notes', e.target.value)}
                      placeholder="Optional"
                      className="h-8 text-sm"
                    />
                  </div>
                </div>
              </div>
            ))}
          </div>

          <Button variant="outline" size="sm" onClick={addSession} className="gap-2 w-full">
            <Plus className="w-4 h-4" />
            Tambah Session
          </Button>
        </div>
      )}

      {step === 2 && (
        <div className="space-y-4">
          <p className="text-sm text-muted-foreground">Siapa yang bisa melihat resep ini?</p>

          <div className="space-y-3">
            {[
              { value: 'private', label: 'Privat', desc: 'Hanya kamu yang bisa melihat.' },
              { value: 'public', label: 'Publik', desc: 'Semua orang bisa melihat di halaman Explore.' },
            ].map((opt) => (
              <button
                key={opt.value}
                onClick={() => setVisibility(opt.value as 'private' | 'public')}
                className={cn(
                  'w-full text-left p-4 rounded-lg border transition-colors',
                  visibility === opt.value
                    ? 'border-primary bg-primary/5'
                    : 'border-border hover:border-ring/40'
                )}
              >
                <p className="text-sm font-medium">{opt.label}</p>
                <p className="text-xs text-muted-foreground mt-0.5">{opt.desc}</p>
              </button>
            ))}
          </div>
        </div>
      )}

      <div className="flex justify-between pt-2">
        <Button
          variant="outline"
          onClick={() => setStep((s) => s - 1)}
          disabled={step === 0}
          className="gap-2"
        >
          <ChevronLeft className="w-4 h-4" />
          Sebelumnya
        </Button>

        {step < STEPS.length - 1 ? (
          <Button onClick={() => setStep((s) => s + 1)} disabled={!canNext} className="gap-2">
            Selanjutnya
            <ChevronRight className="w-4 h-4" />
          </Button>
        ) : (
          <Button className="gap-2">Simpan Resep</Button>
        )}
      </div>
    </div>
  )
}
