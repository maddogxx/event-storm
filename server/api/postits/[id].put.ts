import type { PostIt, UpdatePostItInput } from '~/types'

export default defineEventHandler(async (event): Promise<PostIt> => {
  const id = Number(getRouterParam(event, 'id'))
  const body = await readBody<UpdatePostItInput>(event)
  const db = useDb()

  const current = db.prepare('SELECT * FROM postits WHERE id = ?').get(id) as PostIt | undefined
  if (!current) {
    throw createError({ statusCode: 404, statusMessage: 'Post-it não encontrado' })
  }

  const next = {
    title: body.title ?? current.title,
    description: body.description ?? current.description,
    x: body.x ?? current.x,
    y: body.y ?? current.y,
    width: body.width ?? current.width,
    height: body.height ?? current.height
  }

  db.prepare(`
    UPDATE postits
    SET title = ?, description = ?, x = ?, y = ?, width = ?, height = ?, updatedAt = datetime('now')
    WHERE id = ?
  `).run(next.title, next.description, next.x, next.y, next.width, next.height, id)

  return db.prepare('SELECT * FROM postits WHERE id = ?').get(id) as PostIt
})
