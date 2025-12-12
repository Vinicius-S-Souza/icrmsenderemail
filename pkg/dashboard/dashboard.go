package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/metrics"
	"go.uber.org/zap"
)

// MessageRepository interface para acessar dados de mensagens
type MessageRepository interface {
	CountPendingEmails(ctx context.Context, daysOffset int, maxTentativas int) (int64, error)
}

// Dashboard gerencia o servidor web de métricas
type Dashboard struct {
	server         *http.Server
	logger         *zap.Logger
	metricsSource  *metrics.PerformanceMetrics
	messageRepo    MessageRepository
	daysOffset     int
	maxTentativas  int
	mu             sync.RWMutex
	clients         map[chan []byte]bool
	port            int
	providerName    string
	mux             *http.ServeMux
	manualHandler   ManualHandler
	templateHandler TemplateHandler
}

// Config contém as configurações do dashboard
type Config struct {
	Port            int
	EnableDashboard bool
	ProviderName    string
	DaysOffset      int
	MaxTentativas   int
}

// ManualHandler interface para handlers de disparo manual
type ManualHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	ValidarCliente(w http.ResponseWriter, r *http.Request)
	DispararEmail(w http.ResponseWriter, r *http.Request)
	ConsultarStatus(w http.ResponseWriter, r *http.Request)
	GetProviderInfo(w http.ResponseWriter, r *http.Request)
	PreviewTemplate(w http.ResponseWriter, r *http.Request)
}

// TemplateHandler interface para handlers de templates
type TemplateHandler interface {
	ServeTemplateList(w http.ResponseWriter, r *http.Request)
	ServeTemplateEditor(w http.ResponseWriter, r *http.Request)
	ListTemplates(w http.ResponseWriter, r *http.Request)
	GetTemplate(w http.ResponseWriter, r *http.Request)
	CreateTemplate(w http.ResponseWriter, r *http.Request)
	UpdateTemplate(w http.ResponseWriter, r *http.Request)
	DeleteTemplate(w http.ResponseWriter, r *http.Request)
	GetMacros(w http.ResponseWriter, r *http.Request)
	PreviewTemplate(w http.ResponseWriter, r *http.Request)
	DuplicateTemplate(w http.ResponseWriter, r *http.Request)
}

// MetricsSnapshot representa um snapshot das métricas
type MetricsSnapshot struct {
	Timestamp              time.Time `json:"timestamp"`
	ProviderName           string    `json:"provider_name"`
	TotalMessagesProcessed int64     `json:"total_messages_processed"`
	SuccessCount           int64     `json:"success_count"`
	ErrorCount             int64     `json:"error_count"`
	InvalidEmailCount      int64     `json:"invalid_email_count"`
	SuccessRate            float64   `json:"success_rate"`
	ErrorRate              float64   `json:"error_rate"`
	InvalidEmailRate       float64   `json:"invalid_email_rate"`
	AvgProcessTime         float64   `json:"avg_process_time_ms"`
	AvgSendTime            float64   `json:"avg_send_time_ms"`
	AvgQueryTime           float64   `json:"avg_query_time_ms"`
	QueriesExecuted        int64     `json:"queries_executed"`
	EmailSendSuccessCount  int64     `json:"email_send_success_count"`
	EmailSendErrorCount    int64     `json:"email_send_error_count"`
	EmailSendSuccessRate   float64   `json:"email_send_success_rate"`
	PendingMessagesCount   int64     `json:"pending_messages_count"`
}

// NewDashboard cria uma nova instância do dashboard
func NewDashboard(config Config, metricsSource *metrics.PerformanceMetrics, messageRepo MessageRepository, logger *zap.Logger) *Dashboard {
	return &Dashboard{
		logger:        logger,
		metricsSource: metricsSource,
		messageRepo:   messageRepo,
		daysOffset:    config.DaysOffset,
		maxTentativas: config.MaxTentativas,
		clients:       make(map[chan []byte]bool),
		port:          config.Port,
		providerName:  config.ProviderName,
	}
}

// RegisterManualEndpoints registra os endpoints de disparo manual
func (d *Dashboard) RegisterManualEndpoints(handler ManualHandler) {
	d.manualHandler = handler
}

