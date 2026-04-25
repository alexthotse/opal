package llm

import (
	"context"
	"fmt"
	"iter"
	"os"

	"google.golang.org/genai"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
)

// GenericModelWrapper implements adk's model.LLM interface, making Peregrine vendor-agnostic.
type GenericModelWrapper struct {
	client    *genai.Client
	modelName string
}

func (m *GenericModelWrapper) Name() string {
	return m.modelName
}

func (m *GenericModelWrapper) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		resp, err := m.client.Models.GenerateContent(ctx, m.modelName, req.Contents, req.Config)
		if err != nil {
			yield(nil, err)
			return
		}
		if len(resp.Candidates) > 0 {
			yield(&model.LLMResponse{
				Content: resp.Candidates[0].Content,
			}, nil)
		}
	}
}

// ADKAgentClient wraps the Google Agent Development Kit (ADK).
type ADKAgentClient struct {
	agent agent.Agent
}

func NewADKAgentClient(ctx context.Context) (*ADKAgentClient, error) {
	// Initialize a generic GenAI client that ADK will orchestrate.
	// You can swap this wrapper with ANY provider (Anthropic, Ollama, OpenAI)
	// because ADK uses the model.LLM interface.
	if os.Getenv("GEMINI_API_KEY") == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize GenAI client: %w", err)
	}

	modelWrapper := &GenericModelWrapper{
		client:    client,
		modelName: "gemini-2.5-flash",
	}

	cfg := llmagent.Config{
		Name:        "reasoning_agent",
		Description: "An agent that performs deep reasoning using ADK.",
		Model:       modelWrapper,
	}

	ag, err := llmagent.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create ADK agent: %w", err)
	}

	return &ADKAgentClient{
		agent: ag,
	}, nil
}

func (c *ADKAgentClient) GenerateReasoning(ctx context.Context, prompt string) (string, error) {
	// Trigger the ADK Agent to reason
	result, err := c.agent.Run(ctx, &agent.RunRequest{
		Input: []*genai.Content{
			{
				Role: "user",
				Parts: []any{
					genai.Text(prompt),
				},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("ADK Agent failed: %w", err)
	}

	if len(result.Messages) > 0 {
		for _, msg := range result.Messages {
			if msg.Role == "model" && len(msg.Parts) > 0 {
				if text, ok := msg.Parts[0].(genai.Text); ok {
					return string(text), nil
				}
			}
		}
	}

	return "No reasoning returned from ADK Agent.", nil
}
