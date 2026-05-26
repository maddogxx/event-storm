import type { Dashboard } from '~/types'

export default defineEventHandler((event): Dashboard[] => {
  const projectId = Number(getRouterParam(event, 'id'))
  const db = useDb()
  return db
    .prepare('SELECT * FROM dashboards WHERE projectId = ? ORDER BY updatedAt DESC')
    .all(projectId) as Dashboard[]
})
