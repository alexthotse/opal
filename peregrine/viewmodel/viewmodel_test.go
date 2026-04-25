package viewmodel_test

import (
        "bytes"
        "testing"
        "time"

        "github.com/alexthotse/peregrine/adapters"
        "github.com/alexthotse/peregrine/viewmodel"
        tea "github.com/charmbracelet/bubbletea"
        "github.com/charmbracelet/x/exp/teatest"
)

func TestAppViewModel(t *testing.T) {
        bc := adapters.NewDefaultBackendClient()
        ac, err := adapters.NewADKAgentClient()
        if err != nil {
                t.Fatalf("Failed to initialize ADK Agent: %v", err)
        }
        anim := adapters.NewAnimator()
        app := viewmodel.NewAppViewModel(bc, ac, anim)

        tm := teatest.NewTestModel(t, app, teatest.WithInitialTermSize(80, 24))

        teatest.WaitFor(t, tm.Output(), func(b []byte) bool {
                return bytes.Contains(b, []byte("Peregrine")) || bytes.Contains(b, []byte("Backend")) || bytes.Contains(b, []byte("RPC Error"))
        }, teatest.WithDuration(time.Second*2))

        tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
        tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

        teatest.WaitFor(t, tm.Output(), func(b []byte) bool {
                return bytes.Contains(b, []byte("Voice Mode"))
        }, teatest.WithDuration(time.Second))

        tm.Send(tea.KeyMsg{Type: tea.KeyEsc})
        tm.Send(tea.KeyMsg{Type: tea.KeyEsc})
        
        //tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}
