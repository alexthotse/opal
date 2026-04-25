package main

import (
        "fmt"
        "log"
        "os"

        tea "github.com/charmbracelet/bubbletea"
        "github.com/spf13/cobra"

        "github.com/alexthotse/peregrine/adapters"
        "github.com/alexthotse/peregrine/viewmodel"
)

var (
        theme    string
        provider string
)

func main() {
        var rootCmd = &cobra.Command{
                Use:   "peregrine",
                Short: "Peregrine is a terminal-based AI assistant system",
                Run: func(cmd *cobra.Command, args []string) {
                        runApp(theme, provider)
                },
        }

        rootCmd.PersistentFlags().StringVar(&theme, "theme", "pi", "Theme to use (pi, freecode, crush)")
        rootCmd.PersistentFlags().StringVar(&provider, "provider", "anthropic", "LLM provider to use")

        if err := rootCmd.Execute(); err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
}

func runApp(themeStr, providerStr string) {
        f, err := tea.LogToFile("debug.log", "debug")
        if err != nil {
                fmt.Println("fatal:", err)
                os.Exit(1)
        }
        defer f.Close()

        // Hexagonal Architecture Initialization
        backendClient := adapters.NewDefaultBackendClient()

        agentClient, err := adapters.NewADKAgentClient(providerStr)
        if err != nil {
                log.Fatalf("Failed to initialize ADK Agent: %v", err)
        }

        animatorClient := adapters.NewAnimator()

        // MVVM Bootstrap
        appViewModel := viewmodel.NewAppViewModel(backendClient, agentClient, animatorClient, themeStr)

        p := tea.NewProgram(appViewModel)
        if _, err := p.Run(); err != nil {
                log.Fatal(err)
        }
}
