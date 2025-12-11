package message

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/config"
	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/email"
	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/metrics"
	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/retry"
	"go.uber.org/zap"
)

// Processor processa mensagens de email
type Processor struct {
	repo        *Repository
	sender      *email.Sender
	metrics     *metrics.PerformanceMetrics
	config      *config.PerformanceConfig
	logger      *zap.Logger
	defaultFrom string // Remetente padrão configurado

	ctx            context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup
	jobQueue       chan *Email
	isRunning      bool
	mu             sync.Mutex
	circuitBreaker *CircuitBreaker
}

// CircuitBreaker proteção contra falhas em cascata
type CircuitBreaker struct {
	mu              sync.RWMutex
	failureCount    int
	lastFailureTime time.Time
	state           string // "closed", "open", "half-open"
	threshold       int
	timeout         time.Duration
}

// NewProcessor cria um novo processador
func NewProcessor(
	repo *Repository,
	sender *email.Sender,
	metricsCollector *metrics.PerformanceMetrics,
	config *config.PerformanceConfig,
	defaultFrom string,
	logger *zap.Logger,
) *Processor {
	ctx, cancel := context.WithCancel(context.Background())

	return &Processor{
		repo:        repo,
		sender:      sender,
		metrics:     metricsCollector,
		config:      config,
		defaultFrom: defaultFrom,
		logger:      logger,
		ctx:         ctx,
		cancel:      cancel,
		jobQueue:    make(chan *Email, config.BatchSize*2),
		isRunning:   false,
		circuitBreaker: &CircuitBreaker{
			state:     "closed",
			threshold: config.CircuitBreakerThreshold,
			timeout:   time.Duration(config.CircuitBreakerTimeoutSeconds) * time.Second,
		},
	}
}

// Start inicia o processamento
func (p *Processor) Start() error {
	p.mu.Lock()
	if p.isRunning {
		p.mu.Unlock()
		return fmt.Errorf("processador já está em execução")
	}
	p.isRunning = true
	p.mu.Unlock()

	p.logger.Info("Iniciando processador de Email",
		zap.Int("workers", p.config.WorkerCount),
		zap.Int("batch_size", p.config.BatchSize),
		zap.Int("circuit_breaker_threshold", p.config.CircuitBreakerThreshold))

	// Iniciar workers
	for i := 0; i < p.config.WorkerCount; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}

	// Iniciar dispatcher
	p.wg.Add(1)
	go p.dispatcher()

	return nil
}

// Stop para o processamento gracefully
func (p *Processor) Stop() error {
	p.mu.Lock()
	if !p.isRunning {
		p.mu.Unlock()
		return fmt.Errorf("processador não está em execução")
	}
	p.mu.Unlock()

	p.logger.Info("Parando processador de Email...")

	// Cancelar contexto
	p.cancel()

	// Aguardar workers finalizarem com timeout
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	shutdownTimeout := time.Duration(30) * time.Second
	select {
	case <-done:
		p.logger.Info("Processador parado com sucesso")
	case <-time.After(shutdownTimeout):
		p.logger.Warn("Timeout ao parar processador",
			zap.Duration("timeout", shutdownTimeout))
	}

	p.mu.Lock()
	p.isRunning = false
	p.mu.Unlock()

	return nil
}

// dispatcher busca emails pendentes e distribui para workers
func (p *Processor) dispatcher() {
	defer p.wg.Done()

	fetchInterval := time.Duration(p.config.FetchIntervalSeconds) * time.Second
	ticker := time.NewTicker(fetchInterval)
	defer ticker.Stop()

	p.logger.Info("Dispatcher iniciado", zap.Duration("fetch_interval", fetchInterval))

	for {
		select {
		case <-p.ctx.Done():
			p.logger.Info("Dispatcher finalizado")
			close(p.jobQueue)
			return

		case <-ticker.C:
			// Circuit breaker check
			if !p.circuitBreaker.canExecute() {
				p.logger.Warn("Circuit breaker aberto - pulando busca de emails")
				continue
			}

			p.fetchAndDispatch()
		}
	}
}

