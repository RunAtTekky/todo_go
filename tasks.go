package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type tasks struct {
	entries []task
	index   int

	text_input textinput.Model
	input_mode bool
	show_help  bool

	db *sql.DB
}

type task struct {
	id       int
	done     bool
	details  string
	on_press func() tea.Msg
}
