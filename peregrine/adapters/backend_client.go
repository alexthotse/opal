package adapters

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type BackendStartedMsg struct {
	Cmd    *exec.Cmd
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
}

type BackendMsg struct {
	Content string
}

type JSONRPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      string `json:"id"`
	Method  string `json:"method"`
}

type DefaultBackendClient struct{}

func NewDefaultBackendClient() *DefaultBackendClient {
	return &DefaultBackendClient{}
}

func (c *DefaultBackendClient) Start() tea.Msg {
	cmd := exec.Command("gleam", "run")
	if backendPath := os.Getenv("FALCON_DIR"); backendPath != "" {
		cmd.Dir = backendPath
	} else {
		cmd.Dir = "../falcon"
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return BackendMsg{fmt.Sprintf("Error connecting stdin: %v", err)}
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return BackendMsg{fmt.Sprintf("Error connecting stdout: %v", err)}
	}
	if err := cmd.Start(); err != nil {
		return BackendMsg{fmt.Sprintf("Error starting gleam: %v", err)}
	}

	return BackendStartedMsg{
		Cmd:    cmd,
		Stdin:  stdin,
		Stdout: stdout,
	}
}

func (c *DefaultBackendClient) Listen(stdout io.ReadCloser) tea.Cmd {
	return func() tea.Msg {
		reader := bufio.NewReader(stdout)
		line, err := reader.ReadString('\n')
		if err != nil {
			return BackendMsg{"Backend disconnected."}
		}
		return BackendMsg{line}
	}
}

func (c *DefaultBackendClient) SendRequest(stdin io.WriteCloser, reqID int, method string) tea.Cmd {
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
