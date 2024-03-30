package main

import (
	"log"

	"kancli/database"
	"kancli/models"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	db, err := database.Init("./data.db")
	if err != nil {
		log.Fatalf("unable to open database: %v", err)
	}

	model := *models.NewModel(db)
	m := model[models.BOARD]
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
