package viewmodel

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HistoryModel struct {
	active  bool
	prompts []string
	cursor  int
}

func NewHistoryModel() HistoryModel {
	return HistoryModel{
		active: false,
		prompts: []string{
			"Show me the architecture plan.",
			"Run the benchmark suite.",
			"Explain the Gleam concurrency model.",
		},
		cursor: 0,
	}
}

func (m HistoryModel) Init() tea.Cmd {
	return nil
}

func (m HistoryModel) Update(msg tea.Msg) (HistoryModel, tea.Cmd) {
	if !m.active {
		return m, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "h":
			m.active = false
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.prompts)-1 {
				m.cursor++
			}
		case "enter":
			// For now, just exit history mode when a prompt is selected.
			m.active = false
		}
	}
	return m, nil
}

func (m HistoryModel) View() string {
	if !m.active {
		return ""
	}
	
	s := "📜 Prompt History\n\n"
	for i, prompt := range m.prompts {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, prompt)
	}
	s += "\nPress 'h' or 'esc' to close."

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2)
	return style.Render(s)
}
