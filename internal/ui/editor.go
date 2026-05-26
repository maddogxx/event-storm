package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/madd0gxx/event-storm/internal/domain"
	"github.com/madd0gxx/event-storm/internal/export"
)

const (
	canvasWidth  = 4000
	canvasHeight = 3000
)

// editor owns the state of an open dashboard inside the canvas view.
type editor struct {
	app       *App
	project   *domain.Project
	dashboard *domain.Dashboard

	canvas  *fyne.Container
	widgets map[int64]*postItWidget
}

func (a *App) showEditor(p *domain.Project, d *domain.Dashboard) {
	postits, err := a.store.ListPostIts(d.ID)
	if err != nil {
		dialog.ShowError(fmt.Errorf("carregando post-its: %w", err), a.window)
		return
	}

	ed := &editor{
		app:       a,
		project:   p,
		dashboard: d,
		canvas:    container.NewWithoutLayout(),
		widgets:   map[int64]*postItWidget{},
	}

	background := canvas.NewRectangle(color.NRGBA{R: 0xF5, G: 0xF5, B: 0xF0, A: 0xFF})
	background.Resize(fyne.NewSize(canvasWidth, canvasHeight))
	background.Move(fyne.NewPos(0, 0))
	ed.canvas.Add(background)
	ed.canvas.Resize(fyne.NewSize(canvasWidth, canvasHeight))

	for _, p := range postits {
		ed.spawn(p)
	}

	a.swap(ed.layout())
}

