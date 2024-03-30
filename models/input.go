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

type Form struct {
	db          database.Database
	focused     status
	title       textinput.Model
	description textarea.Model
	width       int
	height      int
}

func NewForm(db database.Database, focused status, width int, height int) *Form {
	form := &Form{db: db, focused: focused, width: width, height: height}

	form.title = textinput.New()
	form.title.Placeholder = "New Task Name"
	form.title.Focus()

	form.description = textarea.New()
	form.description.Placeholder = "A large and important module!"
	return form
}

func (m Form) NewTask() tea.Msg {
	task := NewTask(m.focused, m.title.Value(), m.description.Value())
	return task
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
				return models[BOARD], m.NewTask
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