// fetchAndDispatch busca emails e envia para fila
func (p *Processor) fetchAndDispatch() {
	ctx, cancel := context.WithTimeout(p.ctx, 30*time.Second)
	defer cancel()

	// Medir tempo de execução da query
	queryStartTime := time.Now()
	emailList, err := p.repo.GetPendingEmails(ctx, p.config.BatchSize, p.config.DataDisparoOffset, p.config.MaxTentativas)
	queryDuration := time.Since(queryStartTime)

	// Registrar métrica de query executada
	p.metrics.RecordQueryExecuted(queryDuration)

	if err != nil {
		p.logger.Error("Erro ao buscar emails pendentes", zap.Error(err))
		p.circuitBreaker.recordFailure()
		return
	}

	if len(emailList) == 0 {
		p.logger.Debug("Nenhum email pendente encontrado")
		p.circuitBreaker.recordSuccess()
		return
	}

	p.logger.Info("Emails pendentes encontrados",
		zap.Int("total", len(emailList)))

	// Enviar para fila de processamento
	for i := range emailList {
		select {
		case <-p.ctx.Done():
			return
		case p.jobQueue <- &emailList[i]:
			// Email adicionado à fila
		default:
			// Fila cheia, processar no próximo ciclo
			p.logger.Warn("Fila de processamento cheia, aguardando próximo ciclo")
			return
		}
	}

	p.circuitBreaker.recordSuccess()
}

// worker processa emails da fila
func (p *Processor) worker(id int) {
	defer p.wg.Done()

	p.logger.Info("Worker iniciado", zap.Int("worker_id", id))

	for {
		select {
		case <-p.ctx.Done():
			p.logger.Info("Worker finalizado", zap.Int("worker_id", id))
			return

		case emailMsg, ok := <-p.jobQueue:
			if !ok {
				p.logger.Info("Canal de jobs fechado", zap.Int("worker_id", id))
				return
			}

			p.processEmail(emailMsg, id)
		}
	}
}

// processEmail processa um único email
func (p *Processor) processEmail(message *Email, workerID int) {
	startTime := time.Now()

	p.logger.Info("Processando email",
		zap.Int("worker_id", workerID),
		zap.Int64("email_id", message.ID),
		zap.String("to", maskEmail(message.Destinatario)),
		zap.Int("tentativa", message.QTDTentativas+1))

	// Criar contexto com timeout para envio
	ctx, cancel := context.WithTimeout(p.ctx, time.Duration(p.config.SendTimeoutSeconds)*time.Second)
	defer cancel()

	// Configurar retry
	retryConfig := retry.Config{
		MaxAttempts:     2,
		InitialInterval: 1 * time.Second,
		MaxInterval:     5 * time.Second,
		Multiplier:      2.0,
		MaxElapsedTime:  time.Duration(p.config.SendTimeoutSeconds) * time.Second,
	}

	// Preparar dados para envio
	emailData := email.EmailData{
		ID:          message.ID,
		From:        p.defaultFrom, // Usar o remetente padrão da configuração
		To:          message.Destinatario,
		Subject:     message.Assunto,
		Body:        message.Corpo,
		ContentType: message.TipoCorpo,
	}

	// Carregar anexo se houver
	if message.AnexoReferencia.Valid && message.AnexoReferencia.String != "" {
		// Verificar se é URL (Zenvia) ou base64 (SendGrid, Pontaltech)
		if message.AnexoTipo.Valid && message.AnexoTipo.String == "url" {
			// Anexo via URL pública (Zenvia)
			emailData.Attachment = &email.Attachment{
				Filename:    message.AnexoNome.String,
				ContentType: "application/octet-stream", // Tipo genérico para URL
				URL:         message.AnexoReferencia.String, // URL pública
			}
			p.logger.Debug("Anexo via URL carregado",
				zap.Int64("email_id", message.ID),
				zap.String("filename", message.AnexoNome.String),
				zap.String("url", message.AnexoReferencia.String))
		} else {
			// Anexo em base64 (SendGrid, Pontaltech)
			anexoData, err := base64.StdEncoding.DecodeString(message.AnexoReferencia.String)
			if err != nil {
				p.logger.Error("Erro ao decodificar anexo",
					zap.Int64("email_id", message.ID),
					zap.Error(err))
			} else {
				emailData.Attachment = &email.Attachment{
					Filename:    message.AnexoNome.String,
					ContentType: message.AnexoTipo.String,
					Data:        bytes.NewReader(anexoData),
					Size:        int64(len(anexoData)),
				}
				p.logger.Debug("Anexo em base64 carregado",
					zap.Int64("email_id", message.ID),
					zap.String("filename", message.AnexoNome.String),
					zap.Int("size_bytes", len(anexoData)))
			}
		}
	}

	var result email.SendResult
	err := retry.Retry(ctx, retryConfig, func() error {
		result = p.sender.Send(ctx, emailData)
		if !result.Success {
			return result.Error
		}
		return nil
	}, p.logger)

	// Processar resultado
	if err == nil && result.Success {
		// Sucesso
		providerName := p.sender.GetProvider().GetName()
		providerCode := ProviderStringToCode(providerName)
		if err := p.repo.MarkAsSent(ctx, message.ID, result.ProviderID, providerCode); err != nil {
			p.logger.Error("Erro ao marcar email como enviado", zap.Error(err))
			p.metrics.RecordMessageProcessed(false, false, time.Since(startTime))
		} else {
			sendDuration := time.Since(startTime)
			p.metrics.RecordMessageProcessed(true, false, sendDuration)
			p.metrics.RecordEmailSend(true, sendDuration, 0)
			p.logger.Info("Email enviado com sucesso",
				zap.Int64("email_id", message.ID),
				zap.String("provider_id", result.ProviderID),
				zap.String("provider", providerName),
				zap.Int("provider_code", providerCode),
				zap.Duration("duracao", sendDuration))
		}
	} else {
		// Erro
		processDuration := time.Since(startTime)
		providerName := p.sender.GetProvider().GetName()
		providerCode := ProviderStringToCode(providerName)

		errorMsg := "erro desconhecido"
		if result.Error != nil {
			errorMsg = result.Error.Error()
		} else if err != nil {
			errorMsg = err.Error()
		}

		// Determinar tipo de erro
		isInvalidEmail := isInvalidEmailError(errorMsg)
		if isInvalidEmail {
			if err := p.repo.MarkAsInvalid(ctx, message.ID, errorMsg, providerCode); err != nil {
				p.logger.Error("Erro ao marcar email como inválido", zap.Error(err))
			}
			p.metrics.RecordMessageProcessed(false, true, processDuration)
		} else if message.QTDTentativas+1 >= p.config.MaxTentativas {
			if err := p.repo.MarkAsPermanentFailure(ctx, message.ID, errorMsg, providerCode); err != nil {
				p.logger.Error("Erro ao marcar email como falha permanente", zap.Error(err))
			}
			p.metrics.RecordMessageProcessed(false, false, processDuration)
		} else {
			// Erro temporário, marcar para retry
			if err := p.repo.MarkAsError(ctx, message.ID, errorMsg, providerCode); err != nil {
				p.logger.Error("Erro ao marcar email com erro", zap.Error(err))
			}
			p.metrics.RecordMessageProcessed(false, false, processDuration)
		}

		p.metrics.RecordEmailSend(false, processDuration, 0)

		p.logger.Warn("Falha ao enviar email",
			zap.Int64("email_id", message.ID),
			zap.String("provider", providerName),
			zap.Int("provider_code", providerCode),
			zap.String("erro", errorMsg),
			zap.Int("tentativas", message.QTDTentativas+1))
	}
}

