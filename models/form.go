package models

import (
	"kancli/database"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var inputStyle = lipgloss.NewStyle().
	Padding(1, 2).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("62")).
	Foreground(lipgloss.Color("241"))

type mode int

const (
	NEW mode = iota
	EDIT
)

type Form struct {
	db          database.Database
	name        string
	focused     status
	id          int
	title       textinput.Model
	description textarea.Model
	width       int
	height      int
	mode        mode
}

func NewForm(db database.Database, focused status, width int, height int) *Form {
	form := &Form{
		db:      db,
		focused: focused,
		width:   width,
		height:  height,
		name:    "Create New Task",
		mode:    NEW,
	}

	form.title = textinput.New()
	form.title.Placeholder = "New Task Name"
	form.title.Focus()

	form.description = textarea.New()
	form.description.Placeholder = "A large and important module!"
	return form
}

func (m Form) NewTask() tea.Msg {
	return NewTask(m.focused, m.title.Value(), m.description.Value())
}

func NewFormWithTask(db database.Database, focused status, width int, height int, task Task) *Form {
	form := &Form{
		db:      db,
		focused: focused,
		id:      task.id,
		width:   width,
		height:  height,
		name:    "Edit Task",
		mode:    EDIT,
	}

	form.title = textinput.New()
	form.title.SetValue(task.title)
	form.title.Focus()

	form.description = textarea.New()
	form.description.SetValue(task.description)
	return form
}

func (m Form) EditTask() tea.Msg {
	return EditTask(m.id, m.focused, m.title.Value(), m.description.Value())
}

func (m Form) Init() tea.Cmd {
	return nil
}

func (m Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width-2, msg.Height-5
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.title.Focused() {
				m.title.Blur()
				m.description.Focus()
				return m, textarea.Blink
			} else {
				models[FORM] = m
				switch m.mode {
				case NEW:
					return models[BOARD], m.NewTask
				case EDIT:
					return models[BOARD], m.EditTask
				}
			}
		}
	}

	var cmd tea.Cmd
	if m.title.Focused() {
		m.title, cmd = m.title.Update(msg)
		return m, cmd
	} else {
		m.description, cmd = m.description.Update(msg)
		return m, cmd
	}
}

func (m Form) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.Place(
			m.width,
			m.height/2-1,
			lipgloss.Center,
			lipgloss.Center,
			m.name,
		),
		lipgloss.Place(
			m.width,
			m.height/2-1,
			lipgloss.Center,
			lipgloss.Center,
			inputStyle.Render(m.title.View()),
		),
		lipgloss.Place(
			m.width,
			m.height/2-1,
			lipgloss.Center,
			lipgloss.Center,
			inputStyle.Render(m.description.View()),
		),
	)
}
