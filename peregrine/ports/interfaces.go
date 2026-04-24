package ports

import (
	"io"
	tea "github.com/charmbracelet/bubbletea"
)

// BackendClient represents the side-effectful boundary to our Gleam backend.
type BackendClient interface {
	Start() tea.Msg
	Listen(stdout io.ReadCloser) tea.Cmd
	SendRequest(stdin io.WriteCloser, reqID int, method string) tea.Cmd
}

// AgentClient represents the side-effectful boundary to our ADK agent.
type AgentClient interface {
	ProcessPrompt(prompt string) string
}

// AnimatorClient represents the UI animation engine boundary.
type AnimatorClient interface {
	Update()
	SetTarget(target float64)
	Position() float64
}
