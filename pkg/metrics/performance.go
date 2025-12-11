package metrics

import (
	"sync"
	"time"

	"go.uber.org/zap"
)

// PerformanceMetrics coleta métricas básicas de performance
type PerformanceMetrics struct {
	mu sync.RWMutex

	// Contadores de mensagens Email
	MessagesProcessed     int64
	MessagesSuccessful    int64
	MessagesFailed        int64
	MessagesInvalidEmail  int64

	// Tempos de processamento
	TotalProcessTime    time.Duration
	AverageProcessTime  time.Duration
	LastProcessTime     time.Time

	// Performance por período
	MessagesPerSecond   float64
	MessagesPerMinute   float64
	LastCalculationTime time.Time

	// Queries Oracle
	QueriesExecuted  int64
	TotalQueryTime   time.Duration
	AverageQueryTime time.Duration

	// Envios Email
	EmailSendsAttempted int64
	EmailSendsSuccess   int64
	TotalSendTime       time.Duration
	AverageSendTime     time.Duration

	// Custo (se aplicável)
	TotalCost float64 // em centavos
}

// NewPerformanceMetrics cria uma nova instância de métricas
func NewPerformanceMetrics() *PerformanceMetrics {
	return &PerformanceMetrics{
		LastCalculationTime: time.Now(),
		LastProcessTime:     time.Now(),
	}
}

// RecordMessageProcessed registra uma mensagem Email processada
func (pm *PerformanceMetrics) RecordMessageProcessed(success bool, invalidEmail bool, processDuration time.Duration) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.MessagesProcessed++
	pm.TotalProcessTime += processDuration

	if success {
		pm.MessagesSuccessful++
	} else {
		pm.MessagesFailed++
		if invalidEmail {
			pm.MessagesInvalidEmail++
		}
	}

	pm.LastProcessTime = time.Now()
	pm.calculateAverages()
}

// RecordQueryExecuted registra uma query executada
func (pm *PerformanceMetrics) RecordQueryExecuted(duration time.Duration) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.QueriesExecuted++
	pm.TotalQueryTime += duration

	if pm.QueriesExecuted > 0 {
		pm.AverageQueryTime = pm.TotalQueryTime / time.Duration(pm.QueriesExecuted)
	}
}

// RecordEmailSend registra um envio Email
func (pm *PerformanceMetrics) RecordEmailSend(success bool, duration time.Duration, cost float64) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.EmailSendsAttempted++
	pm.TotalSendTime += duration

	if success {
		pm.EmailSendsSuccess++
		pm.TotalCost += cost
	}

	if pm.EmailSendsAttempted > 0 {
		pm.AverageSendTime = pm.TotalSendTime / time.Duration(pm.EmailSendsAttempted)
	}
}

// calculateAverages calcula médias e taxas (deve ser chamado com lock)
func (pm *PerformanceMetrics) calculateAverages() {
	if pm.MessagesProcessed > 0 {
		pm.AverageProcessTime = pm.TotalProcessTime / time.Duration(pm.MessagesProcessed)
	}

	// Calcular taxa por segundo/minuto se passou tempo suficiente
	now := time.Now()
	elapsed := now.Sub(pm.LastCalculationTime)

	if elapsed >= 10*time.Second { // Recalcular a cada 10 segundos
		elapsedSeconds := elapsed.Seconds()

		if elapsedSeconds > 0 {
			pm.MessagesPerSecond = float64(pm.MessagesProcessed) / elapsedSeconds
			pm.MessagesPerMinute = pm.MessagesPerSecond * 60
		}
	}
}

// GetSnapshot retorna um snapshot das métricas atuais
func (pm *PerformanceMetrics) GetSnapshot() MetricsSnapshot {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	successRate := float64(0)
	if pm.MessagesProcessed > 0 {
		successRate = float64(pm.MessagesSuccessful) / float64(pm.MessagesProcessed) * 100
	}

	emailSuccessRate := float64(0)
	if pm.EmailSendsAttempted > 0 {
		emailSuccessRate = float64(pm.EmailSendsSuccess) / float64(pm.EmailSendsAttempted) * 100
	}

	return MetricsSnapshot{
		MessagesProcessed:    pm.MessagesProcessed,
		MessagesSuccessful:   pm.MessagesSuccessful,
		MessagesFailed:       pm.MessagesFailed,
		MessagesInvalidEmail: pm.MessagesInvalidEmail,
		SuccessRate:          successRate,

		AverageProcessTime: pm.AverageProcessTime,
		MessagesPerSecond:  pm.MessagesPerSecond,
		MessagesPerMinute:  pm.MessagesPerMinute,

		QueriesExecuted:  pm.QueriesExecuted,
		AverageQueryTime: pm.AverageQueryTime,

		EmailSendsAttempted: pm.EmailSendsAttempted,
		EmailSendsSuccess:   pm.EmailSendsSuccess,
		EmailSuccessRate:    emailSuccessRate,
		AverageSendTime:     pm.AverageSendTime,

		TotalCost: pm.TotalCost,

		LastProcessTime: pm.LastProcessTime,
	}
}