// isInvalidEmailError verifica se erro é de email inválido
func isInvalidEmailError(errMsg string) bool {
	invalidPhrases := []string{
		"email inválido",
		"invalid email",
		"invalid address",
		"invalid recipient",
		"malformed email",
	}

	errLower := strings.ToLower(errMsg)
	for _, phrase := range invalidPhrases {
		if strings.Contains(errLower, phrase) {
			return true
		}
	}

	return false
}

// maskEmail mascara email para logs
func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "***"
	}
	if len(parts[0]) <= 2 {
		return "**@" + parts[1]
	}
	return parts[0][:2] + "***@" + parts[1]
}

// GetMetrics retorna as métricas do processador
func (p *Processor) GetMetrics() *metrics.PerformanceMetrics {
	return p.metrics
}

// IsRunning verifica se processador está rodando
func (p *Processor) IsRunning() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.isRunning
}

// Circuit breaker methods
func (cb *CircuitBreaker) canExecute() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	if cb.state == "closed" {
		return true
	}

	if cb.state == "open" {
		if time.Since(cb.lastFailureTime) > cb.timeout {
			cb.mu.RUnlock()
			cb.mu.Lock()
			cb.state = "half-open"
			cb.mu.Unlock()
			cb.mu.RLock()
			return true
		}
		return false
	}

	// half-open state
	return true
}

func (cb *CircuitBreaker) recordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failureCount = 0
	cb.state = "closed"
}

func (cb *CircuitBreaker) recordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failureCount++
	cb.lastFailureTime = time.Now()

	if cb.failureCount >= cb.threshold {
		cb.state = "open"
	}
}

// GetCircuitBreakerState retorna estado atual do circuit breaker
func (p *Processor) GetCircuitBreakerState() string {
	p.circuitBreaker.mu.RLock()
	defer p.circuitBreaker.mu.RUnlock()
	return p.circuitBreaker.state
}