// RegisterTemplateEndpoints registra os endpoints de templates
func (d *Dashboard) RegisterTemplateEndpoints(handler TemplateHandler) {
	d.templateHandler = handler
}

// Start inicia o servidor do dashboard
func (d *Dashboard) Start() error {
	d.mux = http.NewServeMux()

	// Endpoints da API de métricas
	d.mux.HandleFunc("/api/metrics", d.handleMetrics)
	d.mux.HandleFunc("/api/metrics/stream", d.handleMetricsStream)

	// Endpoints de disparo manual (se configurado)
	if d.manualHandler != nil {
		d.mux.HandleFunc("/manual", d.manualHandler.ServeHTTP)
		d.mux.HandleFunc("/api/manual/validar-cliente", d.manualHandler.ValidarCliente)
		d.mux.HandleFunc("/api/manual/disparar", d.manualHandler.DispararEmail)
		d.mux.HandleFunc("/api/manual/status", d.manualHandler.ConsultarStatus)
		d.mux.HandleFunc("/api/manual/provider-info", d.manualHandler.GetProviderInfo)
		d.mux.HandleFunc("/api/manual/preview-template", d.manualHandler.PreviewTemplate)
	}

	// Endpoints de templates (se configurado)
	if d.templateHandler != nil {
		// Páginas web
		d.mux.HandleFunc("/templates", d.templateHandler.ServeTemplateList)
		d.mux.HandleFunc("/templates/novo", d.templateHandler.ServeTemplateEditor)
		d.mux.HandleFunc("/templates/", d.handleTemplateRoute)
		// API REST - ordem importa! Rotas mais específicas primeiro
		d.mux.HandleFunc("/api/templates/macros", d.templateHandler.GetMacros)
		d.mux.HandleFunc("/api/templates/preview", d.templateHandler.PreviewTemplate)
		d.mux.HandleFunc("/api/templates/", d.handleTemplatesAPIWithID)
		d.mux.HandleFunc("/api/templates", d.handleTemplatesAPI)
	}

	// Servir página principal do dashboard
	d.mux.HandleFunc("/", d.handleIndex)

	d.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", d.port),
		Handler: d.corsMiddleware(d.mux),
	}

	d.logger.Info("Dashboard iniciado", zap.Int("port", d.port))

	// Iniciar broadcaster de métricas
	go d.broadcastMetrics()

	return d.server.ListenAndServe()
}

// Stop para o servidor do dashboard
func (d *Dashboard) Stop(ctx context.Context) error {
	d.logger.Info("Parando dashboard...")
	return d.server.Shutdown(ctx)
}

// handleMetrics retorna as métricas atuais em JSON
func (d *Dashboard) handleMetrics(w http.ResponseWriter, r *http.Request) {
	snapshot := d.getMetricsSnapshot()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snapshot)
}

// handleMetricsStream implementa Server-Sent Events para streaming de métricas
func (d *Dashboard) handleMetricsStream(w http.ResponseWriter, r *http.Request) {
	// Configurar headers para SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Criar canal para este cliente
	clientChan := make(chan []byte, 10)

	d.mu.Lock()
	d.clients[clientChan] = true
	d.mu.Unlock()

	// Remover cliente quando desconectar
	defer func() {
		d.mu.Lock()
		delete(d.clients, clientChan)
		close(clientChan)
		d.mu.Unlock()
	}()

	// Enviar métricas quando disponíveis
	for {
		select {
		case data := <-clientChan:
			fmt.Fprintf(w, "data: %s\n\n", data)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case <-r.Context().Done():
			return
		}
	}
}

// broadcastMetrics envia métricas para todos os clientes conectados
func (d *Dashboard) broadcastMetrics() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		snapshot := d.getMetricsSnapshot()
		data, err := json.Marshal(snapshot)
		if err != nil {
			d.logger.Error("Erro ao serializar métricas", zap.Error(err))
			continue
		}

		d.mu.RLock()
		for clientChan := range d.clients {
			select {
			case clientChan <- data:
			default:
				// Canal cheio, pular este cliente
			}
		}
		d.mu.RUnlock()
	}
}

