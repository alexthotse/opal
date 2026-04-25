package viewmodel

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	historyTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(1)

	historyItemStyle = lipgloss.NewStyle().
		PaddingLeft(2)

	historySelectedItemStyle = lipgloss.NewStyle().
		PaddingLeft(2).
		Foreground(lipgloss.Color("205")).
		Bold(true)
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
			"What are the best practices for Go project structure?",
			"How do I use bazel to build this?",
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
	
	s := historyTitleStyle.Render("📜 Prompt History") + "\n"
	for i, prompt := range m.prompts {
		if m.cursor == i {
			s += historySelectedItemStyle.Render(fmt.Sprintf("▶ %s", prompt)) + "\n"
		} else {
			s += historyItemStyle.Render(fmt.Sprintf("  %s", prompt)) + "\n"
		}
	}
	
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1)
	s += helpStyle.Render("Press 'enter' to select, 'h' or 'esc' to close.")

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2)
	return style.Render(s)
}
