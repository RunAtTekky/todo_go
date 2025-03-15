package main

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func setup_DB() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", "./todo.db")

	if err != nil {
		return nil, err
	}

	createTABLEsql := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		done BOOLEAN NOT NULL DEFAULT 0,
		details TEXT NOT NULL,
	);
	`

	_, err = db.Exec(createTABLEsql)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func initialModel(db *sql.DB) tasks {
	ti := textinput.New()
	ti.Placeholder = "Enter a task..."
	ti.Focus()

	return tasks{
		entries:    []task{},
		text_input: ti,
		input_mode: true,
	}
}

// INIT
func (m tasks) Init() tea.Cmd {
	return nil
}

// VIEW
func (m tasks) View() string {

	output := ""

	// Controls
	if m.input_mode {
		output += "Enter task: " + m.text_input.View() + "\n"
	}

	var task_option []string
	for i, t := range m.entries {

		if m.index == i {

			if t.done {
				task_option = append(task_option, fmt.Sprintf("-> [x] %s ", t.details))
			} else {
				task_option = append(task_option, fmt.Sprintf("-> [ ] %s ", t.details))
			}
		} else {
			if t.done {
				task_option = append(task_option, fmt.Sprintf("   [x] %s ", t.details))
			} else {
				task_option = append(task_option, fmt.Sprintf("   [ ] %s ", t.details))
			}
		}
	}
	if len(m.entries) == 0 {
		// output += "No tasks!\n"
		task_option = append(task_option, "No tasks!\n")
	}

	help_options := ""

	if m.show_help {

		help_options += `
Use vim bindings k/j to move up/down
Use ctrl+c or q to quit
Use c to create new task 
	Use esc to return to tasks
Use d to delete task
	`
	} else {
		help_options += `
Use ? to show help
		`
	}

	output += fmt.Sprintf(`
Hiiiiii, this is our TODO app
%s

%s

`, strings.Join(task_option, "\n"), help_options)

	return output
}

type toggle_casing_msg struct{}

func (m tasks) toggle_selected_item() tea.Model {
	// txt := m.entries[m.index].details

	if len(m.entries) == 0 {
		return m
	}

	if m.entries[m.index].done {
		m.entries[m.index].done = false
	} else {
		m.entries[m.index].done = true
	}

	return m
}

// UPDATE
func (m tasks) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case toggle_casing_msg:
		return m.toggle_selected_item(), nil

	case tea.KeyMsg:
		if m.input_mode {
			m.text_input, _ = m.text_input.Update(msg)
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				m.input_mode = false
				m.text_input.Reset()
				return m, tea.ClearScreen
			case "enter", "return":
				new_task := m.text_input.Value()
				if new_task != "" {
					m.entries = append(m.entries, task{
						details:  new_task,
						done:     false,
						on_press: func() tea.Msg { return toggle_casing_msg{} },
					})
				}
				m.input_mode = false
				m.text_input.Reset()

				return m, tea.ClearScreen

			}

		} else {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "enter", "return":
				return m, m.entries[m.index].on_press
			case "j":
				m.index++
				if m.index >= len(m.entries) {
					m.index = len(m.entries) - 1
				}
				return m, nil
			case "k":
				m.index--
				if m.index < 0 {
					m.index = 0
				}
				return m, nil
			case "c":
				m.input_mode = true
				m.text_input.Focus()
				return m, textinput.Blink
			case "d":
				if len(m.entries) == 0 {
					return m, nil
				}
				m.entries = append(m.entries[:m.index], m.entries[m.index+1:]...)
				return m, tea.ClearScreen
			case "?":
				m.show_help = !m.show_help
				return m, tea.ClearScreen
			}

		}
	}

	return m, nil
}
