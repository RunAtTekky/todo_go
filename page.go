package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"

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
		details TEXT NOT NULL
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

	model := tasks{
		entries:    []task{},
		text_input: ti,
		input_mode: false,
		db:         db,
	}

	model.load_DB()

	return model
}

func (m *tasks) load_DB() {

	rows, err := m.db.Query("SELECT id, done, details FROM tasks")

	if err != nil {
		log.Printf("Error loading tasks %v", err)
		return
	}

	defer rows.Close()

	var loaded_tasks []task

	for rows.Next() {
		var t task

		err := rows.Scan(&t.id, &t.done, &t.details)

		if err != nil {
			log.Printf("Error loading this row %v", err)
			continue
		}

		t.on_press = func() tea.Msg { return toggle_casing_msg{} }

		loaded_tasks = append(loaded_tasks, t)
	}

	m.entries = loaded_tasks

}

// ADD task
func (m *tasks) add_task(details string) {

	result, err := m.db.Exec("INSERT INTO tasks (done, details) VALUES (?, ?)", false, details)

	if err != nil {
		log.Printf("Error inserting task: %v", err)
		return
	}

	id, err := result.LastInsertId()

	if err != nil {
		log.Printf("Error getting last id: %v", err)
		return
	}

	m.entries = append(m.entries, task{
		id:       int(id),
		done:     false,
		details:  details,
		on_press: func() tea.Msg { return toggle_casing_msg{} },
	})
}

// UPDATE task
func (m *tasks) update_task(id int, done bool) {

	_, err := m.db.Exec("UPDATE tasks SET done = ? WHERE id = ?", done, id)

	if err != nil {
		log.Printf("Error updating task: %v", err)
	}

}

// DELETE task
func (m *tasks) delete_task(id int) {
	_, err := m.db.Exec("DELETE FROM tasks WHERE ID = ?", id)

	if err != nil {
		log.Printf("Error deleting task: %v", err)
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

	m.entries[m.index].done = !m.entries[m.index].done

	m.update_task(m.entries[m.index].id, m.entries[m.index].done)

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
					m.add_task(new_task)
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
				if len(m.entries) >= 0 {
					return m, m.entries[m.index].on_press
				}
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
				m.delete_task(m.entries[m.index].id)
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
