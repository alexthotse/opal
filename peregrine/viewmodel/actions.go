package viewmodel

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ActionsModel struct {
	active  bool
	actions []string
	cursor  int
}

func NewActionsModel() ActionsModel {
	return ActionsModel{
		active: false,
		actions: []string{
			"Copy to clipboard",
			"Save to file",
			"Regenerate response",
			"Clear chat",
		},
		cursor: 0,
	}
}

func (m ActionsModel) Init() tea.Cmd {
	return nil
}

func (m ActionsModel) Update(msg tea.Msg) (ActionsModel, tea.Cmd) {
	if !m.active {
		return m, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "m":
			m.active = false
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.actions)-1 {
				m.cursor++
			}
		case "enter":
			// In a real app, dispatch an action
			m.active = false
		}
	}
	return m, nil
}

func (m ActionsModel) View() string {
	if !m.active {
		return ""
	}
	
	s := "⚡ Inline Actions\n\n"
	for i, action := range m.actions {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, action)
	}
	s += "\nPress 'm' or 'esc' to close."

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Padding(1, 2)
	return style.Render(s)
}
