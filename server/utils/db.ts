import Database from 'better-sqlite3'
import { existsSync, mkdirSync } from 'node:fs'
import { dirname, resolve } from 'node:path'

let db: Database.Database | null = null

const ensureDir = (filePath: string) => {
  const dir = dirname(filePath)
  if (!existsSync(dir)) {
    mkdirSync(dir, { recursive: true })
  }
}

const runMigrations = (instance: Database.Database) => {
  instance.exec(`
    CREATE TABLE IF NOT EXISTS projects (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      name TEXT NOT NULL,
      description TEXT,
      createdAt TEXT NOT NULL DEFAULT (datetime('now')),
      updatedAt TEXT NOT NULL DEFAULT (datetime('now'))
    );

    CREATE TABLE IF NOT EXISTS dashboards (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      projectId INTEGER NOT NULL,
      name TEXT NOT NULL,
      description TEXT,
      createdAt TEXT NOT NULL DEFAULT (datetime('now')),
      updatedAt TEXT NOT NULL DEFAULT (datetime('now')),
      FOREIGN KEY (projectId) REFERENCES projects(id) ON DELETE CASCADE
    );

    CREATE TABLE IF NOT EXISTS postits (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      dashboardId INTEGER NOT NULL,
      kind TEXT NOT NULL,
      title TEXT NOT NULL DEFAULT '',
      description TEXT,
      x REAL NOT NULL DEFAULT 0,
      y REAL NOT NULL DEFAULT 0,
      width REAL NOT NULL DEFAULT 160,
      height REAL NOT NULL DEFAULT 120,
      createdAt TEXT NOT NULL DEFAULT (datetime('now')),
      updatedAt TEXT NOT NULL DEFAULT (datetime('now')),
      FOREIGN KEY (dashboardId) REFERENCES dashboards(id) ON DELETE CASCADE
    );

    CREATE INDEX IF NOT EXISTS idx_dashboards_project ON dashboards(projectId);
    CREATE INDEX IF NOT EXISTS idx_postits_dashboard ON postits(dashboardId);
  `)
}

export const useDb = (): Database.Database => {
  if (db) return db

  const config = useRuntimeConfig()
  const dbPath = resolve(config.dbPath as string)
  ensureDir(dbPath)

  db = new Database(dbPath)
  db.pragma('journal_mode = WAL')
  db.pragma('foreign_keys = ON')

  runMigrations(db)
  return db
}
