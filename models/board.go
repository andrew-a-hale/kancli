package models

import (
	"kancli/database"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	columnStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.HiddenBorder())
	focusedStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

// CUSTOM ITEM TASK
type Task struct {
	status      status
	title       string
	description string
}

func NewTask(status status, title, description string) Task {
	return Task{title: title, description: description, status: status}
}

// implement the list.Item interface
func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

func (t *Task) Next() {
	if t.status < DONE {
		t.status++
	}
}

func (t *Task) Prev() {
	if t.status > TODO {
		t.status--
	}
}

// MAIN MODEL
type Board struct {
	db      database.Database
	focused status
	lists   []list.Model
	loaded  bool
	width   int
	height  int
}

func NewBoard(db database.Database) *Board {
	return &Board{db: db}
}

func (m *Board) InitLists() {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), m.width/DIVISOR, m.height)
	defaultList.SetShowHelp(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}

	m.lists[TODO].Title = "To do"
	m.lists[TODO].SetItems([]list.Item{})

	m.lists[IN_PROGRESS].Title = "In Progress"
	m.lists[IN_PROGRESS].SetItems([]list.Item{})

	m.lists[DONE].Title = "Done"
	m.lists[DONE].SetItems([]list.Item{})
}

func (m Board) Init() tea.Cmd {
	return nil
}

func (m Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width-5, msg.Height-5

		columnStyle.Width(m.width / DIVISOR)
		columnStyle.Height(m.height)

		focusedStyle.Width(m.width / DIVISOR)
		focusedStyle.Height(m.height)
		m.InitLists()
		m.loaded = true
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "l", "right":
			m.NextBoard()
		case "h", "left":
			m.PrevBoard()
		case "enter":
			m.MoveTask(1)
		case "backspace":
			m.MoveTask(-1)
		case "d":
			m.DeleteTask()
		case "n":
			models[BOARD] = m
			models[FORM] = NewForm(m.db, m.focused, m.width, m.height)
			return models[FORM].Update(nil)
		}
	case Task:
		task := msg
		return m, m.lists[task.status].InsertItem(0, task)
	}

	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m *Board) NextBoard() {
	m.focused++
	if m.focused > DONE {
		m.focused = TODO
	}
}

func (m *Board) PrevBoard() {
	m.focused--
	if m.focused < TODO {
		m.focused = DONE
	}
}

func (m *Board) MoveTask(dir int) tea.Msg {
	selectedIndex := m.lists[m.focused].Index()
	selectedTask, ok := m.lists[m.focused].SelectedItem().(Task)

	if ok {
		if dir > 0 {
			selectedTask.Next()
		} else {
			selectedTask.Prev()
		}

		m.lists[m.focused].RemoveItem(selectedIndex)
		m.lists[selectedTask.status].InsertItem(0, selectedTask)
	}

	return nil
}

func (m *Board) DeleteTask() tea.Msg {
	selectedIndex := m.lists[m.focused].Index()
	m.lists[m.focused].RemoveItem(selectedIndex)
	return nil
}

func (m Board) View() string {
	if !m.loaded {
		return ""
	}

	todoView := columnStyle.Render(m.lists[TODO].View())
	inProgressView := columnStyle.Render(m.lists[IN_PROGRESS].View())
	doneView := columnStyle.Render(m.lists[DONE].View())

	switch m.focused {
	case TODO:
		todoView = focusedStyle.Render(m.lists[TODO].View())
	case IN_PROGRESS:
		inProgressView = focusedStyle.Render(m.lists[IN_PROGRESS].View())
	case DONE:
		doneView = focusedStyle.Render(m.lists[DONE].View())
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		todoView,
		inProgressView,
		doneView,
	)
}
