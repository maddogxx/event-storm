import type { Project } from '~/types'

export default defineEventHandler((): Project[] => {
  const db = useDb()
  return db
    .prepare('SELECT * FROM projects ORDER BY updatedAt DESC')
    .all() as Project[]
})
