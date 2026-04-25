package adapters

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	falconv1 "github.com/alexthotse/peregrine/gen/falcon/v1"
	"github.com/alexthotse/peregrine/gen/falcon/v1/falconv1connect"
	"github.com/alexthotse/peregrine/llm"
)

type BackendClient struct {
	client    falconv1connect.FalconServiceClient
	llmClient *llm.ADKAgentClient
}

func NewDefaultBackendClient() *BackendClient {
	// Initialize the ConnectRPC client using the custom MessagePack Codec over HTTP
	client := falconv1connect.NewFalconServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
		connect.WithCodec(NewMsgPackCodec()),
	)

	llmClient, _ := llm.NewADKAgentClient(context.Background())

	return &BackendClient{
		client:    client,
		llmClient: llmClient,
	}
}

// Ping checks liveness of the Falcon backend
func (b *BackendClient) Ping(ctx context.Context, id string) (string, error) {
	req := connect.NewRequest(&falconv1.PingRequest{Id: id})
	res, err := b.client.Ping(ctx, req)
	if err != nil {
		return "", err
	}
	return res.Msg.Message, nil
}

// StartUltrathink invokes the deep reasoning agent
func (b *BackendClient) StartUltrathink(ctx context.Context, id string, prompt string) (string, error) {
	req := connect.NewRequest(&falconv1.UltrathinkRequest{Id: id, Prompt: prompt})
	res, err := b.client.StartUltrathink(ctx, req)
	if err != nil {
		return "", err
	}
	return res.Msg.Result, nil
}

// StartUltraplan invokes the architecture planner
func (b *BackendClient) StartUltraplan(ctx context.Context, id string, goal string) (string, error) {
	req := connect.NewRequest(&falconv1.UltraplanRequest{Id: id, Goal: goal})
	res, err := b.client.StartUltraplan(ctx, req)
	if err != nil {
		return "", err
	}
	return res.Msg.Result, nil
}

// ExtractMemories retrieves memories
func (b *BackendClient) ExtractMemories(ctx context.Context, id string, sessionID string) (string, error) {
	req := connect.NewRequest(&falconv1.ExtractMemoriesRequest{Id: id, SessionId: sessionID})
	res, err := b.client.ExtractMemories(ctx, req)
	if err != nil {
		return "", err
	}
	return res.Msg.Result, nil
}

// DispatchAction dispatches a Jido action
func (b *BackendClient) DispatchAction(ctx context.Context, id, action string) (string, error) {
	req := connect.NewRequest(&falconv1.ActionRequest{Id: id, Action: action})
	res, err := b.client.DispatchAction(ctx, req)
	if err != nil {
		return "", err
	}
	return res.Msg.Result, nil
}

// GenerateReasoning generates reasoning natively in Go via ADK Agent
func (b *BackendClient) GenerateReasoning(ctx context.Context, prompt string) (string, error) {
	if b.llmClient == nil {
		return "Mock Reasoning (ADK Agent failed to init or API Key missing): " + prompt, nil
	}
	return b.llmClient.GenerateReasoning(ctx, prompt)
}
