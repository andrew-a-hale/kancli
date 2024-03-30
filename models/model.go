package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"kancli/database"
)

// Board Swimlane
type status int

// Model Management
type view int

var models []tea.Model

const (
	BOARD view = iota
	FORM
)

func NewModel(db database.Database) *[]tea.Model {
	models = []tea.Model{
		NewBoard(db),
		NewForm(db, TODO, 0, 0),
	}

	return &models
}
