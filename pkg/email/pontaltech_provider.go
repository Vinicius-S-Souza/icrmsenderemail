package email

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

const pontaltechEmailAPIURL = "https://pointer-email-api.pontaltech.com.br/send"

// PontaltechProvider implementa Provider para Pontaltech Email API
type PontaltechProvider struct {
	username    string
	password    string
	accountID   int
	apiURL      string // URL customizada da API
	callbackURL string // URL de callback para notificações
	logger      *zap.Logger
	httpClient  *http.Client
}

// PontaltechAttachment representa um anexo
type PontaltechAttachment struct {
	Filename string `json:"filename"`
	Data     string `json:"data"` // Base64 encoded
}

// PontaltechMessageVariable representa variáveis da mensagem
type PontaltechMessageVariable struct {
	Nome string `json:"nome"`
}

// PontaltechRecipient representa um destinatário
type PontaltechRecipient struct {
	Email           string                     `json:"email"`
	MessageVariable *PontaltechMessageVariable `json:"messageVariable,omitempty"`
	Attachments     []PontaltechAttachment     `json:"attachments,omitempty"`
}

// PontaltechEmailRequest representa a requisição para a API Pontaltech
// Baseado no código WinDev fornecido
type PontaltechEmailRequest struct {
	To              []PontaltechRecipient `json:"to"`
	FromGroup       string                `json:"fromGroup"`
	MailBody        string                `json:"mailBody"`
	Subject         string                `json:"subject"`
	ReplyTo         string                `json:"replyTo"`
	Sender          string                `json:"sender"`
	AccountID       int                   `json:"accountId"`
	Tracking        bool                  `json:"tracking"`
	AttachmentField bool                  `json:"attachmentField"`
	ReplaceVariable bool                  `json:"replaceVariable"`
	URLCallback     string                `json:"urlCallback,omitempty"`
}

// PontaltechMessage representa uma mensagem enviada na resposta
type PontaltechMessage struct {
	Email string `json:"email"`
	ID    int64  `json:"id"`
}

// PontaltechEmailResponse representa a resposta da API Pontaltech
type PontaltechEmailResponse struct {
	Messages        []PontaltechMessage `json:"messages"`
	InvalidMessages []string            `json:"invalidMessages"`
	CampaignID      int64               `json:"campaignId"`
}

