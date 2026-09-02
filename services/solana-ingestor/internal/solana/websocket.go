package solana

import (
	"context"
	"fmt"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

// WSClient wraps a Solana WebSocket connection for streaming live logs.
type WSClient struct {
	client *ws.Client
	wsURL  string
}

// NewWSClient connects to the given Solana WebSocket URL.
func NewWSClient(ctx context.Context, wsURL string) (*WSClient, error) {
	client, err := ws.Connect(ctx, wsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to solana websocket: %w", err)
	}
	return &WSClient{client: client, wsURL: wsURL}, nil
}

// LogEvent is a simplified representation of a log notification we care about.
type LogEvent struct {
	Signature string
	Err       interface{}
	Logs      []string
}

// SubscribeProgramLogs subscribes to all transaction logs mentioning the given
// program ID (e.g. a DEX program) and streams them to the returned channel.
// The caller is responsible for reading from the channel and calling Close
// when done.
func (w *WSClient) SubscribeProgramLogs(ctx context.Context, programID string, out chan<- LogEvent) error {
	pubkey, err := solanago.PublicKeyFromBase58(programID)
	if err != nil {
		return fmt.Errorf("invalid program id: %w", err)
	}

	sub, err := w.client.LogsSubscribeMentions(pubkey, rpc.CommitmentFinalized)
	if err != nil {
		return fmt.Errorf("failed to subscribe to logs: %w", err)
	}

	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			got, err := sub.Recv(ctx)
			if err != nil {
				return
			}
			if got == nil {
				continue
			}

			out <- LogEvent{
				Signature: got.Value.Signature.String(),
				Err:       got.Value.Err,
				Logs:      got.Value.Logs,
			}
		}
	}()

	return nil
}

// Close closes the underlying websocket connection.
func (w *WSClient) Close() {
	w.client.Close()
}
