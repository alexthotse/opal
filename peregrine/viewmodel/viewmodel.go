package viewmodel

import (
	"io"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"peregrine/adapters"
	"peregrine/domain"
	"peregrine/ports"
	"peregrine/view"
)

type AppViewModel struct {
	State          domain.State
	BackendCmd     *exec.Cmd
	Stdin          io.WriteCloser
	Stdout         io.ReadCloser
	BackendClient  ports.BackendClient
	AgentClient    ports.AgentClient
	AnimatorClient ports.AnimatorClient
}

type TickMsg time.Time

func TickCmd() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func NewAppViewModel(bc ports.BackendClient, ac ports.AgentClient, anim ports.AnimatorClient) AppViewModel {
	anim.SetTarget(20.0)
	return AppViewModel{
		State:          domain.InitialState(),
		BackendClient:  bc,
		AgentClient:    ac,
		AnimatorClient: anim,
	}
}

func (m AppViewModel) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return m.BackendClient.Start() },
		TickCmd(),
	)
}

func (m AppViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.BackendCmd != nil {
				m.BackendCmd.Process.Kill()
			}
			return m, tea.Quit
		case "p":
			if m.Stdin != nil {
				m.State = domain.IncrementRequestID(m.State)
				return m, m.BackendClient.SendRequest(m.Stdin, m.State.RequestID, "ping")
			}
		case "a":
			if m.Stdin != nil {
				out := m.AgentClient.ProcessPrompt("trigger backend action")
				m.State = domain.AddMessage(m.State, out)
				m.State = domain.IncrementRequestID(m.State)
				return m, m.BackendClient.SendRequest(m.Stdin, m.State.RequestID, "jido.action")
			}
		case "u":
			if m.Stdin != nil {
				m.State = domain.IncrementRequestID(m.State)
				return m, m.BackendClient.SendRequest(m.Stdin, m.State.RequestID, "ultrathink.start")
			}
		case "r":
			if m.Stdin != nil {
				m.State = domain.IncrementRequestID(m.State)
				return m, m.BackendClient.SendRequest(m.Stdin, m.State.RequestID, "ultraplan.start")
			}
		case "v":
			m.State = domain.AddMessage(m.State, "VOICE_MODE: Listening for voice input...")
		case "h":
			m.State = domain.AddMessage(m.State, "HISTORY_PICKER: Displaying prompt history...")
		case "m":
			m.State = domain.AddMessage(m.State, "MESSAGE_ACTIONS: Displaying message actions...")
		case "s":
			if m.Stdin != nil {
				out := m.AgentClient.ProcessPrompt("search logs")
				m.State = domain.AddMessage(m.State, out)
				m.State = domain.IncrementRequestID(m.State)
				return m, m.BackendClient.SendRequest(m.Stdin, m.State.RequestID, "search.quick")
			}
		case "t":
			if m.State.ThemeType == "pi.dev" {
				m.State = domain.SetTheme(m.State, "free-code")
			} else if m.State.ThemeType == "free-code" {
				m.State = domain.SetTheme(m.State, "crush")
			} else {
				m.State = domain.SetTheme(m.State, "pi.dev")
			}
		}

	case TickMsg:
		m.AnimatorClient.Update()
		return m, TickCmd()

	case adapters.BackendStartedMsg:
		m.BackendCmd = msg.Cmd
		m.Stdin = msg.Stdin
		m.Stdout = msg.Stdout
		m.State = domain.AddMessage(m.State, "Falcon Backend connected! Keys: 'p' (ping), 'a' (agent), 'u' (think), 'r' (plan), 'v' (voice), 'h' (history), 'm' (actions), 's' (search), 't' (theme), 'q' (quit).")
		return m, m.BackendClient.Listen(m.Stdout)

	case adapters.BackendMsg:
		parsed := strings.TrimSpace(msg.Content)
		if parsed != "" {
			m.State = domain.AddMessage(m.State, "Gleam: "+parsed)
		}
		if m.Stdout != nil && !strings.Contains(parsed, "disconnected") {
			return m, m.BackendClient.Listen(m.Stdout)
		}
	}

	return m, nil
}

func (m AppViewModel) View() string {
	return view.Render(m.State, m.AnimatorClient.Position())
}
