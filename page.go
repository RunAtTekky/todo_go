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
			task_option = append(task_option, fmt.Sprintf("-> %t \t| %s ", t.done, t.details))
		} else {
			task_option = append(task_option, fmt.Sprintf("   %t \t| %s ", t.done, t.details))
		}
	}

	output := fmt.Sprintf(`
Hiiiiii, this is our TODO app
%s

Use ctrl+c or q to quit
`, strings.Join(task_option, "\n"))

	return output

	// return strings.Join(task_option, "\n")
	// return m.entries[m.index].details
	// return "Hola"
}

type toggle_casing_msg struct{}

func (m tasks) toggle_selected_item() tea.Model {
	// txt := m.entries[m.index].details

	if m.entries[m.index].done {
		m.entries[m.index].done = false
	} else {
		m.entries[m.index].done = true
	}

	// if txt == strings.ToUpper(txt) {
	// 	m.entries[m.index].details = strings.ToLower(txt)
	// } else {
	// 	m.entries[m.index].details = strings.ToUpper(txt)
	// }

	return m
}

// UPDATE
func (m tasks) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case toggle_casing_msg:
		return m.toggle_selected_item(), nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter", "return":
			return m, m.entries[m.index].on_press
		case "j":
			m.index++
			return m, nil
		case "k":
			m.index--
			return m, nil
		}
	}

	return m, nil
}
