export default defineEventHandler((event) => {
  const id = Number(getRouterParam(event, 'id'))
  const db = useDb()
  const result = db.prepare('DELETE FROM dashboards WHERE id = ?').run(id)

  if (result.changes === 0) {
    throw createError({ statusCode: 404, statusMessage: 'Dashboard não encontrado' })
  }
  return { success: true }
})
