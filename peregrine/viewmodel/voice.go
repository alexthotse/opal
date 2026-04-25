package viewmodel

import (
	"log"
	"sync"
	"unsafe"

	"github.com/alexthotse/peregrine/adapters"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gen2brain/malgo"
)

type AudioRecorder interface {
	Start() error
	Stop() error
	Read(p []int16) (int, error)
	Close() error
}

type MalgoAudioRecorder struct {
	ctx         *malgo.AllocatedContext
	device      *malgo.Device
	buffer      []int16
	mutex       sync.Mutex
	isRecording bool
}

func NewMalgoAudioRecorder() (*MalgoAudioRecorder, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		// Log suppressed to avoid noise
	})
	if err != nil {
		return nil, err
	}

	recorder := &MalgoAudioRecorder{
		ctx: ctx,
	}

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.SampleRate = 16000
	deviceConfig.Alsa.NoMMap = 1

	callbacks := malgo.DeviceCallbacks{
		Data: func(pOutputSample, pInputSamples []byte, framecount uint32) {
			recorder.mutex.Lock()
			defer recorder.mutex.Unlock()

			if !recorder.isRecording || len(pInputSamples) == 0 {
				return
			}

			samples := len(pInputSamples) / 2
			pInt16 := unsafe.Slice((*int16)(unsafe.Pointer(&pInputSamples[0])), samples)

			recorder.buffer = append(recorder.buffer, pInt16...)
		},
	}

	device, err := malgo.InitDevice(ctx.Context, deviceConfig, callbacks)
	if err != nil {
		ctx.Free()
		return nil, err
	}
	recorder.device = device

	return recorder, nil
}

func (m *MalgoAudioRecorder) Start() error {
	if m == nil || m.device == nil {
		return nil
	}
	m.mutex.Lock()
	m.buffer = m.buffer[:0]
	m.isRecording = true
	m.mutex.Unlock()

	return m.device.Start()
}

func (m *MalgoAudioRecorder) Stop() error {
	if m == nil || m.device == nil {
		return nil
	}
	m.mutex.Lock()
	m.isRecording = false
	m.mutex.Unlock()

	return m.device.Stop()
}

func (m *MalgoAudioRecorder) Read(p []int16) (int, error) {
	if m == nil {
		return 0, nil
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.buffer) == 0 {
		return 0, nil
	}

	n := copy(p, m.buffer)
	m.buffer = m.buffer[n:]
	return n, nil
}

func (m *MalgoAudioRecorder) Close() error {
	if m == nil {
		return nil
	}
	if m.device != nil {
		m.device.Uninit()
	}
	if m.ctx != nil {
		m.ctx.Free()
	}
	return nil
}

type VoiceModel struct {
	active   bool
	recorder AudioRecorder
	tts      adapters.TTSProvider
}

func NewVoiceModel() VoiceModel {
	recorder, err := NewMalgoAudioRecorder()
	if err != nil {
		log.Printf("Warning: Failed to initialize malgo audio recorder: %v", err)
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
