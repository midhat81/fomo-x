package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/midhat81/fomo-x/services/trading-engine/internal/config"
	"github.com/midhat81/fomo-x/services/trading-engine/internal/domain"
	"github.com/midhat81/fomo-x/services/trading-engine/internal/execution"
	"github.com/midhat81/fomo-x/services/trading-engine/internal/repository"
	"github.com/midhat81/fomo-x/services/trading-engine/internal/risk"
)

type createOrderRequest struct {
	IdempotencyKey string  `json:"idempotency_key"`
	Wallet         string  `json:"wallet"`
	Token          string  `json:"token"`
	Side           string  `json:"side"`
	Quantity       float64 `json:"quantity"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on system environment variables")
	}

	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := pgxpool.New(ctx, cfg.PostgresURL)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("postgres ping failed: %v", err)
	}
	log.Println("Connected to Postgres.")

	orderRepo := repository.NewOrderRepository(pool)
	riskEngine := risk.NewEngine(domain.DefaultRiskLimits())

	mux := http.NewServeMux()

	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req createOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		order := domain.NewOrder(req.IdempotencyKey, req.Wallet, req.Token, domain.OrderSide(req.Side), req.Quantity)

		created, err := orderRepo.Create(r.Context(), order)
		if err != nil {
			if err == repository.ErrDuplicateOrder {
				log.Printf("duplicate order detected for idempotency_key=%s, returning existing order\n", req.IdempotencyKey)
				respondJSON(w, http.StatusOK, created)
				return
			}
			log.Printf("failed to create order: %v\n", err)
			http.Error(w, "failed to create order", http.StatusInternalServerError)
			return
		}

		marketPrice := execution.GetMarketPrice(created.Token)

		account := risk.AccountState{Wallet: created.Wallet}
		riskResult := riskEngine.Evaluate(created, marketPrice, account)

		if !riskResult.Passed {
			created.Reject(riskResult.Reason)
			_ = orderRepo.UpdateStatus(r.Context(), created.ID, created.Status, created.RejectReason)
			log.Printf("order %s rejected: %s\n", created.ID, riskResult.Reason)
			respondJSON(w, http.StatusOK, created)
			return
		}

		created.Accept()
		_ = orderRepo.UpdateStatus(r.Context(), created.ID, created.Status, "")

		exec := execution.Simulate(created, marketPrice)
		_ = orderRepo.RecordExecution(r.Context(), created.ID, exec)

		if exec.IsFullyFilled(created.Quantity) {
			created.Fill()
		} else {
			created.Status = domain.StatusPartiallyFilled
		}
		_ = orderRepo.UpdateStatus(r.Context(), created.ID, created.Status, "")

		log.Printf("order %s executed: status=%s filled=%v price=%v\n", created.ID, created.Status, exec.FilledQty, exec.FillPrice)

		respondJSON(w, http.StatusOK, created)
	})

	server := &http.Server{Addr: ":3002", Handler: mux}

	go func() {
		<-ctx.Done()
		log.Println("Shutting down trading engine...")
		_ = server.Close()
	}()

	log.Println("FOMO-X Trading Engine listening on port 3002")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}

	log.Println("Trading engine shut down cleanly.")
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
