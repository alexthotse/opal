package adapters

import (
	"context"
	"fmt"
	"iter"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/genai"
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

type ADKAgentClient struct {
	agent agent.Agent
}

func NewADKAgentClient() (*ADKAgentClient, error) {
	a, err := llmagent.New(llmagent.Config{
		Name:        "peregrine_agent",
		Model:       mockModel{},
		Description: "Frontend agent",
		Instruction: "You process user commands.",
	})
	if err != nil {
		return nil, err
	}
	return &ADKAgentClient{agent: a}, nil
}

func (c *ADKAgentClient) ProcessPrompt(prompt string) string {
	return fmt.Sprintf("ADK Agent Processing Prompt: '%s'", prompt)
}
