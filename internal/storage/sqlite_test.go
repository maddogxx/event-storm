package storage_test

import (
	"path/filepath"
	"testing"

	"github.com/madd0gxx/event-storm/internal/domain"
	"github.com/madd0gxx/event-storm/internal/storage"
)

func openTestStore(t *testing.T) *storage.Store {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.db")
	store, err := storage.Open(path)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	return store
}

func TestProjectLifecycle(t *testing.T) {
	store := openTestStore(t)

	p, err := store.CreateProject("Pagamentos", "Domínio de pagamentos")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if p.ID == 0 {
		t.Fatal("expected assigned id")
	}

	p.Name = "Pagamentos v2"
	if err := store.UpdateProject(p); err != nil {
		t.Fatalf("update: %v", err)
	}

	got, err := store.ProjectByID(p.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Name != "Pagamentos v2" {
		t.Fatalf("got name %q, expected Pagamentos v2", got.Name)
	}

	if err := store.DeleteProject(p.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}

	all, err := store.ListProjects()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(all) != 0 {
		t.Fatalf("expected empty list, got %d", len(all))
	}
}

func TestDashboardCascade(t *testing.T) {
	store := openTestStore(t)

	project, err := store.CreateProject("Checkout", "")
	if err != nil {
		t.Fatalf("project: %v", err)
	}
	dash, err := store.CreateDashboard(project.ID, "Fluxo principal")
	if err != nil {
		t.Fatalf("dashboard: %v", err)
	}

	w, h := domain.DefaultSize()
	_, err = store.CreatePostIt(&domain.PostIt{
		DashboardID: dash.ID,
		Type:        domain.TypeDomainEvent,
		Description: "Pedido criado",
		X:           10, Y: 20, Width: w, Height: h,
	})
	if err != nil {
		t.Fatalf("postit: %v", err)
	}

	if err := store.DeleteProject(project.ID); err != nil {
		t.Fatalf("delete project: %v", err)
	}

	postits, err := store.ListPostIts(dash.ID)
	if err != nil {
		t.Fatalf("list postits: %v", err)
	}
	if len(postits) != 0 {
		t.Fatalf("expected cascade delete, got %d postits", len(postits))
	}
}

func TestPostItUpdate(t *testing.T) {
	store := openTestStore(t)
	project, _ := store.CreateProject("P", "")
	dash, _ := store.CreateDashboard(project.ID, "D")

	w, h := domain.DefaultSize()
	postit, err := store.CreatePostIt(&domain.PostIt{
		DashboardID: dash.ID,
		Type:        domain.TypeCommand,
		Description: "Iniciar checkout",
		X:           5, Y: 6, Width: w, Height: h,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	postit.X = 100
	postit.Y = 200
	postit.Description = "Iniciar pagamento"
	if err := store.UpdatePostIt(postit); err != nil {
		t.Fatalf("update: %v", err)
	}

	postits, _ := store.ListPostIts(dash.ID)
	if len(postits) != 1 {
		t.Fatalf("expected 1 postit, got %d", len(postits))
	}
	if postits[0].X != 100 || postits[0].Y != 200 {
		t.Fatalf("position not persisted: got (%v,%v)", postits[0].X, postits[0].Y)
	}
	if postits[0].Description != "Iniciar pagamento" {
		t.Fatalf("description not persisted: %q", postits[0].Description)
	}
}