// LogMetrics registra as métricas no logger
func (pm *PerformanceMetrics) LogMetrics(logger *zap.Logger) {
	snapshot := pm.GetSnapshot()

	logger.Info("📊 Métricas de Performance",
		zap.Int64("mensagens_processadas", snapshot.MessagesProcessed),
		zap.Int64("sucessos", snapshot.MessagesSuccessful),
		zap.Int64("falhas", snapshot.MessagesFailed),
		zap.Int64("emails_invalidos", snapshot.MessagesInvalidEmail),
		zap.Float64("taxa_sucesso_pct", snapshot.SuccessRate),
		zap.Duration("tempo_medio_processamento", snapshot.AverageProcessTime),
		zap.Float64("emails_por_segundo", snapshot.MessagesPerSecond),
		zap.Float64("emails_por_minuto", snapshot.MessagesPerMinute),
		zap.Int64("queries_executadas", snapshot.QueriesExecuted),
		zap.Duration("tempo_medio_query", snapshot.AverageQueryTime),
		zap.Int64("email_tentativas", snapshot.EmailSendsAttempted),
		zap.Int64("email_sucessos", snapshot.EmailSendsSuccess),
		zap.Float64("email_taxa_sucesso_pct", snapshot.EmailSuccessRate),
		zap.Duration("tempo_medio_email", snapshot.AverageSendTime),
		zap.Float64("custo_total_centavos", snapshot.TotalCost),
		zap.Float64("custo_total_reais", snapshot.TotalCost/100),
	)
}

// Reset limpa todas as métricas
func (pm *PerformanceMetrics) Reset() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.MessagesProcessed = 0
	pm.MessagesSuccessful = 0
	pm.MessagesFailed = 0
	pm.MessagesInvalidEmail = 0
	pm.TotalProcessTime = 0
	pm.AverageProcessTime = 0
	pm.MessagesPerSecond = 0
	pm.MessagesPerMinute = 0
	pm.QueriesExecuted = 0
	pm.TotalQueryTime = 0
	pm.AverageQueryTime = 0
	pm.EmailSendsAttempted = 0
	pm.EmailSendsSuccess = 0
	pm.TotalSendTime = 0
	pm.AverageSendTime = 0
	pm.TotalCost = 0
	pm.LastCalculationTime = time.Now()
	pm.LastProcessTime = time.Now()
}

// MetricsSnapshot representa um snapshot das métricas em um momento específico
type MetricsSnapshot struct {
	MessagesProcessed    int64
	MessagesSuccessful   int64
	MessagesFailed       int64
	MessagesInvalidEmail int64
	SuccessRate          float64

	AverageProcessTime time.Duration
	MessagesPerSecond  float64
	MessagesPerMinute  float64

	QueriesExecuted  int64
	AverageQueryTime time.Duration

	EmailSendsAttempted int64
	EmailSendsSuccess   int64
	EmailSuccessRate    float64
	AverageSendTime     time.Duration

	TotalCost float64

	LastProcessTime time.Time
}

// Stats representa estatísticas detalhadas para o dashboard
type Stats struct {
	TotalMessagesProcessed int64
	SuccessCount           int64
	ErrorCount             int64
	InvalidTokenCount      int64
	AvgProcessTimeMs       float64
	AvgSendTimeMs          float64
	AvgQueryTimeMs         float64
	QueriesExecuted        int64
	PushSendSuccessCount   int64 // Mantido por compatibilidade com dashboard
	PushSendErrorCount     int64 // Mantido por compatibilidade com dashboard
}

// GetStats retorna estatísticas formatadas para o dashboard
func (pm *PerformanceMetrics) GetStats() Stats {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return Stats{
		TotalMessagesProcessed: pm.MessagesProcessed,
		SuccessCount:           pm.MessagesSuccessful,
		ErrorCount:             pm.MessagesFailed,
		InvalidTokenCount:      pm.MessagesInvalidEmail,
		AvgProcessTimeMs:       float64(pm.AverageProcessTime.Microseconds()) / 1000.0,
		AvgSendTimeMs:          float64(pm.AverageSendTime.Microseconds()) / 1000.0,
		AvgQueryTimeMs:         float64(pm.AverageQueryTime.Microseconds()) / 1000.0,
		QueriesExecuted:        pm.QueriesExecuted,
		PushSendSuccessCount:   pm.EmailSendsSuccess,
		PushSendErrorCount:     pm.EmailSendsAttempted - pm.EmailSendsSuccess,
	}
}
