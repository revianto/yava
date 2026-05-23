export type RecipeVisibility = 'private' | 'public' | 'group'

export interface RecipeAuthor {
  name: string
  initials: string
}

export interface RecipeParams {
  dose: string
  yield: string
  temp: string
  grind: string
  ratio: string
}

export interface RecipeSession {
  kind: 'session'
  name: string
  duration: number
  note?: string
}

export interface RecipeNote {
  kind: 'note'
  content: string
}

export type TimelineItem = RecipeSession | RecipeNote

export interface Recipe {
  id: string
  type: string
  subtype: string
  name: string
  description: string
  tags: string[]
  visibility: RecipeVisibility
  isDefault: boolean
  params: RecipeParams
  timeline: TimelineItem[]
  author: RecipeAuthor
  saves: number
  lastBrewed: string
}
