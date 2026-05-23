'use client'

import { useState } from 'react'
import Link from 'next/link'
import { Plus, Search, Filter } from 'lucide-react'
import { buttonVariants } from '@/components/ui/button'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { RecipeCard } from '@/components/recipe-card'
import { mockRecipes, recipeTypes } from '@/lib/mock-data'
import { cn } from '@/lib/utils'

export default function RecipesPage() {
  const [search, setSearch] = useState('')
  const [activeType, setActiveType] = useState<string | null>(null)
  const [showArchived, setShowArchived] = useState(false)

  const filtered = mockRecipes.filter((r) => {
    if (!showArchived && r.isArchived) return false
    if (activeType && r.typeId !== activeType) return false
    if (search) {
      const q = search.toLowerCase()
      return (
        r.title.toLowerCase().includes(q) ||
        r.description?.toLowerCase().includes(q) ||
        r.typeName.toLowerCase().includes(q) ||
        r.coffeeBeans?.toLowerCase().includes(q)
      )
    }
    return true
  })

  return (
    <div className="p-4 md:p-6 space-y-5 max-w-6xl mx-auto">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-semibold">Resep Saya</h1>
          <p className="text-sm text-muted-foreground mt-0.5">
            {mockRecipes.filter((r) => !r.isArchived).length} resep aktif
          </p>
        </div>
        <Link href="/recipes/new" className={cn(buttonVariants({ size: 'sm' }), 'gap-1.5')}>
          <Plus className="w-4 h-4" />
          Tambah Resep
        </Link>
      </div>

      <div className="flex flex-col sm:flex-row gap-3">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
          <Input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Cari resep, biji kopi..."
            className="pl-9"
          />
        </div>
        <Button
          variant="outline"
          size="sm"
          onClick={() => setShowArchived((v) => !v)}
          className={cn('gap-2 shrink-0', showArchived && 'border-primary text-primary')}
        >
          <Filter className="w-4 h-4" />
          {showArchived ? 'Sembunyikan arsip' : 'Tampilkan arsip'}
        </Button>
      </div>

      <div className="flex gap-2 flex-wrap">
        <button
          onClick={() => setActiveType(null)}
          className={cn(
            'px-3 py-1 rounded-full text-xs font-medium border transition-colors',
            activeType === null
              ? 'bg-primary text-primary-foreground border-primary'
              : 'border-border text-muted-foreground hover:text-foreground'
          )}
        >
          Semua
        </button>
        {recipeTypes.map((t) => (
          <button
            key={t.id}
            onClick={() => setActiveType(activeType === t.id ? null : t.id)}
            className={cn(
              'px-3 py-1 rounded-full text-xs font-medium border transition-colors',
              activeType === t.id
                ? 'bg-primary text-primary-foreground border-primary'
                : 'border-border text-muted-foreground hover:text-foreground'
            )}
          >
            {t.name}
          </button>
        ))}
      </div>

      {filtered.length === 0 ? (
        <div className="text-center py-16 text-muted-foreground">
          <p className="text-sm">Tidak ada resep ditemukan.</p>
          {search && (
            <button onClick={() => setSearch('')} className="text-xs underline mt-1">
              Reset pencarian
            </button>
          )}
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          {filtered.map((r) => (
            <RecipeCard key={r.id} recipe={r} />
          ))}
        </div>
      )}
    </div>
  )
}
