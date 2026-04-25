package viewmodel

import (
	"log"

	"github.com/alexthotse/peregrine/adapters"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AudioRecorder interface {
	Start() error
	Stop() error
	Read(p []int16) (int, error)
	Close() error
}

type MockAudioRecorder struct{}

func (m *MockAudioRecorder) Start() error {
	log.Println("MockAudioRecorder: Start recording (hardware absent)")
	return nil
}

func (m *MockAudioRecorder) Stop() error {
	log.Println("MockAudioRecorder: Stop recording")
	return nil
}

func (m *MockAudioRecorder) Read(p []int16) (int, error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

func (m *MockAudioRecorder) Close() error {
	return nil
}

type VoiceModel struct {
	active   bool
	recorder AudioRecorder
	tts      adapters.TTSProvider
}

func NewVoiceModel() VoiceModel {
	return VoiceModel{
		active:   false,
		recorder: &MockAudioRecorder{},
		tts:      adapters.NewTTSProvider("pocket-tts"),
	}
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
