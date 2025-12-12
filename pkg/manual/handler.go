package manual

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/cliente"
	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/message"
	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/template"
	"go.uber.org/zap"
)

// Handler gerencia as requisições HTTP para disparo manual de e-mail
type Handler struct {
	clienteRepo    *cliente.Repository
	emailRepo      *message.Repository
	templateRepo   *template.Repository
	macroProcessor *template.MacroProcessor
	logger         *zap.Logger
	providerName   string // Nome do provider configurado (mock, smtp, sendgrid, zenvia, pontaltech)
}

// NewHandler cria uma nova instância do handler
func NewHandler(clienteRepo *cliente.Repository, emailRepo *message.Repository, templateRepo *template.Repository, macroProcessor *template.MacroProcessor, providerName string) *Handler {
	logger, _ := zap.NewProduction()
	return &Handler{
		clienteRepo:    clienteRepo,
		emailRepo:      emailRepo,
		templateRepo:   templateRepo,
		macroProcessor: macroProcessor,
		logger:         logger,
		providerName:   providerName,
	}
}

// ValidarClienteRequest é a requisição para validar um cliente
type ValidarClienteRequest struct {
	CliCodigo  string `json:"cliCodigo"`
	CliCpfCnpj string `json:"cliCpfCnpj"`
}

// ValidarClienteResponse é a resposta da validação do cliente
type ValidarClienteResponse struct {
	Success      bool   `json:"success"`
	Error        string `json:"error,omitempty"`
	CliCodigo    int    `json:"cliCodigo,omitempty"`
	CliCpfCnpj   string `json:"cliCpfCnpj,omitempty"`
	CliNome      string `json:"cliNome,omitempty"`
	Email        string `json:"email,omitempty"`
	EmailValido  bool   `json:"emailValido,omitempty"`
}

// DispararEmailRequest é a requisição para disparar um e-mail
type DispararEmailRequest struct {
	CliCodigo      int    `json:"cliCodigo"`
	Email          string `json:"email"`
	Assunto        string `json:"assunto"`
	Mensagem       string `json:"mensagem"`
	IsHTML         bool   `json:"isHtml"`
	TemplateID     int64  `json:"templateId"`     // ID do template (opcional)
	AttachmentData string `json:"attachmentData"` // Base64 encoded (SendGrid, Pontaltech)
	AttachmentName string `json:"attachmentName"`
	AttachmentType string `json:"attachmentType"`
	AttachmentURL  string `json:"attachmentUrl"`  // URL pública (Zenvia)
}

// DispararEmailResponse é a resposta do disparo de e-mail
type DispararEmailResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	EmailID int64  `json:"emailId,omitempty"`
	Message string `json:"message,omitempty"`
}

// StatusEmailResponse é a resposta do status de um e-mail
type StatusEmailResponse struct {
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
	Status     int    `json:"status,omitempty"`
	StatusDesc string `json:"statusDesc,omitempty"`
	DataEnvio  string `json:"dataEnvio,omitempty"`
	Tentativas int    `json:"tentativas,omitempty"`
	ErroMsg    string `json:"erroMsg,omitempty"`
	IDProvedor string `json:"idProvedor,omitempty"`
}

// ServeHTTP serve a página HTML de disparo manual
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(manualSendHTML))
}

// ProviderInfoResponse é a resposta com informações do provider
type ProviderInfoResponse struct {
	ProviderName string `json:"providerName"`
	Status       string `json:"status"`
}

// GetProviderInfo retorna informações do provider configurado
func (h *Handler) GetProviderInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	respondJSON(w, http.StatusOK, ProviderInfoResponse{
		ProviderName: h.providerName,
		Status:       "online",
	})
}

