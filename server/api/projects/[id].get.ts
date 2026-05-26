import type { Project } from '~/types'

export default defineEventHandler((event): Project => {
  const id = Number(getRouterParam(event, 'id'))
  const db = useDb()
  const project = db
    .prepare('SELECT * FROM projects WHERE id = ?')
    .get(id) as Project | undefined

  if (!project) {
    throw createError({ statusCode: 404, statusMessage: 'Projeto não encontrado' })
  }
  return project
})
