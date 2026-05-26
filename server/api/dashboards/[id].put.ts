import type { Dashboard } from '~/types'

interface UpdateDashboardBody {
  name?: string
  description?: string | null
}

export default defineEventHandler(async (event): Promise<Dashboard> => {
  const id = Number(getRouterParam(event, 'id'))
  const body = await readBody<UpdateDashboardBody>(event)

  if (!body?.name?.trim()) {
    throw createError({ statusCode: 400, statusMessage: 'Nome do dashboard é obrigatório' })
  }

  const db = useDb()
  const result = db
    .prepare(`
      UPDATE dashboards
      SET name = ?, description = ?, updatedAt = datetime('now')
      WHERE id = ?
    `)
    .run(body.name.trim(), body.description ?? null, id)

  if (result.changes === 0) {
    throw createError({ statusCode: 404, statusMessage: 'Dashboard não encontrado' })
  }

  return db.prepare('SELECT * FROM dashboards WHERE id = ?').get(id) as Dashboard
})
