package main

import tea "github.com/charmbracelet/bubbletea"

type tasks struct {
	entries []task
	index   int
}

type task struct {
	done     bool
	details  string
	on_press func() tea.Msg
}
