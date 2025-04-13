export interface Task {
  id: number
  title: string
  description: string | null
  status: string
  due_date: string | null
  priority: string
  user_id: number
  created_at?: string
  updated_at?: string
  tags: Tag[]
}

export interface Tag {
  id: number
  name: string
}

export interface TaskFilter {
  status?: string
  priority?: string
  dueDateBefore?: string
  dueDateAfter?: string
  tagName?: string
}
