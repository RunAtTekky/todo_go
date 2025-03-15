package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	// m := tasks{
	// 	entries: []task{
	// 		{
	// 			done:    false,
	// 			details: "I'm a task",
	// 			// on_press: ,
	// 			on_press: func() tea.Msg { return toggle_casing_msg{} },
	// 		},
	// 		{
	// 			done:     false,
	// 			details:  "I'm another task",
	// 			on_press: func() tea.Msg { return toggle_casing_msg{} },
	// 		},
	// 	},
	// 	index: 0,
	// }

	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		panic(err)
	}

}
