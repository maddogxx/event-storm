import type { CreateProjectInput, Project } from '~/types'

export default defineEventHandler(async (event): Promise<Project> => {
  const id = Number(getRouterParam(event, 'id'))
  const body = await readBody<CreateProjectInput>(event)

  if (!body?.name?.trim()) {
    throw createError({ statusCode: 400, statusMessage: 'Nome do projeto é obrigatório' })
  }

  const db = useDb()
  const result = db
    .prepare(`
      UPDATE projects
      SET name = ?, description = ?, updatedAt = datetime('now')
      WHERE id = ?
    `)
    .run(body.name.trim(), body.description ?? null, id)

  if (result.changes === 0) {
    throw createError({ statusCode: 404, statusMessage: 'Projeto não encontrado' })
  }

  return db.prepare('SELECT * FROM projects WHERE id = ?').get(id) as Project
})
