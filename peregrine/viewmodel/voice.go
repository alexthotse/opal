package viewmodel

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

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

type ARecordAudioRecorder struct {
	mutex            sync.Mutex
	cmd              *exec.Cmd
	cancel           context.CancelFunc
	waitDone         chan struct{}
	buffer           []int16
	maxBufferSamples int
	running          bool
	stopRequested    bool
	procErr          error
	closed           bool
}

func NewARecordAudioRecorder() (*ARecordAudioRecorder, error) {
	if _, err := exec.LookPath("arecord"); err != nil {
		return nil, err
	}

	return &ARecordAudioRecorder{
		buffer:           make([]int16, 0, 16000),
		maxBufferSamples: 16000 * 60,
	}, nil
}

func (r *ARecordAudioRecorder) Start() error {
	if r == nil {
		return nil
	}

	r.mutex.Lock()
	if r.closed {
		r.mutex.Unlock()
		return errors.New("audio recorder is closed")
	}
	if r.running {
		r.mutex.Unlock()
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	args := []string{"-q", "-t", "raw", "-f", "S16_LE", "-c", "1", "-r", "16000"}
	if device := os.Getenv("PEREGRINE_ARECORD_DEVICE"); device != "" {
		args = append(args, "-D", device)
	}
	args = append(args, "-")

	cmd := exec.CommandContext(ctx, "arecord", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		r.mutex.Unlock()
		return err
	}
	cmd.Stderr = io.Discard

	if err := cmd.Start(); err != nil {
		_ = stdout.Close()
		cancel()
		r.mutex.Unlock()
		return err
	}

	r.cmd = cmd
	r.cancel = cancel
	r.waitDone = make(chan struct{})
	r.buffer = r.buffer[:0]
	r.procErr = nil
	r.running = true
	r.stopRequested = false
	waitDone := r.waitDone
	r.mutex.Unlock()

	go r.readLoop(cmd, stdout, waitDone)
	return nil
}

func (r *ARecordAudioRecorder) readLoop(cmd *exec.Cmd, stdout io.ReadCloser, waitDone chan struct{}) {
	defer close(waitDone)
	defer func() { _ = stdout.Close() }()

	buf := make([]byte, 4096)
	var pendingByte byte
	var hasPending bool
	var readErr error

	for {
		n, err := stdout.Read(buf)
		if n > 0 {
			data := buf[:n]
			if hasPending {
				if len(data) > 0 {
					sample := int16(uint16(pendingByte) | (uint16(data[0]) << 8))
					r.appendSamples([]int16{sample})
					data = data[1:]
					hasPending = false
				}
			}

			if len(data)%2 == 1 {
				pendingByte = data[len(data)-1]
				hasPending = true
				data = data[:len(data)-1]
			}

			if len(data) > 0 {
				samples := len(data) / 2
				converted := make([]int16, samples)
				for i := 0; i < samples; i++ {
					lo := uint16(data[i*2])
					hi := uint16(data[i*2+1]) << 8
					converted[i] = int16(lo | hi)
				}
				r.appendSamples(converted)
			}
		}

		if err != nil {
			readErr = err
			break
		}
	}

	waitErr := cmd.Wait()

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.running = false
	if r.stopRequested {
		r.procErr = nil
	} else if waitErr != nil {
		r.procErr = waitErr
	} else if readErr != nil && !errors.Is(readErr, io.EOF) {
		r.procErr = readErr
	} else {
		r.procErr = io.EOF
	}

	r.cmd = nil
	if r.cancel != nil {
		r.cancel()
	}
	r.cancel = nil
}

func (r *ARecordAudioRecorder) appendSamples(samples []int16) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !r.running || r.closed {
		return
	}

	r.buffer = append(r.buffer, samples...)
	if r.maxBufferSamples > 0 && len(r.buffer) > r.maxBufferSamples {
		overflow := len(r.buffer) - r.maxBufferSamples
		copy(r.buffer, r.buffer[overflow:])
		r.buffer = r.buffer[:len(r.buffer)-overflow]
	}
}

func (r *ARecordAudioRecorder) Stop() error {
	if r == nil {
		return nil
	}

	r.mutex.Lock()
	if !r.running {
		r.mutex.Unlock()
		return nil
	}
	r.stopRequested = true
	cmd := r.cmd
	cancel := r.cancel
	waitDone := r.waitDone
	r.mutex.Unlock()

	if cmd != nil && cmd.Process != nil {
		_ = cmd.Process.Signal(os.Interrupt)
	}

	if waitDone != nil {
		select {
		case <-waitDone:
			return nil
		case <-time.After(1200 * time.Millisecond):
		}
	}

	if cancel != nil {
		cancel()
	}

	if waitDone != nil {
		select {
		case <-waitDone:
			return nil
		case <-time.After(1200 * time.Millisecond):
		}
	}

	if cmd != nil && cmd.Process != nil {
		_ = cmd.Process.Kill()
	}

	if waitDone != nil {
		select {
		case <-waitDone:
		case <-time.After(1200 * time.Millisecond):
		}
	}

	return nil
}

func (r *ARecordAudioRecorder) Read(p []int16) (int, error) {
	if r == nil {
		return 0, nil
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if len(r.buffer) == 0 {
		if r.closed {
			return 0, io.EOF
		}
		if !r.running && r.procErr != nil {
			return 0, r.procErr
		}
		return 0, nil
	}

	n := copy(p, r.buffer)
	r.buffer = r.buffer[n:]
	return n, nil
}

func (r *ARecordAudioRecorder) Close() error {
	if r == nil {
		return nil
	}

	r.mutex.Lock()
	r.closed = true
	r.mutex.Unlock()
	return r.Stop()
}

type VoiceModel struct {
	active   bool
	recorder AudioRecorder
	tts      adapters.TTSProvider
}

func NewVoiceModel() VoiceModel {
	recorder, err := NewARecordAudioRecorder()
	if err != nil {
		log.Printf("Warning: Failed to initialize arecord audio recorder: %v", err)
	}

	return VoiceModel{
		active:   false,
		recorder: recorder,
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
