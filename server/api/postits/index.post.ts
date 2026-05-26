import type { CreatePostItInput, PostIt, PostItKind } from '~/types'

const VALID_KINDS: PostItKind[] = [
  'event',
  'command',
  'actor',
  'aggregate',
  'policy',
  'readmodel',
  'external',
  'hotspot'
]

export default defineEventHandler(async (event): Promise<PostIt> => {
  const body = await readBody<CreatePostItInput>(event)

  if (!body?.dashboardId || !VALID_KINDS.includes(body.kind)) {
    throw createError({ statusCode: 400, statusMessage: 'dashboardId e kind válido são obrigatórios' })
  }

  const db = useDb()
  const dashboard = db.prepare('SELECT id FROM dashboards WHERE id = ?').get(body.dashboardId)
  if (!dashboard) {
    throw createError({ statusCode: 404, statusMessage: 'Dashboard não encontrado' })
  }

  const result = db
    .prepare(`
      INSERT INTO postits (dashboardId, kind, title, description, x, y, width, height)
      VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `)
    .run(
      body.dashboardId,
      body.kind,
      body.title ?? '',
      body.description ?? null,
      body.x ?? 0,
      body.y ?? 0,
      body.width ?? 160,
      body.height ?? 120
    )

  return db
    .prepare('SELECT * FROM postits WHERE id = ?')
    .get(result.lastInsertRowid) as PostIt
})
