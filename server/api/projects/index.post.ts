import type { CreateProjectInput, Project } from '~/types'

export default defineEventHandler(async (event): Promise<Project> => {
  const body = await readBody<CreateProjectInput>(event)

  if (!body?.name?.trim()) {
    throw createError({ statusCode: 400, statusMessage: 'Nome do projeto é obrigatório' })
  }

  const db = useDb()
  const result = db
    .prepare('INSERT INTO projects (name, description) VALUES (?, ?)')
    .run(body.name.trim(), body.description ?? null)

  return db
    .prepare('SELECT * FROM projects WHERE id = ?')
    .get(result.lastInsertRowid) as Project
})
