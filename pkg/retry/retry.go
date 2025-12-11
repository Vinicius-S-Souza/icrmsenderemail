package retry

import (
	"context"
	"fmt"
	"math"
	"time"

	"go.uber.org/zap"
)

// Config configurações para retry com exponential backoff
type Config struct {
	MaxAttempts     int           // Número máximo de tentativas
	InitialInterval time.Duration // Intervalo inicial entre tentativas
	MaxInterval     time.Duration // Intervalo máximo entre tentativas
	Multiplier      float64       // Multiplicador para backoff exponencial
	MaxElapsedTime  time.Duration // Tempo máximo total para todas as tentativas
}

// DefaultConfig retorna configuração padrão
func DefaultConfig() Config {
	return Config{
		MaxAttempts:     3,
		InitialInterval: 1 * time.Second,
		MaxInterval:     30 * time.Second,
		Multiplier:      2.0,
		MaxElapsedTime:  2 * time.Minute,
	}
}

// Operation função que será executada com retry
type Operation func() error

// RetryableError interface para erros que devem ser retentados
type RetryableError interface {
	error
	IsRetryable() bool
}

// Retry executa uma operação com exponential backoff
func Retry(ctx context.Context, config Config, operation Operation, logger *zap.Logger) error {
	if config.MaxAttempts <= 0 {
		return fmt.Errorf("max_attempts deve ser maior que 0")
	}

	startTime := time.Now()
	var lastErr error

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		// Verificar timeout de contexto
		select {
		case <-ctx.Done():
			return fmt.Errorf("contexto cancelado: %w", ctx.Err())
		default:
		}

		// Verificar se excedeu tempo máximo
		if config.MaxElapsedTime > 0 && time.Since(startTime) >= config.MaxElapsedTime {
			return fmt.Errorf("tempo máximo excedido após %d tentativas: %w", attempt-1, lastErr)
		}

		// Executar operação
		err := operation()
		if err == nil {
			// Sucesso
			if attempt > 1 && logger != nil {
				logger.Info("Operação bem-sucedida após retry",
					zap.Int("tentativa", attempt),
					zap.Duration("tempo_total", time.Since(startTime)))
			}
			return nil
		}

		lastErr = err

		// Verificar se erro é retentável
		if retryableErr, ok := err.(RetryableError); ok && !retryableErr.IsRetryable() {
			// Erro não retentável - falhar imediatamente
			if logger != nil {
				logger.Warn("Erro não retentável detectado",
					zap.Error(err),
					zap.Int("tentativa", attempt))
			}
			return err
		}

		// Se é a última tentativa, não faz backoff
		if attempt >= config.MaxAttempts {
			if logger != nil {
				logger.Error("Máximo de tentativas atingido",
					zap.Error(err),
					zap.Int("tentativas", attempt),
					zap.Duration("tempo_total", time.Since(startTime)))
			}
			return fmt.Errorf("após %d tentativas: %w", attempt, err)
		}

		// Calcular tempo de backoff
		backoff := calculateBackoff(attempt, config)

		if logger != nil {
			logger.Warn("Tentativa falhou, aguardando antes de retentar",
				zap.Error(err),
				zap.Int("tentativa", attempt),
				zap.Int("max_tentativas", config.MaxAttempts),
				zap.Duration("backoff", backoff))
		}

		// Aguardar com backoff ou até contexto ser cancelado
		select {
		case <-ctx.Done():
			return fmt.Errorf("contexto cancelado durante backoff: %w", ctx.Err())
		case <-time.After(backoff):
			// Continuar para próxima tentativa
		}
	}

	return lastErr
}

// calculateBackoff calcula o tempo de espera baseado em exponential backoff
func calculateBackoff(attempt int, config Config) time.Duration {
	// Backoff = InitialInterval * (Multiplier ^ (attempt - 1))
	backoff := float64(config.InitialInterval) * math.Pow(config.Multiplier, float64(attempt-1))

	// Limitar ao intervalo máximo
	if backoff > float64(config.MaxInterval) {
		backoff = float64(config.MaxInterval)
	}

	return time.Duration(backoff)
}

// Do é um alias conveniente para Retry com configuração padrão
func Do(ctx context.Context, operation Operation, logger *zap.Logger) error {
	return Retry(ctx, DefaultConfig(), operation, logger)
}

// DoWithConfig executa retry com configuração customizada
func DoWithConfig(ctx context.Context, config Config, operation Operation, logger *zap.Logger) error {
	return Retry(ctx, config, operation, logger)
}
