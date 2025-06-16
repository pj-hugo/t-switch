package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type theme string

type model struct {
	cursor  int
	themes  []theme
	chosen  *theme
}

func initialModel() model {
	return model{
		themes: []theme{
			"Kanagawa",
			"Tokyo Night",
			"Rose Pine",
			"Catppuccin",
		},
		chosen: nil,
	}
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("t-switch")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.themes)-1 {
				m.cursor++
			}
		case "enter":
			selectedTheme := m.themes[m.cursor]
			m.chosen = &selectedTheme
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "Choose a theme:\n\n"

	for i, t := range m.themes {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, t)
	}

	s += "\nPress 'enter' to select, 'q' to quit.\n"

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v\n", err)
		os.Exit(1)
	}

	if m, ok := finalModel.(model); ok && m.chosen != nil {
		fmt.Printf("You chose: %s\n", *m.chosen)
	} else if m.chosen == nil {
		fmt.Println("No theme was selected.")
	}
}
