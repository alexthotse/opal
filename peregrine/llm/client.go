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
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
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
	runner *runner.Runner
}

func NewADKAgentClient(ctx context.Context) (*ADKAgentClient, error) {
	provider := os.Getenv("LLM_PROVIDER")
	if provider == "" {
		provider = "gemini"
	}

	var modelWrapper model.LLM
	var err error

	switch provider {
	case "anthropic":
		modelWrapper, err = NewAnthropicWrapper("claude-3-7-sonnet-20250219")
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Anthropic wrapper: %w", err)
		}
	case "openai":
		modelWrapper, err = NewOpenAIWrapper("gpt-4o")
		if err != nil {
			return nil, fmt.Errorf("failed to initialize OpenAI wrapper: %w", err)
		}
	case "gemini":
		if os.Getenv("GEMINI_API_KEY") == "" {
			return nil, fmt.Errorf("GEMINI_API_KEY is not set")
		}
		client, err := genai.NewClient(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize GenAI client: %w", err)
		}
		modelWrapper = &GenericModelWrapper{
			client:    client,
			modelName: "gemini-2.5-flash",
		}
	default:
		return nil, fmt.Errorf("unsupported LLM_PROVIDER: %s", provider)
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

	sessionSvc := session.InMemoryService()
	r, err := runner.New(runner.Config{
		AppName:        "peregrine",
		Agent:          ag,
		SessionService: sessionSvc,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create runner: %w", err)
	}

	return &ADKAgentClient{
		runner: r,
	}, nil
}

func (c *ADKAgentClient) GenerateReasoning(ctx context.Context, prompt string) (string, error) {
	// Trigger the ADK Agent to reason
	content := genai.NewContentFromParts([]*genai.Part{
		{Text: prompt},
	}, genai.RoleUser)

	events := c.runner.Run(ctx, "user-1", "session-1", content, agent.RunConfig{})

	var reasoning string
	for event, err := range events {
		if err != nil {
			return "", fmt.Errorf("ADK Agent failed: %w", err)
		}
		if event.Content != nil {
			for _, part := range event.Content.Parts {
				if part.Text != "" {
					reasoning += part.Text
				}
			}
		}
	}

	if reasoning == "" {
		return "No reasoning returned from ADK Agent.", nil
	}

	return reasoning, nil
}
