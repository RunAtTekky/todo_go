package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type tasks struct {
	entries []task
	index   int

	text_input textinput.Model
	input_mode bool
}

type task struct {
	done     bool
	details  string
	on_press func() tea.Msg
}
