package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/logger"
	"go.uber.org/zap"
)

// ConfigureOracle configura a conexão Oracle com parâmetros otimizados para performance
func ConfigureOracle(db *sql.DB) error {
	logger.Info("Configurando parâmetros otimizados da conexão Oracle")

	// Pool de conexões otimizado para maior throughput
	db.SetMaxOpenConns(20)                  // Aumentado de 10 para 20
	db.SetMaxIdleConns(10)                  // Aumentado de 5 para 10
	db.SetConnMaxLifetime(30 * time.Minute) // Reduzido de 1h para 30min para rotação mais frequente

	// Testar conexão com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		logger.Error("Erro ao testar conexão Oracle", zap.Error(err))
		return fmt.Errorf("erro ao testar conexão: %v", err)
	}

	logger.Info("Conexão Oracle configurada com sucesso")
	return nil
}
