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

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	backendCmd *exec.Cmd
	stdin      io.WriteCloser
	stdout     io.ReadCloser
	messages   []string
	requestID  int
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

func initialModel() model {
	return model{
		messages: []string{"Starting Opal Backend (Gleam)..."},
	}
}

func (m model) Init() tea.Cmd {
	return startBackend
}

func startBackend() tea.Msg {
	cmd := exec.Command("gleam", "run")
	cmd.Dir = "../opal_backend"
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "a": // Start agent
			if m.stdin != nil {
				m.requestID++
				return m, sendRequest(m.stdin, m.requestID, "agent.start")
			}
		}

	case backendStartedMsg:
		m.backendCmd = msg.cmd
		m.stdin = msg.stdin
		m.stdout = msg.stdout
		m.messages = append(m.messages, "Backend connected! Press 'p' to ping, 'a' to start agent, 'q' to quit.")
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

func (m model) View() string {
	s := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render("✦ Opal Agent (Go + Gleam)") + "\n\n"
	for _, msg := range m.messages {
		s += msg + "\n"
	}
	return s
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
