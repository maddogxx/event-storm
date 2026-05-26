package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"github.com/madd0gxx/event-storm/internal/domain"
)

// postItWidget is a draggable, tappable representation of a domain.PostIt.
type postItWidget struct {
	widget.BaseWidget

	model    *domain.PostIt
	onTap    func(*postItWidget)
	onMoved  func(*postItWidget)
	onCommit func(*postItWidget)
}

func newPostItWidget(model *domain.PostIt, onTap, onMoved, onCommit func(*postItWidget)) *postItWidget {
	p := &postItWidget{
		model:    model,
		onTap:    onTap,
		onMoved:  onMoved,
		onCommit: onCommit,
	}
	p.ExtendBaseWidget(p)
	return p
}

// Dragged is invoked on every drag step; we accumulate the deltas into the model position.
func (p *postItWidget) Dragged(ev *fyne.DragEvent) {
	p.model.X += ev.Dragged.DX
	p.model.Y += ev.Dragged.DY
	if p.model.X < 0 {
		p.model.X = 0
	}
	if p.model.Y < 0 {
		p.model.Y = 0
	}
	if p.onMoved != nil {
		p.onMoved(p)
	}
}

// DragEnd is invoked once the drag gesture finishes; we persist the new position.
func (p *postItWidget) DragEnd() {
	if p.onCommit != nil {
		p.onCommit(p)
	}
}

// Tapped opens the post-it editor dialog.
func (p *postItWidget) Tapped(_ *fyne.PointEvent) {
	if p.onTap != nil {
		p.onTap(p)
	}
}

// TappedSecondary is a no-op but required so right-clicks do not crash the widget.
func (p *postItWidget) TappedSecondary(_ *fyne.PointEvent) {}

// MinSize ensures the widget reserves space matching the underlying model.
func (p *postItWidget) MinSize() fyne.Size {
	return fyne.NewSize(p.model.Width, p.model.Height)
}

// CreateRenderer wires up the Fyne renderer for the post-it.
func (p *postItWidget) CreateRenderer() fyne.WidgetRenderer {
	info := domain.InfoFor(p.model.Type)

	bg := canvas.NewRectangle(info.Fill)
	bg.StrokeColor = info.Stroke
	bg.StrokeWidth = 1.5
	bg.CornerRadius = 4

	header := canvas.NewText(info.Label, info.TextColor)
	header.TextStyle = fyne.TextStyle{Bold: true}
	header.TextSize = 11

	desc := widget.NewLabel(p.model.Description)
	desc.Wrapping = fyne.TextWrapWord
	desc.Alignment = fyne.TextAlignLeading

	return &postItRenderer{
		widget: p,
		bg:     bg,
		header: header,
		desc:   desc,
	}
}

// updateModel re-applies metadata changes (type/description) to the renderer.
func (p *postItWidget) updateModel() {
	p.Refresh()
}

type postItRenderer struct {
	widget *postItWidget
	bg     *canvas.Rectangle
	header *canvas.Text
	desc   *widget.Label
}

func (r *postItRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.bg.Move(fyne.NewPos(0, 0))

	r.header.Move(fyne.NewPos(8, 4))
	r.header.Resize(fyne.NewSize(size.Width-16, 18))

	r.desc.Move(fyne.NewPos(4, 22))
	r.desc.Resize(fyne.NewSize(size.Width-8, size.Height-26))
}

func (r *postItRenderer) MinSize() fyne.Size {
	return fyne.NewSize(r.widget.model.Width, r.widget.model.Height)
}

func (r *postItRenderer) Refresh() {
	info := domain.InfoFor(r.widget.model.Type)
	r.bg.FillColor = info.Fill
	r.bg.StrokeColor = info.Stroke
	r.bg.Refresh()

	r.header.Text = info.Label
	r.header.Color = info.TextColor
	r.header.Refresh()

	r.desc.SetText(r.widget.model.Description)
}

func (r *postItRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.header, r.desc}
}

func (r *postItRenderer) Destroy() {}

// Compile-time guarantees that the widget implements the interfaces Fyne expects.
var (
	_ fyne.Draggable      = (*postItWidget)(nil)
	_ fyne.Tappable       = (*postItWidget)(nil)
	_ fyne.WidgetRenderer = (*postItRenderer)(nil)
)
