package viewmodel

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type VoiceModel struct {
	active bool
}

func NewVoiceModel() VoiceModel {
	return VoiceModel{active: false}
}

func (m VoiceModel) Init() tea.Cmd {
	return nil
}

func (m VoiceModel) Update(msg tea.Msg) (VoiceModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "v":
			m.active = false
			return m, nil
		}
	}
	return m, nil
}

func (m VoiceModel) View() string {
	if !m.active {
		return ""
	}
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2)
	return style.Render("🎤 Voice Mode (Push-to-talk)\n\nPress 'v' or 'esc' to exit.")
}
