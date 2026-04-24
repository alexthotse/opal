package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"peregrine/adapters"
	"peregrine/viewmodel"
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	// Hexagonal Architecture Initialization
	backendClient := adapters.NewDefaultBackendClient()

	agentClient, err := adapters.NewADKAgentClient()
	if err != nil {
		log.Fatalf("Failed to initialize ADK Agent: %v", err)
	}

	animatorClient := adapters.NewAnimator()

	// MVVM Bootstrap
	appViewModel := viewmodel.NewAppViewModel(backendClient, agentClient, animatorClient)

	p := tea.NewProgram(appViewModel)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
