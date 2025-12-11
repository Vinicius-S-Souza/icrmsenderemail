package health

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type HealthStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Database  string    `json:"database"`
	Uptime    string    `json:"uptime"`
}

type HealthChecker struct {
	db        *sql.DB
	logger    *zap.Logger
	startTime time.Time
}

func NewHealthChecker(db *sql.DB, logger *zap.Logger) *HealthChecker {
	return &HealthChecker{
		db:        db,
		logger:    logger,
		startTime: time.Now(),
	}
}

func (h *HealthChecker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := h.Check()

	w.Header().Set("Content-Type", "application/json")
	if status.Status != "healthy" {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	json.NewEncoder(w).Encode(status)
}

func (h *HealthChecker) Check() HealthStatus {
	status := HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Database:  "disconnected",
		Uptime:    time.Since(h.startTime).String(),
	}

	// Verificar conexão com banco de dados
	if h.db != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := h.db.PingContext(ctx); err != nil {
			status.Status = "unhealthy"
			status.Database = fmt.Sprintf("error: %v", err)
		} else {
			status.Database = "connected"
		}
	}

	return status
}

// StartHealthServer inicia o servidor de health check
func StartHealthServer(port int, checker *HealthChecker, logger *zap.Logger) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/health", checker)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		logger.Info("Servidor de health check iniciado", zap.Int("port", port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Erro no servidor de health check", zap.Error(err))
		}
	}()

	return server
}
