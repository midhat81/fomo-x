package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go/rpc"
)

// Client wraps the Solana RPC connection used by the ingestor.
type Client struct {
	rpcClient *rpc.Client
	rpcURL    string
}

// NewClient creates a new Solana RPC client from the given RPC URL.
func NewClient(rpcURL string) *Client {
	return &Client{
		rpcClient: rpc.New(rpcURL),
		rpcURL:    rpcURL,
	}
}

// HealthCheck verifies the RPC endpoint is reachable and responding.
func (c *Client) HealthCheck(ctx context.Context) error {
	out, err := c.rpcClient.GetHealth(ctx)
	if err != nil {
		return fmt.Errorf("solana rpc health check failed: %w", err)
	}
	if out != rpc.HealthOk {
		return fmt.Errorf("solana rpc unhealthy: %s", out)
	}
	return nil
}

// GetLatestBlockhash fetches the most recent blockhash, useful as a basic
// connectivity/sanity check beyond HealthCheck.
func (c *Client) GetLatestBlockhash(ctx context.Context) (string, error) {
	out, err := c.rpcClient.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return "", fmt.Errorf("failed to get latest blockhash: %w", err)
	}
	return out.Value.Blockhash.String(), nil
}