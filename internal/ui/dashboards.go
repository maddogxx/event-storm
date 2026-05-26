package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/madd0gxx/event-storm/internal/domain"
)

func (a *App) showDashboards(p *domain.Project) {
	dashboards, err := a.store.ListDashboards(p.ID)
	if err != nil {
		dialog.ShowError(fmt.Errorf("carregando dashboards: %w", err), a.window)
		return
	}

	back := widget.NewButtonWithIcon("Projetos", theme.NavigateBackIcon(), func() {
		a.showProjects()
	})

	title := widget.NewLabelWithStyle(p.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	subtitle := widget.NewLabel("Dashboards do projeto. Cada dashboard representa um fluxo de Event Storming.")
	subtitle.Wrapping = fyne.TextWrapWord

	newBtn := widget.NewButtonWithIcon("Novo dashboard", theme.ContentAddIcon(), func() {
		a.promptNewDashboard(p)
	})
	newBtn.Importance = widget.HighImportance

	headerTop := container.NewBorder(nil, nil, back, newBtn, title)
	list := a.buildDashboardsList(p, dashboards)

	content := container.NewBorder(
		container.NewVBox(headerTop, subtitle, widget.NewSeparator()),
		nil, nil, nil,
		list,
	)
	a.swap(container.NewPadded(content))
}

func (a *App) buildDashboardsList(p *domain.Project, dashboards []*domain.Dashboard) fyne.CanvasObject {
	if len(dashboards) == 0 {
		empty := widget.NewLabel("Nenhum dashboard. Clique em \"Novo dashboard\" para iniciar um fluxo.")
		empty.Alignment = fyne.TextAlignCenter
		return container.NewCenter(empty)
	}

	items := container.NewVBox()
	for _, d := range dashboards {
		items.Add(a.dashboardCard(p, d))
	}
	return container.NewVScroll(items)
}

func (a *App) dashboardCard(p *domain.Project, d *domain.Dashboard) fyne.CanvasObject {
	name := widget.NewLabelWithStyle(d.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	created := widget.NewLabel(fmt.Sprintf("Criado em %s", d.CreatedAt.Format("02/01/2006 15:04")))

	open := widget.NewButtonWithIcon("Editar", theme.DocumentCreateIcon(), func() {
		a.showEditor(p, d)
	})
	open.Importance = widget.HighImportance

	rename := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		a.promptRenameDashboard(p, d)
	})
	del := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		a.confirmDeleteDashboard(p, d)
	})

	actions := container.NewHBox(rename, del, open)
	body := container.NewBorder(nil, nil, nil, actions, container.NewVBox(name, created))
	return container.NewPadded(widget.NewCard("", "", body))
}

func (a *App) promptNewDashboard(p *domain.Project) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Ex.: Fluxo de Checkout")

	form := dialog.NewForm("Novo dashboard", "Criar", "Cancelar",
		[]*widget.FormItem{{Text: "Nome", Widget: entry}},
		func(ok bool) {
			if !ok {
				return
			}
			if entry.Text == "" {
				dialog.ShowError(fmt.Errorf("informe um nome para o dashboard"), a.window)
				return
			}
			d, err := a.store.CreateDashboard(p.ID, entry.Text)
			if err != nil {
				dialog.ShowError(err, a.window)
				return
			}
			a.showEditor(p, d)
		},
		a.window,
	)
	form.Resize(fyne.NewSize(480, 200))
	form.Show()
}

func (a *App) promptRenameDashboard(p *domain.Project, d *domain.Dashboard) {
	entry := widget.NewEntry()
	entry.SetText(d.Name)

	form := dialog.NewForm("Renomear dashboard", "Salvar", "Cancelar",
		[]*widget.FormItem{{Text: "Nome", Widget: entry}},
		func(ok bool) {
			if !ok {
				return
			}
			if entry.Text == "" {
				dialog.ShowError(fmt.Errorf("informe um nome para o dashboard"), a.window)
				return
			}
			d.Name = entry.Text
			if err := a.store.UpdateDashboard(d); err != nil {
				dialog.ShowError(err, a.window)
				return
			}
			a.showDashboards(p)
		},
		a.window,
	)
	form.Resize(fyne.NewSize(480, 200))
	form.Show()
}

func (a *App) confirmDeleteDashboard(p *domain.Project, d *domain.Dashboard) {
	msg := fmt.Sprintf("Remover o dashboard \"%s\" e seus post-its?", d.Name)
	dialog.ShowConfirm("Remover dashboard", msg, func(ok bool) {
		if !ok {
			return
		}
		if err := a.store.DeleteDashboard(d.ID); err != nil {
			dialog.ShowError(err, a.window)
			return
		}
		a.showDashboards(p)
	}, a.window)
}