func (e *editor) layout() fyne.CanvasObject {
	back := widget.NewButtonWithIcon("Dashboards", theme.NavigateBackIcon(), func() {
		e.app.showDashboards(e.project)
	})

	exportBtn := widget.NewButtonWithIcon("Exportar SVG", theme.DocumentSaveIcon(), e.exportSVG)
	exportBtn.Importance = widget.HighImportance

	title := widget.NewLabelWithStyle(
		fmt.Sprintf("%s · %s", e.project.Name, e.dashboard.Name),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	header := container.NewBorder(nil, nil, back, exportBtn, title)

	palette := e.buildPalette()

	scroll := container.NewScroll(e.canvas)
	scroll.SetMinSize(fyne.NewSize(720, 520))

	body := container.NewBorder(nil, nil, palette, nil, scroll)
	return container.NewPadded(container.NewBorder(
		container.NewVBox(header, widget.NewSeparator()),
		nil, nil, nil,
		body,
	))
}

func (e *editor) buildPalette() fyne.CanvasObject {
	heading := widget.NewLabelWithStyle("Post-its", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	hint := widget.NewLabel("Clique em um tipo para adicionar ao quadro.")
	hint.Wrapping = fyne.TextWrapWord

	items := container.NewVBox(heading, hint, widget.NewSeparator())
	for _, info := range domain.TypeCatalog() {
		items.Add(e.paletteEntry(info))
	}

	scroll := container.NewVScroll(items)
	scroll.SetMinSize(fyne.NewSize(240, 0))
	return scroll
}

func (e *editor) paletteEntry(info domain.TypeInfo) fyne.CanvasObject {
	swatch := canvas.NewRectangle(info.Fill)
	swatch.StrokeColor = info.Stroke
	swatch.StrokeWidth = 1.5
	swatch.CornerRadius = 4
	swatch.SetMinSize(fyne.NewSize(28, 28))

	label := widget.NewLabelWithStyle(info.Label, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	desc := widget.NewLabel(info.Description)
	desc.Wrapping = fyne.TextWrapWord

	addBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		e.createPostIt(info.Type)
	})

	textCol := container.NewVBox(label, desc)
	row := container.NewBorder(nil, nil, swatch, addBtn, textCol)
	return container.NewPadded(row)
}

func (e *editor) createPostIt(t domain.Type) {
	w, h := domain.DefaultSize()
	p := &domain.PostIt{
		DashboardID: e.dashboard.ID,
		Type:        t,
		Description: "",
		X:           80 + float32(len(e.widgets)%6)*40,
		Y:           80 + float32(len(e.widgets)%6)*40,
		Width:       w,
		Height:      h,
	}
	stored, err := e.app.store.CreatePostIt(p)
	if err != nil {
		dialog.ShowError(err, e.app.window)
		return
	}
	e.spawn(stored)
	e.promptEditPostIt(e.widgets[stored.ID])
}

func (e *editor) spawn(p *domain.PostIt) {
	w := newPostItWidget(p,
		func(pw *postItWidget) { e.promptEditPostIt(pw) },
		func(pw *postItWidget) {
			pw.Move(fyne.NewPos(pw.model.X, pw.model.Y))
		},
		func(pw *postItWidget) {
			if err := e.app.store.UpdatePostIt(pw.model); err != nil {
				dialog.ShowError(err, e.app.window)
			}
		},
	)
	w.Resize(fyne.NewSize(p.Width, p.Height))
	w.Move(fyne.NewPos(p.X, p.Y))
	e.canvas.Add(w)
	e.canvas.Refresh()
	e.widgets[p.ID] = w
}

func (e *editor) promptEditPostIt(w *postItWidget) {
	descEntry := widget.NewMultiLineEntry()
	descEntry.SetText(w.model.Description)
	descEntry.SetPlaceHolder("Descreva o evento, comando, política, etc.")
	descEntry.Wrapping = fyne.TextWrapWord

	catalog := domain.TypeCatalog()
	options := make([]string, len(catalog))
	current := ""
	for i, info := range catalog {
		options[i] = info.Label
		if info.Type == w.model.Type {
			current = info.Label
		}
	}
	typeSelect := widget.NewSelect(options, nil)
	typeSelect.SetSelected(current)

	var form dialog.Dialog
	deleteBtn := widget.NewButtonWithIcon("Remover post-it", theme.DeleteIcon(), func() {
		if form != nil {
			form.Hide()
		}
		e.confirmDeletePostIt(w)
	})

	form = dialog.NewForm("Editar post-it", "Salvar", "Cancelar",
		[]*widget.FormItem{
			{Text: "Tipo", Widget: typeSelect},
			{Text: "Descrição", Widget: descEntry},
			{Text: "", Widget: deleteBtn},
		},
		func(ok bool) {
			if !ok {
				return
			}
			w.model.Description = descEntry.Text
			for _, info := range catalog {
				if info.Label == typeSelect.Selected {
					w.model.Type = info.Type
					break
				}
			}
			if err := e.app.store.UpdatePostIt(w.model); err != nil {
				dialog.ShowError(err, e.app.window)
				return
			}
			w.updateModel()
		},
		e.app.window,
	)
	form.Resize(fyne.NewSize(540, 380))
	form.Show()
}

func (e *editor) confirmDeletePostIt(w *postItWidget) {
	dialog.ShowConfirm("Remover post-it", "Deseja remover este post-it?", func(ok bool) {
		if !ok {
			return
		}
		if err := e.app.store.DeletePostIt(w.model.ID); err != nil {
			dialog.ShowError(err, e.app.window)
			return
		}
		delete(e.widgets, w.model.ID)
		e.canvas.Remove(w)
		e.canvas.Refresh()
	}, e.app.window)
}

func (e *editor) exportSVG() {
	postits, err := e.app.store.ListPostIts(e.dashboard.ID)
	if err != nil {
		dialog.ShowError(err, e.app.window)
		return
	}

	save := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, e.app.window)
			return
		}
		if writer == nil {
			return
		}
		path := writer.URI().Path()
		_ = writer.Close()

		if err := export.Dashboard(path, e.dashboard, postits); err != nil {
			dialog.ShowError(err, e.app.window)
			return
		}
		dialog.ShowInformation("Exportado", fmt.Sprintf("Dashboard exportado para:\n%s", path), e.app.window)
	}, e.app.window)

	save.SetFileName(sanitizeFilename(e.dashboard.Name) + ".svg")
	save.SetFilter(storage.NewExtensionFileFilter([]string{".svg"}))
	save.Resize(fyne.NewSize(720, 520))
	save.Show()
}

func sanitizeFilename(name string) string {
	out := make([]rune, 0, len(name))
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z',
			r >= 'A' && r <= 'Z',
			r >= '0' && r <= '9',
			r == '-', r == '_':
			out = append(out, r)
		case r == ' ':
			out = append(out, '_')
		}
	}
	if len(out) == 0 {
		return "dashboard"
	}
	return string(out)
}
