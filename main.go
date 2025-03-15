package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	db, err := setup_DB()

	if err != nil {
		log.Fatalf("Failed to setup database %v", err)
	}

	defer db.Close()

	p := tea.NewProgram(initialModel(db))

	if _, err := p.Run(); err != nil {
		panic(err)
	}

}
