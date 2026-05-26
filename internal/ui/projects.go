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

func (a *App) showProjects() {
	projects, err := a.store.ListProjects()
	if err != nil {
		dialog.ShowError(fmt.Errorf("carregando projetos: %w", err), a.window)
		return
	}

	title := widget.NewLabelWithStyle("Projetos", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	subtitle := widget.NewLabel("Selecione um projeto para abrir seus dashboards ou crie um novo.")
	subtitle.Wrapping = fyne.TextWrapWord

	newBtn := widget.NewButtonWithIcon("Novo projeto", theme.ContentAddIcon(), func() {
		a.promptNewProject()
	})
	newBtn.Importance = widget.HighImportance

	list := a.buildProjectsList(projects)

	header := container.NewBorder(nil, nil, nil, newBtn, title)
	content := container.NewBorder(
		container.NewVBox(header, subtitle, widget.NewSeparator()),
		nil, nil, nil,
		list,
	)
	a.swap(container.NewPadded(content))
}

func (a *App) buildProjectsList(projects []*domain.Project) fyne.CanvasObject {
	if len(projects) == 0 {
		empty := widget.NewLabel("Nenhum projeto ainda. Clique em \"Novo projeto\" para começar.")
		empty.Alignment = fyne.TextAlignCenter
		return container.NewCenter(empty)
	}

	items := container.NewVBox()
	for _, p := range projects {
		items.Add(a.projectCard(p))
	}
	return container.NewVScroll(items)
}

func (a *App) projectCard(p *domain.Project) fyne.CanvasObject {
	name := widget.NewLabelWithStyle(p.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	desc := widget.NewLabel(p.Description)
	desc.Wrapping = fyne.TextWrapWord

	open := widget.NewButtonWithIcon("Abrir", theme.NavigateNextIcon(), func() {
		a.showDashboards(p)
	})
	open.Importance = widget.HighImportance

	edit := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		a.promptEditProject(p)
	})
	del := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		a.confirmDeleteProject(p)
	})

	actions := container.NewHBox(edit, del, open)
	body := container.NewBorder(nil, nil, nil, actions, container.NewVBox(name, desc))

	return container.NewPadded(widget.NewCard("", "", body))
}

func (a *App) promptNewProject() {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Ex.: Plataforma de Pagamentos")
	descEntry := widget.NewMultiLineEntry()
	descEntry.SetPlaceHolder("Contexto, objetivos, escopo...")
	descEntry.Wrapping = fyne.TextWrapWord

	form := dialog.NewForm("Novo projeto", "Criar", "Cancelar",
		[]*widget.FormItem{
			{Text: "Nome", Widget: nameEntry},
			{Text: "Descrição", Widget: descEntry},
		},
		func(ok bool) {
			if !ok {
				return
			}
			if nameEntry.Text == "" {
				dialog.ShowError(fmt.Errorf("informe um nome para o projeto"), a.window)
				return
			}
			if _, err := a.store.CreateProject(nameEntry.Text, descEntry.Text); err != nil {
				dialog.ShowError(err, a.window)
				return
			}
			a.showProjects()
		},
		a.window,
	)
	form.Resize(fyne.NewSize(520, 320))
	form.Show()
}

func (a *App) promptEditProject(p *domain.Project) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(p.Name)
	descEntry := widget.NewMultiLineEntry()
	descEntry.SetText(p.Description)
	descEntry.Wrapping = fyne.TextWrapWord

	form := dialog.NewForm("Editar projeto", "Salvar", "Cancelar",
		[]*widget.FormItem{
			{Text: "Nome", Widget: nameEntry},
			{Text: "Descrição", Widget: descEntry},
		},
		func(ok bool) {
			if !ok {
				return
			}
			if nameEntry.Text == "" {
				dialog.ShowError(fmt.Errorf("informe um nome para o projeto"), a.window)
				return
			}
			p.Name = nameEntry.Text
			p.Description = descEntry.Text
			if err := a.store.UpdateProject(p); err != nil {
				dialog.ShowError(err, a.window)
				return
			}
			a.showProjects()
		},
		a.window,
	)
	form.Resize(fyne.NewSize(520, 320))
	form.Show()
}

func (a *App) confirmDeleteProject(p *domain.Project) {
	msg := fmt.Sprintf("Remover o projeto \"%s\" e todos os seus dashboards?", p.Name)
	dialog.ShowConfirm("Remover projeto", msg, func(ok bool) {
		if !ok {
			return
		}
		if err := a.store.DeleteProject(p.ID); err != nil {
			dialog.ShowError(err, a.window)
			return
		}
		a.showProjects()
	}, a.window)
}