// PontaltechErrorResponse representa uma resposta de erro da API Pontaltech
type PontaltechErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewPontaltechProvider cria um novo provider Pontaltech
func NewPontaltechProvider(username, password string, accountID int, apiURL, callbackURL string, logger *zap.Logger) *PontaltechProvider {
	// Se URL customizada não foi fornecida, usa a padrão
	if apiURL == "" {
		apiURL = pontaltechEmailAPIURL
		logger.Info("📧 Usando URL padrão da API Pontaltech",
			zap.String("url", apiURL))
	} else {
		logger.Info("📧 Usando URL customizada da API Pontaltech",
			zap.String("url", apiURL))
	}

	if callbackURL != "" {
		logger.Info("📞 URL de callback Pontaltech configurada",
			zap.String("callback_url", callbackURL))
	}

	return &PontaltechProvider{
		username:    username,
		password:    password,
		accountID:   accountID,
		apiURL:      apiURL,
		callbackURL: callbackURL,
		logger:      logger,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Send envia um email via Pontaltech API
func (p *PontaltechProvider) Send(ctx context.Context, email EmailData) (SendResult, error) {
	p.logger.Debug("Enviando email via Pontaltech",
		zap.String("to", email.To),
		zap.String("from", email.From),
		zap.Bool("has_attachment", email.Attachment != nil))

	// Preparar destinatário
	recipient := PontaltechRecipient{
		Email: email.To,
	}

	// Adicionar anexo se houver
	hasAttachment := false
	if email.Attachment != nil && email.Attachment.Data != nil {
		// Ler dados do anexo
		attachmentData, err := io.ReadAll(email.Attachment.Data)
		if err != nil {
			p.logger.Error("Erro ao ler dados do anexo", zap.Error(err))
			return SendResult{
				Success: false,
				Error:   fmt.Errorf("erro ao ler anexo: %w", err),
			}, err
		}

		if len(attachmentData) > 0 {
			// Codificar anexo em base64
			base64Data := base64.StdEncoding.EncodeToString(attachmentData)

			recipient.Attachments = []PontaltechAttachment{
				{
					Filename: email.Attachment.Filename,
					Data:     base64Data,
				},
			}

			// Adicionar variável de mensagem (exemplo do WinDev)
			recipient.MessageVariable = &PontaltechMessageVariable{
				Nome: "nometeste",
			}

			hasAttachment = true

			p.logger.Debug("Anexo adicionado ao email",
				zap.String("filename", email.Attachment.Filename),
				zap.String("content_type", email.Attachment.ContentType),
				zap.Int("size_bytes", len(attachmentData)))
		}
	}

	// Preparar requisição conforme o código WinDev
	req := PontaltechEmailRequest{
		To: []PontaltechRecipient{recipient},
		FromGroup: "Padrão", // Grupo de origem padrão
		MailBody:  email.Body,
		Subject:   email.Subject,
		ReplyTo:   "", // Pode ser configurado se necessário
		Sender:    email.From,
		Tracking:  true, // Habilitar tracking
		AttachmentField: hasAttachment,
		ReplaceVariable: hasAttachment, // Conforme lógica do WinDev
	}

	// Adicionar Account ID se configurado
	if p.accountID > 0 {
		req.AccountID = p.accountID
	}

	// Adicionar URL de callback se configurada
	if p.callbackURL != "" {
		req.URLCallback = p.callbackURL
	}

	// Serializar para JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		p.logger.Error("Erro ao serializar requisição Pontaltech",
			zap.Error(err))
		return SendResult{
			Success: false,
			Error:   fmt.Errorf("erro ao serializar requisição: %w", err),
		}, err
	}

	p.logger.Debug("JSON enviado para Pontaltech",
		zap.String("json", string(jsonData)))

	// Criar requisição HTTP usando URL customizada
	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		p.logger.Error("Erro ao criar requisição HTTP",
			zap.String("url", p.apiURL),
			zap.Error(err))
		return SendResult{
			Success: false,
			Error:   fmt.Errorf("erro ao criar requisição HTTP: %w", err),
		}, err
	}

	// Configurar headers (conforme o código WinDev)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Cache-Control", "no-cache")

	// Configurar Basic Authentication
	auth := base64.StdEncoding.EncodeToString([]byte(p.username + ":" + p.password))
	httpReq.Header.Set("Authorization", "Basic "+auth)

	// Enviar requisição
	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		// Verificar se é erro de DNS
		if strings.Contains(err.Error(), "no such host") {
			p.logger.Error("❌ Erro de DNS ao acessar API Pontaltech",
				zap.String("url", p.apiURL),
				zap.Error(err),
				zap.String("solucao", "Verifique se a URL da API está correta. Configure 'pontaltech_api_url' no dbinit.ini"))
			return SendResult{
				Success: false,
				Error:   fmt.Errorf("erro de DNS - domínio não encontrado: %w (verifique a URL da API Pontaltech)", err),
			}, err
		}

		p.logger.Error("Erro ao enviar requisição para Pontaltech",
			zap.String("url", p.apiURL),
			zap.Error(err))
		return SendResult{
			Success: false,
			Error:   fmt.Errorf("erro ao enviar requisição: %w", err),
		}, err
	}
	defer resp.Body.Close()

	// Ler resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		p.logger.Error("Erro ao ler resposta da Pontaltech",
			zap.Error(err),
			zap.Int("status_code", resp.StatusCode))
		return SendResult{
			Success: false,
			Error:   fmt.Errorf("erro ao ler resposta: %w", err),
		}, err
	}

	// SEMPRE logar resposta para debug (mudado de Debug para Info)
	p.logger.Info("📩 Resposta da API Pontaltech",
		zap.Int("status_code", resp.StatusCode),
		zap.String("body", string(body)))

	// Verificar status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		// Tentar parsear erro
		var errorResp PontaltechErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil {
			p.logger.Error("Erro na API Pontaltech",
				zap.Int("status_code", resp.StatusCode),
				zap.Int("error_code", errorResp.Code),
				zap.String("error_message", errorResp.Message))

			errorMsg := fmt.Sprintf("Pontaltech API error %d: %s", errorResp.Code, errorResp.Message)

			// Verificar se é erro de email inválido
			if resp.StatusCode == http.StatusBadRequest {
				errorMsg = "email inválido: " + errorResp.Message
			}

			return SendResult{
				Success: false,
				Error:   fmt.Errorf(errorMsg),
			}, fmt.Errorf(errorMsg)
		}

		// Erro genérico
		errorMsg := fmt.Sprintf("Pontaltech API returned status %d: %s", resp.StatusCode, string(body))
		p.logger.Error("Erro na resposta da Pontaltech",
			zap.Int("status_code", resp.StatusCode),
			zap.String("body", string(body)))

		return SendResult{
			Success: false,
			Error:   fmt.Errorf(errorMsg),
		}, fmt.Errorf(errorMsg)
	}

	// Parsear resposta de sucesso
	var pontaltechResp PontaltechEmailResponse
	if err := json.Unmarshal(body, &pontaltechResp); err != nil {
		// Se não conseguir parsear mas o status HTTP é sucesso, aceitar mesmo assim
		p.logger.Warn("⚠️ Não foi possível parsear resposta JSON, mas status HTTP indica sucesso",
			zap.Error(err),
			zap.Int("status_code", resp.StatusCode),
			zap.String("body", string(body)))

		// Considerar sucesso baseado apenas no status HTTP
		return SendResult{
			Success:    true,
			ProviderID: fmt.Sprintf("pontaltech-%d", time.Now().Unix()),
			Error:      nil,
		}, nil
	}

	// Verificar se há mensagens inválidas
	if len(pontaltechResp.InvalidMessages) > 0 {
		p.logger.Error("❌ Email rejeitado pela API Pontaltech - endereço inválido",
			zap.Strings("invalid_messages", pontaltechResp.InvalidMessages),
			zap.String("to", email.To))

		return SendResult{
			Success: false,
			Error:   fmt.Errorf("email inválido: %v", pontaltechResp.InvalidMessages),
		}, fmt.Errorf("email inválido: %v", pontaltechResp.InvalidMessages)
	}

	// Verificar se há mensagens enviadas
	if len(pontaltechResp.Messages) == 0 {
		p.logger.Error("❌ API Pontaltech não retornou nenhuma mensagem enviada",
			zap.String("to", email.To),
			zap.String("body", string(body)))

		return SendResult{
			Success: false,
			Error:   fmt.Errorf("API não retornou mensagens enviadas"),
		}, fmt.Errorf("API não retornou mensagens enviadas")
	}

	// Extrair ID da primeira mensagem
	messageID := pontaltechResp.Messages[0].ID
	providerID := fmt.Sprintf("%d", messageID)

	p.logger.Info("✅ Email enviado via Pontaltech com sucesso",
		zap.String("message_id", providerID),
		zap.Int64("campaign_id", pontaltechResp.CampaignID),
		zap.String("to", email.To),
		zap.Bool("has_attachment", hasAttachment))

	return SendResult{
		Success:    true,
		ProviderID: providerID,
		Error:      nil,
	}, nil
}

// GetName retorna o nome do provider
func (p *PontaltechProvider) GetName() string {
	return "Pontaltech"
}

// ValidateEmail valida o formato do email
func (p *PontaltechProvider) ValidateEmail(email string) error {
	return ValidateEmail(email)
}
