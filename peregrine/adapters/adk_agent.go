package adapters

import (
	"context"
	"errors"
	"fmt"
	"io"
	"iter"
	"os"

	"github.com/liushuangls/go-anthropic/v2"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

type geminiModel struct {
	client *genai.Client
}

func (m *geminiModel) Name() string { return "gemini" }

func (m *geminiModel) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		modelName := req.Model
		if modelName == "" {
			modelName = "gemini-2.5-flash"
		}

		if stream {
			resStream := m.client.Models.GenerateContentStream(ctx, modelName, req.Contents, req.Config)
			for res, err := range resStream {
				if err != nil {
					yield(nil, err)
					return
				}
				if len(res.Candidates) > 0 {
					candidate := res.Candidates[0]
					resp := &model.LLMResponse{
						Content:           candidate.Content,
						CitationMetadata:  candidate.CitationMetadata,
						GroundingMetadata: candidate.GroundingMetadata,
						UsageMetadata:     res.UsageMetadata,
						LogprobsResult:    candidate.LogprobsResult,
						ModelVersion:      res.ModelVersion,
						FinishReason:      candidate.FinishReason,
					}
					if !yield(resp, nil) {
						return
					}
				}
			}
		} else {
			res, err := m.client.Models.GenerateContent(ctx, modelName, req.Contents, req.Config)
			if err != nil {
				yield(nil, err)
				return
			}
			if len(res.Candidates) > 0 {
				candidate := res.Candidates[0]
				resp := &model.LLMResponse{
					Content:           candidate.Content,
					CitationMetadata:  candidate.CitationMetadata,
					GroundingMetadata: candidate.GroundingMetadata,
					UsageMetadata:     res.UsageMetadata,
					LogprobsResult:    candidate.LogprobsResult,
					ModelVersion:      res.ModelVersion,
					FinishReason:      candidate.FinishReason,
				}
				yield(resp, nil)
			}
		}
	}
}

type openAIModel struct {
	client *openai.Client
}

func (m *openAIModel) Name() string { return "openai" }

func convertToOpenAIMessages(contents []*genai.Content) []openai.ChatCompletionMessage {
	var msgs []openai.ChatCompletionMessage
	for _, c := range contents {
		role := openai.ChatMessageRoleUser
		if c.Role == "model" {
			role = openai.ChatMessageRoleAssistant
		} else if c.Role == "system" {
			role = openai.ChatMessageRoleSystem
		}

		var text string
		for _, p := range c.Parts {
			text += p.Text
		}

		msgs = append(msgs, openai.ChatCompletionMessage{
			Role:    role,
			Content: text,
		})
	}
	return msgs
}

func (m *openAIModel) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		modelName := req.Model
		if modelName == "" {
			modelName = openai.GPT4o
		}

		messages := convertToOpenAIMessages(req.Contents)

		if stream {
			streamRes, err := m.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
				Model:    modelName,
				Messages: messages,
				Stream:   true,
			})
			if err != nil {
				yield(nil, err)
				return
			}
			defer streamRes.Close()

			for {
				resp, err := streamRes.Recv()
				if errors.Is(err, io.EOF) {
					break
				}
				if err != nil {
					yield(nil, err)
					break
				}
				if len(resp.Choices) > 0 {
					text := resp.Choices[0].Delta.Content
					if text != "" {
						content := &genai.Content{
							Parts: []*genai.Part{{Text: text}},
							Role:  "model",
						}
						llmResp := &model.LLMResponse{
							Content: content,
						}
						if !yield(llmResp, nil) {
							return
						}
					}
				}
			}
		} else {
			resp, err := m.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
				Model:    modelName,
				Messages: messages,
			})
			if err != nil {
				yield(nil, err)
				return
			}
			if len(resp.Choices) > 0 {
				text := resp.Choices[0].Message.Content
				content := &genai.Content{
					Parts: []*genai.Part{{Text: text}},
					Role:  "model",
				}
				llmResp := &model.LLMResponse{
					Content: content,
				}
				yield(llmResp, nil)
			}
		}
	}
}