// ValidarCliente valida um cliente pelo código ou CPF/CNPJ
func (h *Handler) ValidarCliente(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req ValidarClienteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Erro ao decodificar requisição", zap.Error(err))
		respondJSON(w, http.StatusBadRequest, ValidarClienteResponse{
			Success: false,
			Error:   "Requisição inválida",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var cli *cliente.Cliente
	var err error

	// Limpar e normalizar os campos de entrada
	codigoStr := strings.TrimSpace(req.CliCodigo)
	cpfCnpjStr := strings.TrimSpace(req.CliCpfCnpj)

	h.logger.Info("Requisição de validação de cliente recebida",
		zap.String("cliCodigo_original", req.CliCodigo),
		zap.String("cliCodigo_trimmed", codigoStr),
		zap.String("cliCpfCnpj_original", req.CliCpfCnpj),
		zap.String("cliCpfCnpj_trimmed", cpfCnpjStr))

	// Tenta buscar por código primeiro, se fornecido e não vazio
	if codigoStr != "" {
		codigo, parseErr := strconv.Atoi(codigoStr)
		if parseErr != nil || codigo <= 0 {
			h.logger.Warn("Código do cliente inválido",
				zap.String("codigo_str", codigoStr),
				zap.Error(parseErr))
			respondJSON(w, http.StatusBadRequest, ValidarClienteResponse{
				Success: false,
				Error:   "Código do cliente inválido",
			})
			return
		}
		h.logger.Info("Buscando cliente por código", zap.Int("codigo", codigo))
		cli, err = h.clienteRepo.FindByCodigo(ctx, codigo)
	} else if cpfCnpjStr != "" {
		// Busca por CPF/CNPJ
		h.logger.Info("Buscando cliente por CPF/CNPJ", zap.String("cpfCnpj", cpfCnpjStr))
		cli, err = h.clienteRepo.FindByCpfCnpj(ctx, cpfCnpjStr)
	} else {
		h.logger.Warn("Nenhum critério de busca fornecido")
		respondJSON(w, http.StatusBadRequest, ValidarClienteResponse{
			Success: false,
			Error:   "Informe o código ou CPF/CNPJ do cliente",
		})
		return
	}

	if err != nil {
		h.logger.Error("Erro ao validar cliente", zap.Error(err))
		respondJSON(w, http.StatusNotFound, ValidarClienteResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Log do cliente retornado
	h.logger.Info("Cliente retornado pelo repository",
		zap.Int("cliCodigo", cli.CliCodigo),
		zap.String("cliNome", cli.CliNome),
		zap.String("email", cli.Email),
		zap.Int("tamanho_email", len(cli.Email)))

	// Valida o e-mail
	emailValido := false
	if cli.Email != "" {
		h.logger.Info("Validando e-mail do cliente",
			zap.String("email", cli.Email),
			zap.Int("tamanho", len(cli.Email)))

		if err := cliente.ValidarEmail(cli.Email); err == nil {
			emailValido = true
			h.logger.Info("E-mail válido")
		} else {
			h.logger.Warn("E-mail inválido", zap.Error(err))
		}
	} else {
		h.logger.Warn("E-mail vazio")
	}

	respondJSON(w, http.StatusOK, ValidarClienteResponse{
		Success:     true,
		CliCodigo:   cli.CliCodigo,
		CliCpfCnpj:  cli.CliCpfCnpj,
		CliNome:     cli.CliNome,
		Email:       cli.Email,
		EmailValido: emailValido,
	})
}

// DispararEmail dispara um e-mail manual
func (h *Handler) DispararEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req DispararEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Erro ao decodificar requisição", zap.Error(err))
		respondJSON(w, http.StatusBadRequest, DispararEmailResponse{
			Success: false,
			Error:   "Requisição inválida",
		})
		return
	}

	// Validações
	if req.CliCodigo <= 0 {
		respondJSON(w, http.StatusBadRequest, DispararEmailResponse{
			Success: false,
			Error:   "Código do cliente inválido",
		})
		return
	}

	if req.Email == "" {
		respondJSON(w, http.StatusBadRequest, DispararEmailResponse{
			Success: false,
			Error:   "E-mail não informado",
		})
		return
	}

	if err := cliente.ValidarEmail(req.Email); err != nil {
		respondJSON(w, http.StatusBadRequest, DispararEmailResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Se TemplateID fornecido, processar template
	var assunto, mensagem, tipoCorpo string
	var templateID sql.NullInt64

	if req.TemplateID > 0 {
		ctx2, cancel2 := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel2()

		// Buscar template
		tmpl, err := h.templateRepo.GetByID(ctx2, req.TemplateID)
		if err != nil {
			h.logger.Error("Erro ao buscar template", zap.Error(err), zap.Int64("templateId", req.TemplateID))
			respondJSON(w, http.StatusNotFound, DispararEmailResponse{
				Success: false,
				Error:   "Template não encontrado",
			})
			return
		}

		if !tmpl.Ativo {
			respondJSON(w, http.StatusBadRequest, DispararEmailResponse{
				Success: false,
				Error:   "Template está inativo",
			})
			return
		}

		// Processar template com macros do cliente
		cliCodigoNull := sql.NullInt64{Int64: int64(req.CliCodigo), Valid: true}
		assuntoProcessado, corpoProcessado, err := h.macroProcessor.ProcessTemplate(ctx2, tmpl, cliCodigoNull)
		if err != nil {
			h.logger.Error("Erro ao processar template", zap.Error(err))
			respondJSON(w, http.StatusInternalServerError, DispararEmailResponse{
				Success: false,
				Error:   "Erro ao processar template",
			})
			return
		}

		assunto = assuntoProcessado
		mensagem = corpoProcessado
		tipoCorpo = "text/html" // Templates são sempre HTML
		templateID = sql.NullInt64{Int64: req.TemplateID, Valid: true}

		h.logger.Info("Template processado com sucesso",
			zap.Int64("templateId", req.TemplateID),
			zap.Int("cliente", req.CliCodigo))
	} else {
		// Modo manual: usar dados do request
		if req.Assunto == "" {
			respondJSON(w, http.StatusBadRequest, DispararEmailResponse{
				Success: false,
				Error:   "Assunto não informado",
			})
			return
		}

		if req.Mensagem == "" {
			respondJSON(w, http.StatusBadRequest, DispararEmailResponse{
				Success: false,
				Error:   "Mensagem não informada",
			})
			return
		}

		assunto = req.Assunto
		mensagem = req.Mensagem
		tipoCorpo = "text/plain"
		if req.IsHTML {
			tipoCorpo = "text/html"
		}
	}

	// Processar anexo se fornecido
	var anexoReferencia, anexoNome, anexoTipo sql.NullString

	// Priorizar URL de anexo (para Zenvia)
	if req.AttachmentURL != "" {
		// Validar URL
		if !strings.HasPrefix(req.AttachmentURL, "http://") && !strings.HasPrefix(req.AttachmentURL, "https://") {
			respondJSON(w, http.StatusBadRequest, DispararEmailResponse{
				Success: false,
				Error:   "URL do anexo inválida. Deve começar com http:// ou https://",
			})
			return
		}

		anexoReferencia = sql.NullString{String: req.AttachmentURL, Valid: true} // URL pública
		if req.AttachmentName != "" {
			anexoNome = sql.NullString{String: req.AttachmentName, Valid: true}
		}
		anexoTipo = sql.NullString{String: "url", Valid: true} // Marcar como URL

		h.logger.Info("Anexo via URL recebido para envio",
			zap.String("url", req.AttachmentURL),
			zap.String("filename", req.AttachmentName))
	} else if req.AttachmentData != "" && req.AttachmentName != "" {
		// Anexo em base64 (SendGrid, Pontaltech)
		anexoReferencia = sql.NullString{String: req.AttachmentData, Valid: true} // Base64 data
		anexoNome = sql.NullString{String: req.AttachmentName, Valid: true}
		anexoTipo = sql.NullString{String: req.AttachmentType, Valid: true}

		h.logger.Info("Anexo em base64 recebido para envio",
			zap.String("filename", req.AttachmentName),
			zap.String("type", req.AttachmentType),
			zap.Int("base64_length", len(req.AttachmentData)))
	}

	// Extrair IP do cliente HTTP
	clientIP := getClientIP(r)
	h.logger.Debug("IP do cliente extraído",
		zap.String("ip", clientIP),
		zap.Int("cliente", req.CliCodigo))

	// Converte o provider configurado para código numérico
	providerCode := message.ProviderStringToCode(h.providerName)

	// Cria o registro de e-mail
	email := &message.Email{
		CliCodigo:       sql.NullInt64{Int64: int64(req.CliCodigo), Valid: true},
		Remetente:       "noreply@sistema.com.br", // TODO: tornar configurável
		Destinatario:    req.Email,
		Assunto:         assunto,
		Corpo:           mensagem,
		TipoCorpo:       tipoCorpo,
		StatusEnvio:     message.StatusPending,
		DataCadastro:    time.Now(),
		DataAgendamento: sql.NullTime{Time: time.Now(), Valid: true},
		Prioridade:      2, // Normal
		MetodoEnvio:     sql.NullInt64{Int64: int64(providerCode), Valid: true},
		IPOrigem:        sql.NullString{String: clientIP, Valid: clientIP != ""},
		AnexoReferencia: anexoReferencia,
		AnexoNome:       anexoNome,
		AnexoTipo:       anexoTipo,
		TemplateID:      templateID,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Insere o e-mail no banco
	emailID, err := h.emailRepo.InsertEmail(ctx, email)
	if err != nil {
		h.logger.Error("Erro ao inserir e-mail", zap.Error(err))
		respondJSON(w, http.StatusInternalServerError, DispararEmailResponse{
			Success: false,
			Error:   "Erro ao inserir e-mail no banco de dados",
		})
		return
	}

	h.logger.Info("E-mail manual inserido com sucesso",
		zap.Int64("emailId", emailID),
		zap.Int("cliente", req.CliCodigo),
		zap.String("destinatario", req.Email),
	)

	respondJSON(w, http.StatusOK, DispararEmailResponse{
		Success: true,
		EmailID: emailID,
		Message: "E-mail inserido com sucesso e será enviado em instantes",
	})
}

// ConsultarStatus consulta o status de um e-mail
func (h *Handler) ConsultarStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	emailIDStr := r.URL.Query().Get("id")
	if emailIDStr == "" {
		respondJSON(w, http.StatusBadRequest, StatusEmailResponse{
			Success: false,
			Error:   "ID do e-mail não informado",
		})
		return
	}

	emailID, err := strconv.ParseInt(emailIDStr, 10, 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, StatusEmailResponse{
			Success: false,
			Error:   "ID do e-mail inválido",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Busca o e-mail no banco
	email, err := h.emailRepo.GetByID(ctx, emailID)
	if err != nil {
		h.logger.Error("Erro ao buscar e-mail", zap.Error(err), zap.Int64("emailId", emailID))
		respondJSON(w, http.StatusNotFound, StatusEmailResponse{
			Success: false,
			Error:   "E-mail não encontrado",
		})
		return
	}

	// Log do status lido do banco
	h.logger.Info("Status do e-mail consultado",
		zap.Int64("emailId", emailID),
		zap.Int("status", int(email.StatusEnvio)),
		zap.Int("tentativas", email.QTDTentativas),
		zap.String("provider_id", email.IDProvider.String))

	statusDesc := getStatusDescription(email.StatusEnvio)
	dataEnvio := ""
	if email.DataEnvio.Valid {
		dataEnvio = email.DataEnvio.Time.Format("02/01/2006 15:04:05")
	}

	erroMsg := ""
	if email.DetalhesErro.Valid {
		erroMsg = email.DetalhesErro.String
	}

	idProvedor := ""
	if email.IDProvider.Valid {
		idProvedor = email.IDProvider.String
	}

	respondJSON(w, http.StatusOK, StatusEmailResponse{
		Success:    true,
		Status:     int(email.StatusEnvio),
		StatusDesc: statusDesc,
		DataEnvio:  dataEnvio,
		Tentativas: email.QTDTentativas,
		ErroMsg:    erroMsg,
		IDProvedor: idProvedor,
	})
}

// getStatusDescription retorna a descrição do status
func getStatusDescription(status message.EmailStatus) string {
	switch status {
	case message.StatusPending:
		return "Pendente"
	case message.StatusSent:
		return "Enviado com sucesso"
	case message.StatusInvalidEmail:
		return "E-mail inválido"
	case message.StatusError:
		return "Erro no envio"
	case message.StatusPermanentFailure:
		return "Falha permanente"
	default:
		return "Desconhecido"
	}
}

// respondJSON envia uma resposta JSON
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// getClientIP extrai o IP do cliente HTTP, considerando proxies
func getClientIP(r *http.Request) string {
	// Verificar headers de proxy primeiro
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		// X-Forwarded-For pode conter múltiplos IPs: "client, proxy1, proxy2"
		// Pegamos o primeiro (IP real do cliente)
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	// Fallback: usar RemoteAddr
	// RemoteAddr está no formato "IP:Port", extrair apenas o IP
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}

	return ip
}

// PreviewTemplateRequest é a requisição para preview de template
type PreviewTemplateRequest struct {
	TemplateID int64 `json:"templateId"`
	CliCodigo  int   `json:"cliCodigo"`
}

// PreviewTemplateResponse é a resposta do preview de template
type PreviewTemplateResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Assunto string `json:"assunto,omitempty"`
	Corpo   string `json:"corpo,omitempty"`
}

// PreviewTemplate processa um template com os dados do cliente e retorna o preview
func (h *Handler) PreviewTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req PreviewTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Erro ao decodificar requisição", zap.Error(err))
		respondJSON(w, http.StatusBadRequest, PreviewTemplateResponse{
			Success: false,
			Error:   "Requisição inválida",
		})
		return
	}

	if req.TemplateID <= 0 {
		respondJSON(w, http.StatusBadRequest, PreviewTemplateResponse{
			Success: false,
			Error:   "ID do template inválido",
		})
		return
	}

	if req.CliCodigo <= 0 {
		respondJSON(w, http.StatusBadRequest, PreviewTemplateResponse{
			Success: false,
			Error:   "Código do cliente inválido",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Buscar template
	tmpl, err := h.templateRepo.GetByID(ctx, req.TemplateID)
	if err != nil {
		h.logger.Error("Erro ao buscar template", zap.Error(err), zap.Int64("templateId", req.TemplateID))
		respondJSON(w, http.StatusNotFound, PreviewTemplateResponse{
			Success: false,
			Error:   "Template não encontrado",
		})
		return
	}

	if !tmpl.Ativo {
		respondJSON(w, http.StatusBadRequest, PreviewTemplateResponse{
			Success: false,
			Error:   "Template está inativo",
		})
		return
	}

	// Processar template com macros do cliente
	cliCodigoNull := sql.NullInt64{Int64: int64(req.CliCodigo), Valid: true}
	assunto, corpo, err := h.macroProcessor.ProcessTemplate(ctx, tmpl, cliCodigoNull)
	if err != nil {
		h.logger.Error("Erro ao processar template", zap.Error(err))
		respondJSON(w, http.StatusInternalServerError, PreviewTemplateResponse{
			Success: false,
			Error:   "Erro ao processar template",
		})
		return
	}

	h.logger.Info("Preview de template gerado com sucesso",
		zap.Int64("templateId", req.TemplateID),
		zap.Int("cliente", req.CliCodigo))

	respondJSON(w, http.StatusOK, PreviewTemplateResponse{
		Success: true,
		Assunto: assunto,
		Corpo:   corpo,
	})
}
