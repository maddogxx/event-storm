// Event Storm is a Fyne-based desktop tool for designing Event Storming flows.
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/madd0gxx/event-storm/internal/storage"
	"github.com/madd0gxx/event-storm/internal/ui"
)

func main() {
	dbPath, err := resolveDBPath()
	if err != nil {
		log.Fatalf("resolving database path: %v", err)
	}

	store, err := storage.Open(dbPath)
	if err != nil {
		log.Fatalf("opening database at %s: %v", dbPath, err)
	}
	defer store.Close()

	ui.New(store).Run()
}

// resolveDBPath chooses an OS-appropriate location for the SQLite file and
// ensures the parent directory exists.
func resolveDBPath() (string, error) {
	if override := os.Getenv("EVENT_STORM_DB"); override != "" {
		if err := os.MkdirAll(filepath.Dir(override), 0o755); err != nil {
			return "", fmt.Errorf("mkdir db dir: %w", err)
		}
		return override, nil
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("user config dir: %w", err)
	}
	dir := filepath.Join(configDir, "event-storm")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("mkdir config dir: %w", err)
	}
	return filepath.Join(dir, "event-storm.db"), nil
}
