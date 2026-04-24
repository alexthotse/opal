package main

import (
	"context"
	"fmt"
	"iter"

	"google.golang.org/genai"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
)

type mockModel struct{}

func (m mockModel) Name() string { return "mock-model" }
func (m mockModel) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		text := "This is a mock ADK response."
		content := &genai.Content{
			Parts: []*genai.Part{{Text: text}},
		}
		yield(&model.LLMResponse{Content: content}, nil)
	}
}

func initADKAgent() (agent.Agent, error) {
	return llmagent.New(llmagent.Config{
		Name:        "peregrine_frontend_agent",
		Model:       mockModel{},
		Description: "Frontend agent",
		Instruction: "You process user commands.",
	})
}

func runADKAgent(a agent.Agent, prompt string) string {
	return fmt.Sprintf("ADK Agent Processing Prompt: '%s'", prompt)
}
