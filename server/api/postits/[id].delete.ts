export default defineEventHandler((event) => {
  const id = Number(getRouterParam(event, 'id'))
  const db = useDb()
  const result = db.prepare('DELETE FROM postits WHERE id = ?').run(id)
  if (result.changes === 0) {
    throw createError({ statusCode: 404, statusMessage: 'Post-it não encontrado' })
  }
  return { success: true }
})
