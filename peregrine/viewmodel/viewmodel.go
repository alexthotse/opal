package viewmodel

import (
        "context"
        "fmt"
        "strings"
        "time"

        "github.com/charmbracelet/bubbles/textinput"
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

        textInput textinput.Model
        voice     VoiceModel
        history   HistoryModel
        actions   ActionsModel
}

func NewAppViewModel(bc *adapters.BackendClient, ac *adapters.ADKAgentClient, anim *adapters.Animator, themeStr string) *AppViewModel {
        ti := textinput.New()
        ti.Placeholder = "Type a command (/budget, /teammem, /bridge, etc)..."
        ti.Focus()
        ti.CharLimit = 156
        ti.Width = 40

        initialState := domain.InitialState()
        switch themeStr {
        case "pi":
                initialState = domain.SetTheme(initialState, 0) // ThemePIDev
        case "freecode":
                initialState = domain.SetTheme(initialState, 1) // ThemeFreeCode
        case "crush":
                initialState = domain.SetTheme(initialState, 2) // ThemeCrush
        }

        return &AppViewModel{
                state:          initialState,
                backendClient:  bc,
                animatorClient: anim,
                agentClient:    ac,
                textInput:      ti,
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
                textinput.Blink,
                tick(),
                func() tea.Msg {
                        res, err := m.backendClient.Ping(context.Background(), "init-ping")
                        if err != nil {
                                return rpcResultMsg(fmt.Sprintf("RPC Error: %v", err))
                        }
                        return rpcResultMsg(fmt.Sprintf("Backend says: %s", res))
                },
        )
}

func (m *AppViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        var cmds []tea.Cmd
        var cmd tea.Cmd

        switch msg := msg.(type) {
        case tea.KeyMsg:
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

                switch msg.Type {
                case tea.KeyCtrlC, tea.KeyEsc:
                        return m, tea.Quit
                case tea.KeyEnter:
                        val := m.textInput.Value()
                        m.textInput.SetValue("")
                        if val != "" {
                                return m.handleCommand(val)
                        }
                default:
                        m.textInput, cmd = m.textInput.Update(msg)
                        cmds = append(cmds, cmd)
                }

                // Handle global shortcuts if they are not eating input
                // Or maybe don't if textInput is focused.
                // Actually, if we just type, it goes to textInput.

        case rpcResultMsg:
                m.state = domain.AddMessage(m.state, string(msg))
                m.animatorClient.SetTarget(0.0)
                return m, nil

        case tickMsg:
                m.animatorClient.Update()
                cmds = append(cmds, tick())
        }

        return m, tea.Batch(cmds...)
}

