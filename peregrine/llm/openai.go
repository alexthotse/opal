package llm

import (
	"context"
	"fmt"
	"iter"
	"os"

	"github.com/sashabaranov/go-openai"
	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

type OpenAIWrapper struct {
	client    *openai.Client
	modelName string
}

func NewOpenAIWrapper(modelName string) (*OpenAIWrapper, error) {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is not set")
	}
	return &OpenAIWrapper{
		client:    openai.NewClient(key),
		modelName: modelName,
	}, nil
}

func (m *OpenAIWrapper) Name() string {
	return m.modelName
}

func (m *OpenAIWrapper) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		var messages []openai.ChatCompletionMessage
		for _, content := range req.Contents {
			role := openai.ChatMessageRoleUser
			if content.Role == "model" {
				role = openai.ChatMessageRoleAssistant
			} else if content.Role == "system" {
				role = openai.ChatMessageRoleSystem
			}

			var text string
			for _, part := range content.Parts {
				if part.Text != "" {
					text += part.Text
				}
			}
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    role,
				Content: text,
			})
		}

		resp, err := m.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model:    m.modelName,
				Messages: messages,
			},
		)

		if err != nil {
			yield(nil, err)
			return
		}

		if len(resp.Choices) > 0 {
			yield(&model.LLMResponse{
				Content: &genai.Content{
					Role: "model",
					Parts: []*genai.Part{
						genai.NewPartFromText(resp.Choices[0].Message.Content),
					},
				},
			}, nil)
		}
	}
}
