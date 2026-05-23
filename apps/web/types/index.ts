export type RecipeVisibility = 'private' | 'public' | 'group'

export interface User {
  id: string
  name: string
  email: string
  avatarUrl?: string
}

export interface RecipeType {
  id: string
  name: string
  slug: string
}

export interface RecipeSubtype {
  id: string
  typeId: string
  name: string
}

export interface RecipeSession {
  id: string
  recipeId: string
  order: number
  label: string
  durationSeconds: number
  waterMl?: number
  notes?: string
}

export interface RecipeNote {
  id: string
  recipeId: string
  order: number
  content: string
}

export interface Recipe {
  id: string
  owner?: User
  typeId: string
  typeName: string
  subtypeName?: string
  title: string
  description?: string
  coffeeBeans?: string
  coffeeGrams?: number
  waterMl?: number
  grindSize?: string
  waterTempC?: number
  visibility: RecipeVisibility
  isDefault: boolean
  isArchived: boolean
  sessions: RecipeSession[]
  notes: RecipeNote[]
  totalDurationSeconds: number
  createdAt: string
}
