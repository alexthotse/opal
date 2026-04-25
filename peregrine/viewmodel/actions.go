package viewmodel

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	actionTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		MarginBottom(1)

	actionItemStyle = lipgloss.NewStyle().
		PaddingLeft(2)

	actionSelectedItemStyle = lipgloss.NewStyle().
		PaddingLeft(2).
		Foreground(lipgloss.Color("39")).
		Bold(true)
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
			"Edit",
			"Copy",
			"Retry",
			"Explain",
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
	
	s := actionTitleStyle.Render("⚡ Inline Actions") + "\n"
	for i, action := range m.actions {
		if m.cursor == i {
			s += actionSelectedItemStyle.Render(fmt.Sprintf("▶ %s", action)) + "\n"
		} else {
			s += actionItemStyle.Render(fmt.Sprintf("  %s", action)) + "\n"
		}
	}
	
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1)
	s += helpStyle.Render("Press 'enter' to select, 'm' or 'esc' to close.")

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Padding(1, 2)
	return style.Render(s)
}
