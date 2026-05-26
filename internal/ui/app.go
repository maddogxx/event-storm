// Package ui wires the Fyne views, navigation and storage interactions.
package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"github.com/madd0gxx/event-storm/internal/storage"
)

// App is the top-level controller that owns the Fyne window and the storage backend.
type App struct {
	fyne   fyne.App
	window fyne.Window
	store  *storage.Store
	stack  *fyne.Container
}

// New builds a wired-up application instance ready to be run.
func New(store *storage.Store) *App {
	fyneApp := app.NewWithID("br.dev.eventstorm")
	w := fyneApp.NewWindow("Event Storm")
	w.Resize(fyne.NewSize(1200, 780))
	w.CenterOnScreen()

	a := &App{
		fyne:   fyneApp,
		window: w,
		store:  store,
		stack:  container.NewStack(),
	}
	w.SetContent(a.stack)
	a.showProjects()
	return a
}

// Run blocks until the user closes the window.
func (a *App) Run() { a.window.ShowAndRun() }

// swap replaces the current view inside the main window stack.
func (a *App) swap(view fyne.CanvasObject) {
	a.stack.Objects = []fyne.CanvasObject{view}
	a.stack.Refresh()
}
