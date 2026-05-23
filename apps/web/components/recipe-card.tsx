import Link from 'next/link'
import { Clock, Droplets, Bean, Eye, EyeOff, Lock, Archive } from 'lucide-react'
import { Card, CardContent, CardFooter } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { formatDuration } from '@/lib/utils'
import type { Recipe } from '@/types'

const visibilityConfig = {
  public: { label: 'Publik', icon: Eye, className: 'bg-green-100 text-green-700 border-green-200' },
  private: { label: 'Privat', icon: Lock, className: 'bg-muted text-muted-foreground' },
  group: { label: 'Grup', icon: EyeOff, className: 'bg-blue-100 text-blue-700 border-blue-200' },
}

export function RecipeCard({ recipe }: { recipe: Recipe }) {
  const vis = visibilityConfig[recipe.visibility]
  const VisIcon = vis.icon

  return (
    <Link href={`/recipes/${recipe.id}`} className="block group">
      <Card className="h-full transition-shadow hover:shadow-md border-border group-hover:border-ring/30">
        <CardContent className="p-4 space-y-3">
          <div className="flex items-start justify-between gap-2">
            <div className="flex-1 min-w-0">
              <p className="text-xs font-medium text-muted-foreground mb-1">{recipe.typeName}</p>
              <h3 className="font-semibold text-sm leading-snug line-clamp-2 group-hover:text-primary transition-colors">
                {recipe.title}
              </h3>
            </div>
            <div className="flex gap-1 shrink-0 flex-col items-end">
              {recipe.isDefault && (
                <Badge variant="outline" className="text-[10px] px-1.5 py-0 bg-accent/30 text-accent-foreground border-accent/40">
                  Default
                </Badge>
              )}
              {recipe.isArchived && (
                <Badge variant="outline" className="text-[10px] px-1.5 py-0 gap-0.5">
                  <Archive className="w-2.5 h-2.5" />
                  Arsip
                </Badge>
              )}
            </div>
          </div>

          {recipe.description && (
            <p className="text-xs text-muted-foreground line-clamp-2 leading-relaxed">
              {recipe.description}
            </p>
          )}

          <div className="flex flex-wrap gap-x-3 gap-y-1 text-xs text-muted-foreground">
            {recipe.coffeeGrams && (
              <span className="flex items-center gap-1">
                <Bean className="w-3 h-3" />
                {recipe.coffeeGrams}g
              </span>
            )}
            {recipe.waterMl && (
              <span className="flex items-center gap-1">
                <Droplets className="w-3 h-3" />
                {recipe.waterMl}ml
              </span>
            )}
            {recipe.totalDurationSeconds > 0 && (
              <span className="flex items-center gap-1">
                <Clock className="w-3 h-3" />
                {formatDuration(recipe.totalDurationSeconds)}
              </span>
            )}
          </div>
        </CardContent>

        <CardFooter className="px-4 pb-3 pt-0">
          <span className={`inline-flex items-center gap-1 text-[10px] font-medium px-2 py-0.5 rounded-full border ${vis.className}`}>
            <VisIcon className="w-2.5 h-2.5" />
            {vis.label}
          </span>
        </CardFooter>
      </Card>
    </Link>
  )
}
