package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"google.golang.org/adk/agent"
)

type teaModel struct {
	backendCmd *exec.Cmd
	stdin      io.WriteCloser
	stdout     io.ReadCloser
	messages   []string
	requestID  int
	theme      Theme
	animator   *Animator
	adkAgent   agent.Agent
}

type JSONRPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      string `json:"id"`
	Method  string `json:"method"`
}

type JSONRPCResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  string `json:"result"`
	Method  string `json:"method"` // For events
}

type backendMsg struct {
	content string
}

func initialModel() teaModel {
	anim := NewAnimator()
	anim.SetTarget(20.0) // Target width for an animation reveal effect
	
	adkAgent, _ := initADKAgent()

	return teaModel{
		messages: []string{"Starting Peregrine Backend (Gleam)..."},
		theme:    GetTheme(ThemePI),
		animator: anim,
		adkAgent: adkAgent,
	}
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m teaModel) Init() tea.Cmd {
	return tea.Batch(startBackend, tickCmd())
}

func startBackend() tea.Msg {
	cmd := exec.Command("gleam", "run")
	if backendPath := os.Getenv("OPAL_BACKEND_DIR"); backendPath != "" {
		cmd.Dir = backendPath
	} else {
		cmd.Dir = "../peregrine_backend"
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return backendMsg{fmt.Sprintf("Error connecting stdin: %v", err)}
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return backendMsg{fmt.Sprintf("Error connecting stdout: %v", err)}
	}
	if err := cmd.Start(); err != nil {
		return backendMsg{fmt.Sprintf("Error starting gleam: %v", err)}
	}

	return backendStartedMsg{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
	}
}

type backendStartedMsg struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
}

func listenBackend(stdout io.ReadCloser) tea.Cmd {
	return func() tea.Msg {
		reader := bufio.NewReader(stdout)
		line, err := reader.ReadString('\n')
		if err != nil {
			return backendMsg{"Backend disconnected."}
		}
		return backendMsg{line}
	}
}

func sendRequest(stdin io.WriteCloser, reqID int, method string) tea.Cmd {
	return func() tea.Msg {
		req := JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      fmt.Sprintf("%d", reqID),
			Method:  method,
		}
		b, _ := json.Marshal(req)
		fmt.Fprintf(stdin, "%s\n", b)
		return nil
	}
}

func (m teaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.backendCmd != nil {
				m.backendCmd.Process.Kill()
			}
			return m, tea.Quit
		case "p": // Send ping
			if m.stdin != nil {
				m.requestID++
				return m, sendRequest(m.stdin, m.requestID, "ping")
			}
		case "s": // ADK mock search
			if m.stdin != nil {
				m.messages = append(m.messages, runADKAgent(m.adkAgent, "search logs"))
				m.requestID++
				return m, sendRequest(m.stdin, m.requestID, "search.quick")
			}
		case "a": // Start agent
			if m.stdin != nil {
				// 1:1 Matching Workflow:
				// Frontend ADK Agent delegates to Backend Jido Agent
				adkOutput := runADKAgent(m.adkAgent, "trigger backend action")
				m.messages = append(m.messages, adkOutput)

				m.requestID++
				return m, sendRequest(m.stdin, m.requestID, "jido.action")
			}
		case "u": // Start ultrathink
			if m.stdin != nil {
				m.requestID++
				return m, sendRequest(m.stdin, m.requestID, "ultrathink.start")
			}
		case "r": // Start ultraplan
			if m.stdin != nil {
				m.requestID++
				return m, sendRequest(m.stdin, m.requestID, "ultraplan.start")
			}
		case "t": // Toggle theme
			if m.theme.Primary.GetForeground() == lipgloss.Color("39") {
				m.theme = GetTheme(ThemeFreeCode)
			} else if m.theme.Primary.GetForeground() == lipgloss.Color("0") {
				m.theme = GetTheme(ThemeCrush)
			} else {
				m.theme = GetTheme(ThemePI)
			}
		case "v": // Start voice mode
			m.messages = append(m.messages, StartVoiceMode())
		case "h": // Show history
			m.messages = append(m.messages, ShowHistoryPicker())
		case "m": // Show message actions
			m.messages = append(m.messages, ShowMessageActions())
		}

	case tickMsg:
		m.animator.Update()
		return m, tickCmd()

	case backendStartedMsg:
		m.backendCmd = msg.cmd
		m.stdin = msg.stdin
		m.stdout = msg.stdout
		m.messages = append(m.messages, "Backend connected! Keys: 'p' (ping), 'a' (agent), 'u' (think), 'r' (plan), 'v' (voice), 'h' (history), 'm' (actions), 's' (search), 't' (theme), 'q' (quit).")
		return m, listenBackend(m.stdout)

	case backendMsg:
		parsed := strings.TrimSpace(msg.content)
		if parsed != "" {
			m.messages = append(m.messages, "Gleam: "+parsed)
		}
		if m.stdout != nil && !strings.Contains(parsed, "disconnected") {
			return m, listenBackend(m.stdout)
		}
	}

	return m, nil
}

func (m teaModel) View() string {
	pad := int(m.animator.position)
	if pad < 0 {
		pad = 0
	}
	paddingStyle := lipgloss.NewStyle().PaddingLeft(pad)

	s := m.theme.Primary.Render("✦ Peregrine Agent") + "\n\n"
	for _, msg := range m.messages {
		s += m.theme.Text.Render(msg) + "\n"
	}
	return paddingStyle.Render(s)
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
