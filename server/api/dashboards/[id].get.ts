import type { Dashboard } from '~/types'

export default defineEventHandler((event): Dashboard => {
  const id = Number(getRouterParam(event, 'id'))
  const db = useDb()
  const dashboard = db
    .prepare('SELECT * FROM dashboards WHERE id = ?')
    .get(id) as Dashboard | undefined

  if (!dashboard) {
    throw createError({ statusCode: 404, statusMessage: 'Dashboard não encontrado' })
  }
  return dashboard
})