// getMetricsSnapshot obtém um snapshot das métricas atuais
func (d *Dashboard) getMetricsSnapshot() MetricsSnapshot {
	// Obter estatísticas das métricas
	stats := d.metricsSource.GetStats()

	var successRate, errorRate, invalidEmailRate, emailSendSuccessRate float64
	total := float64(stats.TotalMessagesProcessed)

	if total > 0 {
		successRate = (float64(stats.SuccessCount) / total) * 100
		errorRate = (float64(stats.ErrorCount) / total) * 100
		invalidEmailRate = (float64(stats.InvalidTokenCount) / total) * 100
	}

	totalEmailSends := float64(stats.PushSendSuccessCount + stats.PushSendErrorCount)
	if totalEmailSends > 0 {
		emailSendSuccessRate = (float64(stats.PushSendSuccessCount) / totalEmailSends) * 100
	}

	// Buscar contagem de mensagens pendentes
	var pendingCount int64
	if d.messageRepo != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		count, err := d.messageRepo.CountPendingEmails(ctx, d.daysOffset, d.maxTentativas)
		if err != nil {
			d.logger.Error("Erro ao contar mensagens pendentes", zap.Error(err))
			pendingCount = 0
		} else {
			pendingCount = count
		}
	}

	return MetricsSnapshot{
		Timestamp:              time.Now(),
		ProviderName:           d.providerName,
		TotalMessagesProcessed: stats.TotalMessagesProcessed,
		SuccessCount:           stats.SuccessCount,
		ErrorCount:             stats.ErrorCount,
		InvalidEmailCount:      stats.InvalidTokenCount,
		SuccessRate:            successRate,
		ErrorRate:              errorRate,
		InvalidEmailRate:       invalidEmailRate,
		AvgProcessTime:         stats.AvgProcessTimeMs,
		AvgSendTime:            stats.AvgSendTimeMs,
		AvgQueryTime:           stats.AvgQueryTimeMs,
		QueriesExecuted:        stats.QueriesExecuted,
		EmailSendSuccessCount:  stats.PushSendSuccessCount,
		EmailSendErrorCount:    stats.PushSendErrorCount,
		EmailSendSuccessRate:   emailSendSuccessRate,
		PendingMessagesCount:   pendingCount,
	}
}

// corsMiddleware adiciona headers CORS
func (d *Dashboard) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// handleIndex serve a página principal do dashboard
func (d *Dashboard) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(dashboardHTML))
}

// handleTemplateRoute roteia requisições de templates (GET /templates/:id/editar)
func (d *Dashboard) handleTemplateRoute(w http.ResponseWriter, r *http.Request) {
	// /templates/:id/editar
	if strings.HasSuffix(r.URL.Path, "/editar") {
		d.templateHandler.ServeTemplateEditor(w, r)
		return
	}
	http.NotFound(w, r)
}

// handleTemplatesAPI roteia requisições da API de templates (sem ID)
func (d *Dashboard) handleTemplatesAPI(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// GET /api/templates - listar todos
		d.templateHandler.ListTemplates(w, r)
	case http.MethodPost:
		// POST /api/templates - criar novo
		d.templateHandler.CreateTemplate(w, r)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

// handleTemplatesAPIWithID roteia requisições da API de templates (com ID)
func (d *Dashboard) handleTemplatesAPIWithID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// GET /api/templates/:id ou /api/templates/:id/duplicate
		if strings.HasSuffix(r.URL.Path, "/duplicate") {
			// Não permitir duplicate via GET
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		} else {
			d.templateHandler.GetTemplate(w, r)
		}
	case http.MethodPost:
		// POST /api/templates/:id/duplicate
		if strings.HasSuffix(r.URL.Path, "/duplicate") {
			d.templateHandler.DuplicateTemplate(w, r)
		} else {
			http.Error(w, "Rota não encontrada", http.StatusNotFound)
		}
	case http.MethodPut:
		// PUT /api/templates/:id
		d.templateHandler.UpdateTemplate(w, r)
	case http.MethodDelete:
		// DELETE /api/templates/:id
		d.templateHandler.DeleteTemplate(w, r)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}
