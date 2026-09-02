package worker

import (
	"context"
	"errors"
	"log"

	"github.com/midhat81/fomo-x/services/portfolio-engine/internal/cache"
	"github.com/midhat81/fomo-x/services/portfolio-engine/internal/domain"
	"github.com/midhat81/fomo-x/services/portfolio-engine/internal/kafka"
	"github.com/midhat81/fomo-x/services/portfolio-engine/internal/repository"
)

// Processor consumes trade events and updates portfolio state.
type Processor struct {
	consumer      *kafka.Consumer
	portfolioRepo *repository.PortfolioRepository
	positionRepo  *repository.PositionRepository
	cache         *cache.Cache
}

// New creates a new portfolio processor.
func New(consumer *kafka.Consumer, portfolioRepo *repository.PortfolioRepository, positionRepo *repository.PositionRepository, c *cache.Cache) *Processor {
	return &Processor{
		consumer:      consumer,
		portfolioRepo: portfolioRepo,
		positionRepo:  positionRepo,
		cache:         c,
	}
}

// Run blocks, consuming trade events and applying them to portfolio state,
// until ctx is cancelled.
func (p *Processor) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		evt, err := p.consumer.ReadTrade(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			log.Printf("failed to read trade event: %v\n", err)
			continue
		}

		if err := p.processTrade(ctx, evt); err != nil {
			if errors.Is(err, repository.ErrDuplicateTrade) {
				log.Printf("skipped duplicate trade event: %s\n", evt.EventID)
				continue
			}
			log.Printf("failed to process trade event %s: %v\n", evt.EventID, err)
			continue
		}

		log.Printf("processed trade: wallet=%s token=%s side=%s event_id=%s\n", evt.Wallet, evt.Token, evt.Side, evt.EventID)
	}
}

// processTrade applies a single trade event to the wallet's position and
// persists the result.
func (p *Processor) processTrade(ctx context.Context, evt kafka.TradeEvent) error {
	// Day 1's decoder currently emits Wallet/Token as empty and Side as
	// "UNKNOWN" until instruction-level decoding is added — guard against
	// processing incomplete events so we don't corrupt position data.
	if evt.Wallet == "" || evt.Token == "" || evt.Side == "UNKNOWN" {
		return nil
	}

	if err := p.portfolioRepo.EnsureWallet(ctx, evt.Wallet); err != nil {
		return err
	}
	if err := p.portfolioRepo.EnsureToken(ctx, evt.Token); err != nil {
		return err
	}

	txID, err := p.portfolioRepo.RecordTransaction(ctx, evt.Signature, evt.Wallet, true)
	if err != nil {
		return err
	}

	if err := p.portfolioRepo.RecordTrade(ctx, evt.EventID, txID, evt.Wallet, evt.Token, evt.Side, evt.Quantity, evt.Price, evt.Timestamp, evt.Signature); err != nil {
		return err
	}

	pos, exists, err := p.positionRepo.GetOne(ctx, evt.Wallet, evt.Token)
	if err != nil {
		return err
	}
	if !exists {
		pos = domain.Position{Wallet: evt.Wallet, Token: evt.Token}
	}

	switch domain.TradeSide(evt.Side) {
	case domain.SideBuy:
		pos.ApplyBuy(evt.Quantity, evt.Price)
	case domain.SideSell:
		pos.ApplySell(evt.Quantity)
	}

	if err := p.positionRepo.Upsert(ctx, pos); err != nil {
		return err
	}

	return p.cache.InvalidatePortfolio(ctx, evt.Wallet)
}