type anthropicModel struct {
	client *anthropic.Client
}

func (m *anthropicModel) Name() string { return "anthropic" }

func convertToAnthropicMessages(contents []*genai.Content) []anthropic.Message {
	var msgs []anthropic.Message
	for _, c := range contents {
		role := anthropic.RoleUser
		if c.Role == "model" {
			role = anthropic.RoleAssistant
		}

		var text string
		for _, p := range c.Parts {
			text += p.Text
		}

		if c.Role != "system" {
			msgs = append(msgs, anthropic.Message{
				Role:    role,
				Content: []anthropic.MessageContent{anthropic.NewTextMessageContent(text)},
			})
		}
	}
	return msgs
}

func extractAnthropicSystem(contents []*genai.Content) string {
	var sys string
	for _, c := range contents {
		if c.Role == "system" {
			for _, p := range c.Parts {
				sys += p.Text
			}
		}
	}
	return sys
}

func (m *anthropicModel) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		modelName := req.Model
		if modelName == "" {
			modelName = string(anthropic.ModelClaude3Dot5SonnetLatest)
		}

		messages := convertToAnthropicMessages(req.Contents)
		system := extractAnthropicSystem(req.Contents)

		if stream {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			messagesReq := anthropic.MessagesStreamRequest{
				MessagesRequest: anthropic.MessagesRequest{
					Model:     anthropic.Model(modelName),
					Messages:  messages,
					System:    system,
					MaxTokens: 4096,
				},
				OnContentBlockDelta: func(data anthropic.MessagesEventContentBlockDeltaData) {
					if data.Delta.Text != nil && *data.Delta.Text != "" {
						content := &genai.Content{
							Parts: []*genai.Part{{Text: *data.Delta.Text}},
							Role:  "model",
						}
						llmResp := &model.LLMResponse{Content: content}
						if !yield(llmResp, nil) {
							cancel()
						}
					}
				},
			}

			_, err := m.client.CreateMessagesStream(ctx, messagesReq)
			if err != nil && !errors.Is(err, context.Canceled) {
				yield(nil, err)
			}
		} else {
			resp, err := m.client.CreateMessages(ctx, anthropic.MessagesRequest{
				Model:     anthropic.Model(modelName),
				Messages:  messages,
				System:    system,
				MaxTokens: 4096,
			})
			if err != nil {
				yield(nil, err)
				return
			}
			if len(resp.Content) > 0 {
				text := resp.Content[0].Text
				if text != nil {
					content := &genai.Content{
						Parts: []*genai.Part{{Text: *text}},
						Role:  "model",
					}
					llmResp := &model.LLMResponse{Content: content}
					yield(llmResp, nil)
				}
			}
		}
	}
}

type ADKAgentClient struct {
	agent agent.Agent
}

func NewADKAgentClient(provider string) (*ADKAgentClient, error) {
	var m model.LLM
	switch provider {
	case "gemini":
		key := os.Getenv("GEMINI_API_KEY")
		if key == "" {
			return nil, errors.New("GEMINI_API_KEY is not set")
		}
		client, err := genai.NewClient(context.Background(), &genai.ClientConfig{APIKey: key})
		if err != nil {
			return nil, err
		}
		m = &geminiModel{client: client}
	case "anthropic":
		key := os.Getenv("ANTHROPIC_API_KEY")
		if key == "" {
			return nil, errors.New("ANTHROPIC_API_KEY is not set")
		}
		client := anthropic.NewClient(key)
		m = &anthropicModel{client: client}
	case "openai":
		key := os.Getenv("OPENAI_API_KEY")
		if key == "" {
			return nil, errors.New("OPENAI_API_KEY is not set")
		}
		client := openai.NewClient(key)
		m = &openAIModel{client: client}
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        "peregrine_agent_" + provider,
		Model:       m,
		Description: "Frontend agent using " + provider,
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
