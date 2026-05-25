export type RecipeVisibility = 'private' | 'public' | 'group'
export type GroupRole = 'founder' | 'admin' | 'member'
export type GroupRecipeStatus = 'pending' | 'active' | 'rejected'

export interface GroupMember {
  id: string
  name: string
  initials: string
  role: GroupRole
  joinedAt: string
}

export interface GroupRecipeItem {
  id: string
  recipeId: string
  recipeName: string
  recipeType: string
  recipeSubtype: string
  submittedBy: string
  submittedByInitials: string
  submittedAt: string
  status: GroupRecipeStatus
  rejectionReason?: string
}

export interface Group {
  id: string
  name: string
  description: string
  inviteCode: string
  myRole: GroupRole
  createdAt: string
  members: GroupMember[]
  recipes: GroupRecipeItem[]
}

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

export interface DiscussionReply {
  id: string
  authorName: string
  authorInitials: string
  content: string
  createdAt: string
}

export interface Discussion {
  id: string
  recipeId: string
  authorName: string
  authorInitials: string
  content: string
  createdAt: string
  pinned: boolean
  replies: DiscussionReply[]
}

export type NotificationType = 'approved' | 'rejected' | 'reply' | 'joined'

export interface Notification {
  id: string
  type: NotificationType
  title: string
  body: string
  link?: string
  read: boolean
  createdAt: string
}

export interface Recipe {
  id: string
  type: string
  subtype: string
  name: string
  description: string
  tags: string[]
  visibility: RecipeVisibility
  isDefault: boolean
  isArchived?: boolean
  params: RecipeParams
  timeline: TimelineItem[]
  author: RecipeAuthor
  saves: number
  lastBrewed: string
}
