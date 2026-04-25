package adapters

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	falconv1 "github.com/alexthotse/peregrine/gen/falcon/v1"
	"github.com/alexthotse/peregrine/gen/falcon/v1/falconv1connect"
)

type BackendClient struct {
	client falconv1connect.FalconServiceClient
}

func NewDefaultBackendClient() *BackendClient {
	// Initialize the ConnectRPC client using the custom MessagePack Codec over HTTP
	client := falconv1connect.NewFalconServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
		connect.WithCodec(NewMsgPackCodec()),
	)
	return &BackendClient{
		client: client,
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
func (b *BackendClient) StartUltrathink(ctx context.Context, id string) (string, error) {
	req := connect.NewRequest(&falconv1.UltrathinkRequest{Id: id})
	res, err := b.client.StartUltrathink(ctx, req)
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
