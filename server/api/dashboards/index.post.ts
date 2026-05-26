import type { CreateDashboardInput, Dashboard } from '~/types'

export default defineEventHandler(async (event): Promise<Dashboard> => {
  const body = await readBody<CreateDashboardInput>(event)

  if (!body?.name?.trim() || !body.projectId) {
    throw createError({ statusCode: 400, statusMessage: 'projectId e nome são obrigatórios' })
  }

  const db = useDb()
  const project = db.prepare('SELECT id FROM projects WHERE id = ?').get(body.projectId)
  if (!project) {
    throw createError({ statusCode: 404, statusMessage: 'Projeto não encontrado' })
  }

  const result = db
    .prepare('INSERT INTO dashboards (projectId, name, description) VALUES (?, ?, ?)')
    .run(body.projectId, body.name.trim(), body.description ?? null)

  return db
    .prepare('SELECT * FROM dashboards WHERE id = ?')
    .get(result.lastInsertRowid) as Dashboard
})
