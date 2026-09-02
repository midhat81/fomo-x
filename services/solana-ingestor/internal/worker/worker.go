package worker

import (
	"context"
	"log"

	"github.com/midhat81/fomo-x/services/solana-ingestor/internal/checkpoint"
	"github.com/midhat81/fomo-x/services/solana-ingestor/internal/kafka"
	"github.com/midhat81/fomo-x/services/solana-ingestor/internal/solana"
)

// DefaultProgramID is the Solana program we watch for trade activity.
// This is the Raydium AMM v4 program — one of the most active swap
// programs on Solana, good for generating real Day 1 test traffic.
const DefaultProgramID = "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"

// Worker ties together the websocket subscription, parsing, decoding,
// checkpointing, and Kafka publishing into a single pipeline.
type Worker struct {
	ws         *solana.WSClient
	producer   *kafka.Producer
	checkpoint *checkpoint.Store
	programID  string
}

// New creates a new ingestion worker.
func New(ws *solana.WSClient, producer *kafka.Producer, cp *checkpoint.Store) *Worker {
	return &Worker{
		ws:         ws,
		producer:   producer,
		checkpoint: cp,
		programID:  DefaultProgramID,
	}
}

// Run starts the ingestion pipeline: subscribe to logs, parse, decode,
// publish to Kafka, checkpoint. Blocks until ctx is cancelled.
func (w *Worker) Run(ctx context.Context) error {
	logs := make(chan solana.LogEvent, 100)

	if err := w.ws.SubscribeProgramLogs(ctx, w.programID, logs); err != nil {
		return err
	}

	log.Printf("Worker subscribed to program logs: %s\n", w.programID)

	for {
		select {
		case <-ctx.Done():
			return nil

		case evt := <-logs:
			tx := solana.ParseLogEvent(evt)

			trade, ok := solana.DecodeTrade(tx)
			if !ok {
				continue
			}

			if err := w.producer.PublishTrade(ctx, trade); err != nil {
				log.Printf("failed to publish trade event: %v\n", err)
				continue
			}

			if err := w.checkpoint.Save(tx.Signature); err != nil {
				log.Printf("failed to save checkpoint: %v\n", err)
			}

			log.Printf("Published trade event: signature=%s event_id=%s\n", tx.Signature, trade.EventID)
		}
	}
}