package domain

import "time"

// Project groups related Event Storming dashboards under a single business context.
type Project struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   time.Time
}

// Dashboard is a single Event Storming board belonging to a project.
type Dashboard struct {
	ID        int64
	ProjectID int64
	Name      string
	CreatedAt time.Time
}
