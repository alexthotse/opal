package viewmodel

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/alexthotse/peregrine/adapters"
	"github.com/alexthotse/peregrine/domain"
	"github.com/alexthotse/peregrine/view"
)

type AppViewModel struct {
	state          domain.State
	backendClient  *adapters.BackendClient
	animatorClient *adapters.Animator
	agentClient    *adapters.ADKAgentClient

	voice   VoiceModel
	history HistoryModel
	actions ActionsModel
}

func NewAppViewModel(bc *adapters.BackendClient, ac *adapters.ADKAgentClient, anim *adapters.Animator) *AppViewModel {
	return &AppViewModel{
		state:          domain.InitialState(),
		backendClient:  bc,
		animatorClient: anim,
		agentClient:    ac,
		voice:          NewVoiceModel(),
		history:        NewHistoryModel(),
		actions:        NewActionsModel(),
	}
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type rpcResultMsg string

func (m *AppViewModel) Init() tea.Cmd {
	return tea.Batch(
		tick(),
		func() tea.Msg {
			// Initially ping the backend to ensure connectivity
			res, err := m.backendClient.Ping(context.Background(), "init-ping")
			if err != nil {
				return rpcResultMsg(fmt.Sprintf("RPC Error: %v", err))
			}
			return rpcResultMsg(fmt.Sprintf("Backend says: %s", res))
		},
	)
}

func (m *AppViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		// Let sub-models handle keys first if they are active
		if m.voice.active {
			m.voice, cmd = m.voice.Update(msg)
			return m, cmd
		}
		if m.history.active {
			m.history, cmd = m.history.Update(msg)
			return m, cmd
		}
		if m.actions.active {
			m.actions, cmd = m.actions.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "v":
			m.voice.active = true
			return m, nil
		case "h":
			m.history.active = true
			return m, nil
		case "m":
			m.actions.active = true
			return m, nil
		case "t":
			// Cycle themes
			newTheme := (m.state.ThemeType + 1) % 3
			m.state = domain.SetTheme(m.state, newTheme)
			return m, nil
		case "u":
			// Trigger Ultrathink via RPC
			m.state = domain.AddMessage(m.state, "Triggering ULTRATHINK...")
			m.animatorClient.SetTarget(20.0) // Bounce right
			return m, func() tea.Msg {
				res, err := m.backendClient.StartUltrathink(context.Background(), "u-req", "Test reasoning prompt")
				if err != nil {
					return rpcResultMsg(fmt.Sprintf("Ultrathink failed: %v", err))
				}
				return rpcResultMsg(fmt.Sprintf("Ultrathink result: %s", res))
			}
		case "j":
			// Trigger Jido Action via RPC
			m.state = domain.AddMessage(m.state, "Dispatching Jido Action...")
			return m, func() tea.Msg {
				res, err := m.backendClient.DispatchAction(context.Background(), "j-req", "test_action")
				if err != nil {
					return rpcResultMsg(fmt.Sprintf("Jido Action failed: %v", err))
				}
				return rpcResultMsg(fmt.Sprintf("Jido Action result: %s", res))
			}
		case "g":
			// Generate reasoning via Gemini
			m.state = domain.AddMessage(m.state, "Generating reasoning...")
			return m, func() tea.Msg {
				res, err := m.backendClient.GenerateReasoning("Explain the concept of Hexagonal Architecture in one short sentence.")
				if err != nil {
					return rpcResultMsg(fmt.Sprintf("Reasoning failed: %v", err))
				}
				return rpcResultMsg(fmt.Sprintf("Reasoning: %s", res))
			}
		case "space":
			m.animatorClient.SetTarget(0.0) // Bounce back
			m.state = domain.IncrementRequestID(m.state)
			m.state = domain.AddMessage(m.state, fmt.Sprintf("Ping %d", m.state.RequestID))
			return m, nil
		}

	case rpcResultMsg:
		m.state = domain.AddMessage(m.state, string(msg))
		m.animatorClient.SetTarget(0.0) // Reset position
		return m, nil

	case tickMsg:
		m.animatorClient.Update()
		return m, tick()
	}

	return m, nil
}

func (m *AppViewModel) View() string {
	baseView := view.Render(m.state, m.animatorClient.Position())
	
	if m.voice.active {
		baseView += "\n\n" + m.voice.View()
	}
	if m.history.active {
		baseView += "\n\n" + m.history.View()
	}
	if m.actions.active {
		baseView += "\n\n" + m.actions.View()
	}

	return baseView
}
