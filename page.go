package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// INIT
func (m tasks) Init() tea.Cmd {
	return nil
}

// VIEW
func (m tasks) View() string {
	var task_option []string
	for i, t := range m.entries {

		if m.index == i {
			task_option = append(task_option, fmt.Sprintf("-> %s", t.details))
		} else {
			task_option = append(task_option, fmt.Sprintf("   %s", t.details))
		}
	}

	output := fmt.Sprintf(`
Hiiiiii, this is our TODO app
%s

Use ctrl+c to quit
`, strings.Join(task_option, "\n"))

	return output

	// return strings.Join(task_option, "\n")
	// return m.entries[m.index].details
	// return "Hola"
}

// UPDATE
func (m tasks) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}
