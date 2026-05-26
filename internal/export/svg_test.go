package export_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/madd0gxx/event-storm/internal/domain"
	"github.com/madd0gxx/event-storm/internal/export"
)

func TestDashboardWritesSVG(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.svg")

	board := &domain.Dashboard{ID: 1, Name: "Fluxo Pedido"}
	postits := []*domain.PostIt{
		{Type: domain.TypeDomainEvent, Description: "Pedido criado", X: 0, Y: 0, Width: 160, Height: 100},
		{Type: domain.TypeCommand, Description: "Criar pedido", X: 200, Y: 50, Width: 160, Height: 100},
	}

	if err := export.Dashboard(path, board, postits); err != nil {
		t.Fatalf("export: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	content := string(data)

	for _, fragment := range []string{
		"<?xml",
		"<svg",
		"Fluxo Pedido",
		"Pedido criado",
		"Criar pedido",
		"</svg>",
	} {
		if !strings.Contains(content, fragment) {
			t.Errorf("expected svg to contain %q", fragment)
		}
	}
}

func TestDashboardEscapesXML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.svg")

	board := &domain.Dashboard{ID: 1, Name: "A & B <c>"}
	if err := export.Dashboard(path, board, nil); err != nil {
		t.Fatalf("export: %v", err)
	}
	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "A &amp; B &lt;c&gt;") {
		t.Errorf("expected escaped title in svg, got %s", string(data))
	}
}
