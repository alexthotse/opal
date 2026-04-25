package llm

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

type Client struct {
	client *genai.Client
}

func NewClient() (*Client, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}
	return &Client{
		client: client,
	}, nil
}

func (c *Client) GenerateReasoning(prompt string) (string, error) {
	ctx := context.Background()
	// Usually gemini-2.5-flash or gemini-2.0-flash-thinking-exp for reasoning
	resp, err := c.client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate reasoning: %w", err)
	}
	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response generated")
	}
	
	// Assuming it's a Text part
	part := resp.Candidates[0].Content.Parts[0]
	if part.Text != "" {
		return part.Text, nil
	}
	
	return "", fmt.Errorf("unexpected response format")
}