func (m *AppViewModel) handleCommand(val string) (tea.Model, tea.Cmd) {
        m.state = domain.AddMessage(m.state, "> "+val)
        m.animatorClient.SetTarget(20.0)

        id := "cmd-" + fmt.Sprint(m.state.RequestID)
        m.state = domain.IncrementRequestID(m.state)

        var cmd tea.Cmd
        ctx := context.Background()

        switch {
        case strings.HasPrefix(val, "/budget"):
                cmd = func() tea.Msg {
                        res, err := m.backendClient.GetTokenBudget(ctx, id)
                        if err != nil {
                                return rpcResultMsg("Error: " + err.Error())
                        }
                        return rpcResultMsg(res)
                }
        case strings.HasPrefix(val, "/teammem"):
                cmd = func() tea.Msg {
                        res, err := m.backendClient.GetTeamMem(ctx, id)
                        if err != nil {
                                return rpcResultMsg("Error: " + err.Error())
                        }
                        return rpcResultMsg(res)
                }
        case strings.HasPrefix(val, "/bridge"):
                cmd = func() tea.Msg {
                        res, err := m.backendClient.StartBridgeMode(ctx, id)
                        if err != nil {
                                return rpcResultMsg("Error: " + err.Error())
                        }
                        return rpcResultMsg(res)
                }
        case strings.HasPrefix(val, "/ultraplan"):
                cmd = func() tea.Msg {
                        res, err := m.backendClient.StartUltraplan(ctx, id)
                        if err != nil {
                                return rpcResultMsg("Error: " + err.Error())
                        }
                        return rpcResultMsg(res)
                }
        case strings.HasPrefix(val, "/ultrathink"):
                cmd = func() tea.Msg {
                        res, err := m.backendClient.StartUltrathink(ctx, id)
                        if err != nil {
                                return rpcResultMsg("Error: " + err.Error())
                        }
                        return rpcResultMsg(res)
                }
        case strings.HasPrefix(val, "/triggerremote"):
                cmd = func() tea.Msg {
                        res, err := m.backendClient.SetAgentTriggerRemote(ctx, id)
                        if err != nil {
                                return rpcResultMsg("Error: " + err.Error())
                        }
                        return rpcResultMsg(res)
                }
        case strings.HasPrefix(val, "/trigger"):
                cmd = func() tea.Msg {
                        res, err := m.backendClient.SetAgentTrigger(ctx, id)
                        if err != nil {
                                return rpcResultMsg("Error: " + err.Error())
                        }
                        return rpcResultMsg(res)
                }
        case strings.HasPrefix(val, "/extract"):
                cmd = func() tea.Msg {
                        res, err := m.backendClient.ExtractMemories(ctx, id)
                        if err != nil {
                                return rpcResultMsg("Error: " + err.Error())
                        }
                        return rpcResultMsg(res)
                }
        case strings.HasPrefix(val, "/microcompact"):
                cmd = func() tea.Msg {
                        res, err := m.backendClient.CachedMicrocompact(ctx, id)
                        if err != nil {
                                return rpcResultMsg("Error: " + err.Error())
                        }
                        return rpcResultMsg(res)
                }
        case strings.HasPrefix(val, "/verify"):
                cmd = func() tea.Msg {
                        res, err := m.backendClient.StartVerification(ctx, id)
                        if err != nil {
                                return rpcResultMsg("Error: " + err.Error())
                        }
                        return rpcResultMsg(res)
                }
        case strings.HasPrefix(val, "/stats"):
                cmd = func() tea.Msg {
                        res, err := m.backendClient.GetShotStats(ctx, id)
                        if err != nil {
                                return rpcResultMsg("Error: " + err.Error())
                        }
                        return rpcResultMsg(res)
                }
        case strings.HasPrefix(val, "/classify"):
                cmd = func() tea.Msg {
                        res, err := m.backendClient.ClassifyBash(ctx, id)
                        if err != nil {
                                return rpcResultMsg("Error: " + err.Error())
                        }
                        return rpcResultMsg(res)
                }
        case strings.HasPrefix(val, "/compact"):
                cmd = func() tea.Msg {
                        res, err := m.backendClient.GetCompactionReminders(ctx, id)
                        if err != nil {
                                return rpcResultMsg("Error: " + err.Error())
                        }
                        return rpcResultMsg(res)
                }
        case strings.HasPrefix(val, "/search"):
                parts := strings.SplitN(val, " ", 2)
                q := ""
                if len(parts) > 1 {
                        q = parts[1]
                }
                cmd = func() tea.Msg {
                        res, err := m.backendClient.QuickSearch(ctx, id, q)
                        if err != nil {
                                return rpcResultMsg("Error: " + err.Error())
                        }
                        return rpcResultMsg(res)
                }
        case val == "v":
                m.voice.active = true
                return m, nil
        case val == "h":
                m.history.active = true
                return m, nil
        case val == "m":
                m.actions.active = true
                return m, nil
        case val == "t":
                newTheme := (m.state.ThemeType + 1) % 3
                m.state = domain.SetTheme(m.state, newTheme)
                return m, nil
        case val == "j":
                cmd = func() tea.Msg {
                        res, err := m.backendClient.DispatchAction(ctx, id, "test_action")
                        if err != nil {
                                return rpcResultMsg("Error: " + err.Error())
                        }
                        return rpcResultMsg(res)
                }
        default:
                m.animatorClient.SetTarget(0.0)
                cmd = func() tea.Msg {
                        return rpcResultMsg("Unknown command: " + val)
                }
        }
        return m, cmd
}

func (m *AppViewModel) View() string {
        baseView := view.Render(m.state, m.animatorClient.Position())

        baseView += "\n\n" + m.textInput.View()

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
