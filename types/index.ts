export interface Project {
  id: number
  name: string
  description: string | null
  createdAt: string
  updatedAt: string
}

export interface Dashboard {
  id: number
  projectId: number
  name: string
  description: string | null
  createdAt: string
  updatedAt: string
}

export type PostItKind =
  | 'event'
  | 'command'
  | 'actor'
  | 'aggregate'
  | 'policy'
  | 'readmodel'
  | 'external'
  | 'hotspot'

export interface PostItType {
  kind: PostItKind
  label: string
  description: string
  color: string
  textColor: string
}

export interface PostIt {
  id: number
  dashboardId: number
  kind: PostItKind
  title: string
  description: string | null
  x: number
  y: number
  width: number
  height: number
  createdAt: string
  updatedAt: string
}

export interface CreateProjectInput {
  name: string
  description?: string | null
}

export interface CreateDashboardInput {
  projectId: number
  name: string
  description?: string | null
}

export interface CreatePostItInput {
  dashboardId: number
  kind: PostItKind
  title?: string
  description?: string | null
  x: number
  y: number
  width?: number
  height?: number
}

export interface UpdatePostItInput {
  title?: string
  description?: string | null
  x?: number
  y?: number
  width?: number
  height?: number
}
