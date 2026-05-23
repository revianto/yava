'use client'

import { use } from 'react'
import Link from 'next/link'
import { notFound } from 'next/navigation'
import {
  ArrowLeft, Clock, Droplets, Bean, Thermometer, Coffee,
  Play, Pencil, Copy, Archive, Eye, Lock, MoreVertical
} from 'lucide-react'
import { Button, buttonVariants } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { mockRecipes } from '@/lib/mock-data'
import { formatDuration, formatDate } from '@/lib/utils'
import { cn } from '@/lib/utils'

const visibilityLabel: Record<string, string> = {
  public: 'Publik',
  private: 'Privat',
  group: 'Grup',
}

export default function RecipeDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params)
  const recipe = mockRecipes.find((r) => r.id === id)

  if (!recipe) notFound()

  const allSteps = [
    ...recipe.sessions.map((s) => ({ type: 'session' as const, ...s })),
    ...recipe.notes.map((n) => ({ type: 'note' as const, id: n.id, order: n.order, content: n.content })),
  ].sort((a, b) => a.order - b.order)

  return (
    <div className="max-w-2xl mx-auto p-4 md:p-6 space-y-6">
      <Link
        href="/recipes"
        className={cn(buttonVariants({ variant: 'ghost', size: 'sm' }), '-ml-2 gap-2 text-muted-foreground')}
      >
        <ArrowLeft className="w-4 h-4" />
        Kembali
      </Link>

      <div className="space-y-3">
        <div className="flex items-start justify-between gap-3">
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-1.5">
              <span className="text-xs font-medium text-muted-foreground">{recipe.typeName}</span>
              {recipe.subtypeName && (
                <>
                  <span className="text-muted-foreground/40">·</span>
                  <span className="text-xs text-muted-foreground">{recipe.subtypeName}</span>
                </>
              )}
            </div>
            <h1 className="text-2xl font-semibold leading-tight">{recipe.title}</h1>
          </div>

          {!recipe.isDefault && (
            <DropdownMenu>
              <DropdownMenuTrigger className="inline-flex items-center justify-center w-8 h-8 rounded-lg hover:bg-muted transition-colors outline-none focus-visible:ring-2 focus-visible:ring-ring">
                <MoreVertical className="w-4 h-4" />
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem onClick={() => {}}>
                  <Link href={`/recipes/${id}/edit`} className="flex items-center gap-2 w-full">
                    <Pencil className="w-3.5 h-3.5" />
                    Edit
                  </Link>
                </DropdownMenuItem>
                <DropdownMenuItem className="gap-2">
                  <Copy className="w-3.5 h-3.5" />
                  Duplikasi
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem className="gap-2 text-muted-foreground">
                  <Archive className="w-3.5 h-3.5" />
                  {recipe.isArchived ? 'Restore' : 'Arsipkan'}
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          )}
        </div>

        {recipe.description && (
          <p className="text-sm text-muted-foreground leading-relaxed">{recipe.description}</p>
        )}

        <div className="flex flex-wrap gap-2">
          {recipe.isDefault && (
            <Badge variant="outline" className="text-xs bg-accent/30 border-accent/40">
              Default
            </Badge>
          )}
          {recipe.isArchived && (
            <Badge variant="outline" className="text-xs">
              <Archive className="w-3 h-3 mr-1" />
              Diarsipkan
            </Badge>
          )}
          <Badge variant="outline" className="text-xs gap-1">
            {recipe.visibility === 'public' ? <Eye className="w-3 h-3" /> : <Lock className="w-3 h-3" />}
            {visibilityLabel[recipe.visibility]}
          </Badge>
          <span className="text-xs text-muted-foreground self-center">
            {formatDate(recipe.createdAt)}
          </span>
        </div>
      </div>

      <Separator />

      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
        {recipe.coffeeBeans && (
          <div className="space-y-1">
            <p className="text-xs text-muted-foreground flex items-center gap-1">
              <Bean className="w-3 h-3" />
              Biji Kopi
            </p>
            <p className="text-sm font-medium">{recipe.coffeeBeans}</p>
          </div>
        )}
        {recipe.coffeeGrams && (
          <div className="space-y-1">
            <p className="text-xs text-muted-foreground flex items-center gap-1">
              <Coffee className="w-3 h-3" />
              Kopi
            </p>
            <p className="text-sm font-medium">{recipe.coffeeGrams}g</p>
          </div>
        )}
        {recipe.waterMl && (
          <div className="space-y-1">
            <p className="text-xs text-muted-foreground flex items-center gap-1">
              <Droplets className="w-3 h-3" />
              Air
            </p>
            <p className="text-sm font-medium">{recipe.waterMl}ml</p>
          </div>
        )}
        {recipe.waterTempC && (
          <div className="space-y-1">
            <p className="text-xs text-muted-foreground flex items-center gap-1">
              <Thermometer className="w-3 h-3" />
              Suhu
            </p>
            <p className="text-sm font-medium">{recipe.waterTempC}°C</p>
          </div>
        )}
        {recipe.grindSize && (
          <div className="space-y-1">
            <p className="text-xs text-muted-foreground">Gilingan</p>
            <p className="text-sm font-medium">{recipe.grindSize}</p>
          </div>
        )}
        {recipe.totalDurationSeconds > 0 && (
          <div className="space-y-1">
            <p className="text-xs text-muted-foreground flex items-center gap-1">
              <Clock className="w-3 h-3" />
              Total Waktu
            </p>
            <p className="text-sm font-medium">{formatDuration(recipe.totalDurationSeconds)}</p>
          </div>
        )}
      </div>

      <Separator />

      {allSteps.length > 0 && (
        <div className="space-y-3">
          <h2 className="text-sm font-semibold">Langkah Brewing</h2>
          <div className="space-y-2">
            {allSteps.map((step, idx) =>
              step.type === 'session' ? (
                <div key={step.id} className="flex items-start gap-3 p-3 rounded-lg border border-border bg-card">
                  <div className="flex items-center justify-center w-6 h-6 rounded-full bg-primary/10 text-primary text-xs font-semibold shrink-0 mt-0.5">
                    {idx + 1}
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center justify-between gap-2">
                      <p className="text-sm font-medium">{step.label}</p>
                      <div className="flex items-center gap-2 text-xs text-muted-foreground shrink-0">
                        {step.waterMl && (
                          <span className="flex items-center gap-0.5">
                            <Droplets className="w-3 h-3" />
                            {step.waterMl}ml
                          </span>
                        )}
                        <span className="flex items-center gap-0.5">
                          <Clock className="w-3 h-3" />
                          {formatDuration(step.durationSeconds)}
                        </span>
                      </div>
                    </div>
                    {step.notes && (
                      <p className="text-xs text-muted-foreground mt-0.5">{step.notes}</p>
                    )}
                  </div>
                </div>
              ) : (
                <div key={step.id} className="flex items-start gap-3 p-3 rounded-lg bg-accent/20 border border-accent/30">
                  <div className="w-1.5 h-1.5 rounded-full bg-accent-foreground/40 shrink-0 mt-2" />
                  <p className="text-xs text-foreground/70 leading-relaxed">{step.content}</p>
                </div>
              )
            )}
          </div>
        </div>
      )}

      {recipe.sessions.length > 0 && (
        <div className="pt-2">
          <Link
            href={`/recipes/${id}/brew`}
            className={cn(buttonVariants({ size: 'lg' }), 'w-full gap-2')}
          >
            <Play className="w-4 h-4" />
            Mulai Brewing
          </Link>
        </div>
      )}
    </div>
  )
}
