package email

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const sendGridAPIURL = "https://api.sendgrid.com/v3/mail/send"

// SendGridProvider implementa Provider para SendGrid API v3
type SendGridProvider struct {
	apiKey     string
	logger     *zap.Logger
	httpClient *http.Client
}

// SendGridRequest representa a requisição para a API SendGrid (conforme código WinDev)
type SendGridRequest struct {
	Personalizations []SendGridPersonalization `json:"personalizations"`
	From             SendGridEmail             `json:"from"`
	Subject          string                    `json:"subject"`
	Content          []SendGridContent         `json:"content"`
	Attachments      []SendGridAttachment      `json:"attachments,omitempty"`
}

type SendGridPersonalization struct {
	To []SendGridEmail `json:"to"`
}

type SendGridEmail struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type SendGridContent struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// SendGridAttachment representa um anexo (conforme código WinDev)
type SendGridAttachment struct {
	Content     string `json:"content"`     // Base64 encoded
	Filename    string `json:"filename"`    // Nome do arquivo
	Type        string `json:"type"`        // MIME type
	Disposition string `json:"disposition"` // "attachment"
}

// SendGridErrorResponse representa uma resposta de erro da API SendGrid
type SendGridErrorResponse struct {
	Errors []struct {
		Message string `json:"message"`
		Field   string `json:"field"`
		Help    string `json:"help"`
	} `json:"errors"`
}

// NewSendGridProvider cria um novo provider SendGrid
func NewSendGridProvider(apiKey string, logger *zap.Logger) *SendGridProvider {
	return &SendGridProvider{
		apiKey: apiKey,
		logger: logger,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Send envia um email via SendGrid API
func (sg *SendGridProvider) Send(ctx context.Context, email EmailData) (SendResult, error) {
	sg.logger.Info("📧 Enviando email via SendGrid",
		zap.String("to", email.To),
		zap.String("from", email.From),
		zap.String("subject", email.Subject),
		zap.String("content_type", email.ContentType),
		zap.Bool("has_attachment", email.Attachment != nil))

	// IMPORTANTE: O email "from" deve estar verificado no SendGrid
	// Acesse: https://app.sendgrid.com/settings/sender_auth/senders
	// E verifique se o domínio ou email está autenticado

	// Preparar requisição (conforme estrutura WinDev)
	req := SendGridRequest{
		Personalizations: []SendGridPersonalization{
			{
				To: []SendGridEmail{
					{Email: email.To},
				},
			},
		},
		From: SendGridEmail{
			Email: email.From,
		},
		Subject: email.Subject,
		Content: []SendGridContent{
			{
				Type:  email.ContentType,
				Value: email.Body,
			},
		},
	}

	// Adicionar anexo se houver (conforme código WinDev)
	if email.Attachment != nil && email.Attachment.Data != nil {
		// Ler dados do anexo
		attachmentData, err := io.ReadAll(email.Attachment.Data)
		if err != nil {
			sg.logger.Error("Erro ao ler dados do anexo", zap.Error(err))
			return SendResult{
				Success: false,
				Error:   fmt.Errorf("erro ao ler anexo: %w", err),
			}, err
		}

		if len(attachmentData) > 0 {
			// Codificar em base64
			base64Data := base64.StdEncoding.EncodeToString(attachmentData)

			req.Attachments = []SendGridAttachment{
				{
					Content:     base64Data,
					Filename:    email.Attachment.Filename,
					Type:        email.Attachment.ContentType,
					Disposition: "attachment",
				},
			}

			sg.logger.Debug("Anexo adicionado ao email",
				zap.String("filename", email.Attachment.Filename),
				zap.String("content_type", email.Attachment.ContentType),
				zap.Int("size_bytes", len(attachmentData)))
		}
	}

	// Serializar para JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		sg.logger.Error("Erro ao serializar requisição SendGrid",
			zap.Error(err))
		return SendResult{
			Success: false,
			Error:   fmt.Errorf("erro ao serializar requisição: %w", err),
		}, err
	}

	// SEMPRE logar JSON enviado para debug (mudado para Info)
	sg.logger.Info("📤 JSON enviado para SendGrid",
		zap.String("json", string(jsonData)))

	// Criar requisição HTTP
	httpReq, err := http.NewRequestWithContext(ctx, "POST", sendGridAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		sg.logger.Error("Erro ao criar requisição HTTP",
			zap.Error(err))
		return SendResult{
			Success: false,
			Error:   fmt.Errorf("erro ao criar requisição HTTP: %w", err),
		}, err
	}

	// Configurar headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+sg.apiKey)

	// Enviar requisição
	resp, err := sg.httpClient.Do(httpReq)
	if err != nil {
		sg.logger.Error("Erro ao enviar requisição para SendGrid",
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
		sg.logger.Error("Erro ao ler resposta da SendGrid",
			zap.Error(err),
			zap.Int("status_code", resp.StatusCode))
		return SendResult{
			Success: false,
			Error:   fmt.Errorf("erro ao ler resposta: %w", err),
		}, err
	}

	// SEMPRE logar resposta para debug
	sg.logger.Info("📩 Resposta da API SendGrid",
		zap.Int("status_code", resp.StatusCode),
		zap.String("body", string(body)),
		zap.String("x-message-id", resp.Header.Get("X-Message-Id")))

	// Verificar status code (SendGrid retorna 200 ou 202 para sucesso, conforme WinDev)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		// Tentar parsear erro
		var errorResp SendGridErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil && len(errorResp.Errors) > 0 {
			errorMsg := errorResp.Errors[0].Message
			sg.logger.Error("Erro na API SendGrid",
				zap.Int("status_code", resp.StatusCode),
				zap.String("error_message", errorMsg))

			return SendResult{
				Success: false,
				Error:   fmt.Errorf("SendGrid API error: %s", errorMsg),
			}, fmt.Errorf("SendGrid API error: %s", errorMsg)
		}

		// Erro genérico
		errorMsg := fmt.Sprintf("SendGrid API returned status %d: %s", resp.StatusCode, string(body))
		sg.logger.Error("Erro na resposta da SendGrid",
			zap.Int("status_code", resp.StatusCode),
			zap.String("body", string(body)))

		return SendResult{
			Success: false,
			Error:   fmt.Errorf(errorMsg),
		}, fmt.Errorf(errorMsg)
	}

	// SendGrid retorna X-Message-Id no header
	messageID := resp.Header.Get("X-Message-Id")
	if messageID == "" {
		messageID = fmt.Sprintf("sendgrid-%d-%d", email.ID, time.Now().Unix())
	}

	sg.logger.Info("✅ Email aceito pelo SendGrid (status 202)",
		zap.String("message_id", messageID),
		zap.String("to", email.To),
		zap.String("from", email.From),
		zap.Bool("has_attachment", len(req.Attachments) > 0))

	// Aviso: Se o email não está chegando, verifique:
	// 1. Email/domínio "from" está verificado no SendGrid
	// 2. Conta não está em Sandbox Mode
	// 3. Verifique Activity Feed: https://app.sendgrid.com/email_activity
	sg.logger.Info("ℹ️  Para rastrear entrega, acesse SendGrid Activity Feed",
		zap.String("url", "https://app.sendgrid.com/email_activity"),
		zap.String("message_id", messageID))

	return SendResult{
		Success:    true,
		ProviderID: messageID,
		Error:      nil,
	}, nil
}

// GetName retorna o nome do provider
func (sg *SendGridProvider) GetName() string {
	return "SendGrid"
}

// ValidateEmail valida o formato do email
func (sg *SendGridProvider) ValidateEmail(email string) error {
	return ValidateEmail(email)
}
