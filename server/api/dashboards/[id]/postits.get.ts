import type { PostIt } from '~/types'

export default defineEventHandler((event): PostIt[] => {
  const dashboardId = Number(getRouterParam(event, 'id'))
  const db = useDb()
  return db
    .prepare('SELECT * FROM postits WHERE dashboardId = ? ORDER BY id ASC')
    .all(dashboardId) as PostIt[]
})
