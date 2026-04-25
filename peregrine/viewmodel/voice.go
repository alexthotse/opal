package viewmodel

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type VoiceModel struct {
	active    bool
	recording bool
	err       error
}

func NewVoiceModel() VoiceModel {
	return VoiceModel{active: false, recording: false}
}

func (m VoiceModel) Init() tea.Cmd {
	return nil
}

type recordFinishedMsg struct{ err error }

func recordAudio() tea.Msg {
	// Dummy OS command to simulate audio capture.
	cmd := exec.Command("arecord", "-d", "3", "/dev/null")
	err := cmd.Run()
	return recordFinishedMsg{err: err}
}

func (m VoiceModel) Update(msg tea.Msg) (VoiceModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "v":
			m.active = false
			return m, nil
		case "r":
			if m.active && !m.recording {
				m.recording = true
				return m, recordAudio
			}
		}
	case recordFinishedMsg:
		m.recording = false
		m.err = msg.err
	}
	return m, nil
}

func (m VoiceModel) View() string {
	if !m.active {
		return ""
	}

	status := "Press 'r' to start recording."
	if m.recording {
		status = "Recording... (simulating with arecord for 3s)"
	} else if m.err != nil {
		status = "Recording finished (or failed: " + m.err.Error() + ")"
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2)
	return style.Render("🎤 Voice Mode\n\n" + status + "\n\nPress 'v' or 'esc' to exit.")
}
