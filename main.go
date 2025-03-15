package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	m := tasks{
		entries: []task{
			{
				done:    false,
				details: "I'm a task",
			},
			{
				done:    false,
				details: "I'm another task",
			},
		},
		index: 0,
	}

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		panic(err)
	}

}
