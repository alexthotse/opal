package llm

import (
	"context"
	"fmt"
	"iter"
	"os"

	"github.com/liushuangls/go-anthropic/v2"
	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

type AnthropicWrapper struct {
	client    *anthropic.Client
	modelName string
}

func NewAnthropicWrapper(modelName string) (*AnthropicWrapper, error) {
	key := os.Getenv("ANTHROPIC_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY is not set")
	}
	return &AnthropicWrapper{
		client:    anthropic.NewClient(key),
		modelName: modelName,
	}, nil
}

func (m *AnthropicWrapper) Name() string {
	return m.modelName
}

func (m *AnthropicWrapper) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		var messages []anthropic.Message
		var systemPrompt string

		for _, content := range req.Contents {
			var text string
			for _, part := range content.Parts {
				if part.Text != "" {
					text += part.Text
				}
			}

			if content.Role == "system" {
				systemPrompt += text
				continue
			}

			role := anthropic.RoleUser
			if content.Role == "model" {
				role = anthropic.RoleAssistant
			}

			messages = append(messages, anthropic.Message{
				Role: role,
				Content: []anthropic.MessageContent{
					anthropic.NewTextMessageContent(text),
				},
			})
		}

		anthropicReq := anthropic.MessagesRequest{
			Model:     anthropic.Model(m.modelName),
			Messages:  messages,
			MaxTokens: 4096,
		}
		if systemPrompt != "" {
			anthropicReq.System = systemPrompt
		}

		resp, err := m.client.CreateMessages(ctx, anthropicReq)
		if err != nil {
			yield(nil, err)
			return
		}

		if len(resp.Content) > 0 {
			var respText string
			for _, c := range resp.Content {
				if c.Type == anthropic.MessagesContentTypeText {
					respText += *c.Text
				}
			}
			yield(&model.LLMResponse{
				Content: &genai.Content{
					Role: "model",
					Parts: []*genai.Part{
						genai.NewPartFromText(respText),
					},
				},
			}, nil)
		}
	}
}
