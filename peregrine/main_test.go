package main

import (
	"bytes"
	"testing"
	"time"

	"github.com/alexthotse/peregrine/adapters"
	"github.com/alexthotse/peregrine/viewmodel"
	"github.com/charmbracelet/x/exp/teatest"
)

func TestMainApp(t *testing.T) {
	bc := adapters.NewDefaultBackendClient()
	ac, err := adapters.NewADKAgentClient()
	if err != nil {
		t.Fatalf("Failed to initialize ADK Agent: %v", err)
	}
	anim := adapters.NewAnimator()
	app := viewmodel.NewAppViewModel(bc, ac, anim)

	tm := teatest.NewTestModel(t, app, teatest.WithInitialTermSize(80, 24))

	// Ensure the app starts and renders something
	teatest.WaitFor(t, tm.Output(), func(b []byte) bool {
		return bytes.Contains(b, []byte("Peregrine")) || bytes.Contains(b, []byte("Backend")) || bytes.Contains(b, []byte("RPC Error"))
	}, teatest.WithDuration(time.Second*3))

	tm.Quit()
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}
