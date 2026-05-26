// Package storage persists Event Storming projects, dashboards and post-its.
package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/madd0gxx/event-storm/internal/domain"

	_ "modernc.org/sqlite"
)

// Store is the SQLite-backed repository for the application.
type Store struct {
	db *sql.DB
}

// Open initialises a SQLite database at the given path and runs the migrations.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return s, nil
}

// Close releases the underlying database handle.
func (s *Store) Close() error { return s.db.Close() }

func (s *Store) migrate() error {
	schema := `
CREATE TABLE IF NOT EXISTS projects (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS dashboards (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name       TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS postits (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    dashboard_id INTEGER NOT NULL REFERENCES dashboards(id) ON DELETE CASCADE,
    type         TEXT NOT NULL,
    description  TEXT NOT NULL DEFAULT '',
    x            REAL NOT NULL,
    y            REAL NOT NULL,
    width        REAL NOT NULL,
    height       REAL NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_dashboards_project ON dashboards(project_id);
CREATE INDEX IF NOT EXISTS idx_postits_dashboard  ON postits(dashboard_id);
`
	_, err := s.db.Exec(schema)
	return err
}

// ---------- Projects ----------

// CreateProject persists a new project and returns it with id/timestamps filled.
func (s *Store) CreateProject(name, description string) (*domain.Project, error) {
	res, err := s.db.Exec(`INSERT INTO projects(name, description) VALUES(?, ?)`, name, description)
	if err != nil {
		return nil, fmt.Errorf("insert project: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("last id: %w", err)
	}
	return s.ProjectByID(id)
}

// UpdateProject saves edits to an existing project row.
func (s *Store) UpdateProject(p *domain.Project) error {
	_, err := s.db.Exec(`UPDATE projects SET name = ?, description = ? WHERE id = ?`,
		p.Name, p.Description, p.ID)
	if err != nil {
		return fmt.Errorf("update project: %w", err)
	}
	return nil
}

// DeleteProject removes the project and cascades into its dashboards and post-its.
func (s *Store) DeleteProject(id int64) error {
	_, err := s.db.Exec(`DELETE FROM projects WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete project: %w", err)
	}
	return nil
}

// ProjectByID loads a single project by primary key.
func (s *Store) ProjectByID(id int64) (*domain.Project, error) {
	row := s.db.QueryRow(`SELECT id, name, description, created_at FROM projects WHERE id = ?`, id)
	p := &domain.Project{}
	var created string
	if err := row.Scan(&p.ID, &p.Name, &p.Description, &created); err != nil {
		return nil, fmt.Errorf("scan project: %w", err)
	}
	p.CreatedAt = parseTime(created)
	return p, nil
}

// ListProjects returns every project, most recent first.
func (s *Store) ListProjects() ([]*domain.Project, error) {
	rows, err := s.db.Query(`SELECT id, name, description, created_at FROM projects ORDER BY created_at DESC, id DESC`)
	if err != nil {
		return nil, fmt.Errorf("query projects: %w", err)
	}
	defer rows.Close()

	var out []*domain.Project
	for rows.Next() {
		p := &domain.Project{}
		var created string
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &created); err != nil {
			return nil, fmt.Errorf("scan project: %w", err)
		}
		p.CreatedAt = parseTime(created)
		out = append(out, p)
	}
	return out, rows.Err()
}

// ---------- Dashboards ----------

// CreateDashboard creates a dashboard inside a project.
func (s *Store) CreateDashboard(projectID int64, name string) (*domain.Dashboard, error) {
	res, err := s.db.Exec(`INSERT INTO dashboards(project_id, name) VALUES(?, ?)`, projectID, name)
	if err != nil {
		return nil, fmt.Errorf("insert dashboard: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("last id: %w", err)
	}
	return s.DashboardByID(id)
}

// UpdateDashboard renames an existing dashboard.
func (s *Store) UpdateDashboard(d *domain.Dashboard) error {
	_, err := s.db.Exec(`UPDATE dashboards SET name = ? WHERE id = ?`, d.Name, d.ID)
	if err != nil {
		return fmt.Errorf("update dashboard: %w", err)
	}
	return nil
}

// DeleteDashboard removes a dashboard and its post-its.
func (s *Store) DeleteDashboard(id int64) error {
	_, err := s.db.Exec(`DELETE FROM dashboards WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete dashboard: %w", err)
	}
	return nil
}

// DashboardByID loads a dashboard by primary key.
func (s *Store) DashboardByID(id int64) (*domain.Dashboard, error) {
	row := s.db.QueryRow(`SELECT id, project_id, name, created_at FROM dashboards WHERE id = ?`, id)
	d := &domain.Dashboard{}
	var created string
	if err := row.Scan(&d.ID, &d.ProjectID, &d.Name, &created); err != nil {
		return nil, fmt.Errorf("scan dashboard: %w", err)
	}
	d.CreatedAt = parseTime(created)
	return d, nil
}

// ListDashboards returns every dashboard for a given project.
func (s *Store) ListDashboards(projectID int64) ([]*domain.Dashboard, error) {
	rows, err := s.db.Query(`SELECT id, project_id, name, created_at FROM dashboards WHERE project_id = ? ORDER BY created_at DESC, id DESC`, projectID)
	if err != nil {
		return nil, fmt.Errorf("query dashboards: %w", err)
	}
	defer rows.Close()

	var out []*domain.Dashboard
	for rows.Next() {
		d := &domain.Dashboard{}
		var created string
		if err := rows.Scan(&d.ID, &d.ProjectID, &d.Name, &created); err != nil {
			return nil, fmt.Errorf("scan dashboard: %w", err)
		}
		d.CreatedAt = parseTime(created)
		out = append(out, d)
	}
	return out, rows.Err()
}

// ---------- Post-its ----------

// CreatePostIt persists a new post-it on a dashboard.
func (s *Store) CreatePostIt(p *domain.PostIt) (*domain.PostIt, error) {
	res, err := s.db.Exec(
		`INSERT INTO postits(dashboard_id, type, description, x, y, width, height) VALUES(?, ?, ?, ?, ?, ?, ?)`,
		p.DashboardID, string(p.Type), p.Description, p.X, p.Y, p.Width, p.Height,
	)
	if err != nil {
		return nil, fmt.Errorf("insert postit: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("last id: %w", err)
	}
	p.ID = id
	return p, nil
}

// UpdatePostIt persists position, size, description and type changes for a post-it.
func (s *Store) UpdatePostIt(p *domain.PostIt) error {
	_, err := s.db.Exec(
		`UPDATE postits SET type = ?, description = ?, x = ?, y = ?, width = ?, height = ? WHERE id = ?`,
		string(p.Type), p.Description, p.X, p.Y, p.Width, p.Height, p.ID,
	)
	if err != nil {
		return fmt.Errorf("update postit: %w", err)
	}
	return nil
}

// DeletePostIt removes a single post-it by id.
func (s *Store) DeletePostIt(id int64) error {
	_, err := s.db.Exec(`DELETE FROM postits WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete postit: %w", err)
	}
	return nil
}

// ListPostIts returns every post-it belonging to a dashboard.
func (s *Store) ListPostIts(dashboardID int64) ([]*domain.PostIt, error) {
	rows, err := s.db.Query(
		`SELECT id, dashboard_id, type, description, x, y, width, height FROM postits WHERE dashboard_id = ? ORDER BY id ASC`,
		dashboardID,
	)
	if err != nil {
		return nil, fmt.Errorf("query postits: %w", err)
	}
	defer rows.Close()

	var out []*domain.PostIt
	for rows.Next() {
		p := &domain.PostIt{}
		var typ string
		if err := rows.Scan(&p.ID, &p.DashboardID, &typ, &p.Description, &p.X, &p.Y, &p.Width, &p.Height); err != nil {
			return nil, fmt.Errorf("scan postit: %w", err)
		}
		p.Type = domain.Type(typ)
		out = append(out, p)
	}
	return out, rows.Err()
}

// parseTime tolerates the multiple text formats SQLite uses for CURRENT_TIMESTAMP.
func parseTime(raw string) time.Time {
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04:05.999999999-07:00",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t
		}
	}
	return time.Time{}
}
