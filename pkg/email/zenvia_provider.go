package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const zenviaEmailAPIURL = "https://api.zenvia.com/v2/channels/email/messages"

// ZenviaProvider implementa Provider para Zenvia Email API
type ZenviaProvider struct {
	apiToken   string
	logger     *zap.Logger
	httpClient *http.Client
}

// ZenviaAttachment representa um anexo
// IMPORTANTE: Zenvia só aceita anexos via URL pública, não base64
type ZenviaAttachment struct {
	FileURL  string `json:"fileUrl"`           // URL pública do arquivo
	FileName string `json:"fileName,omitempty"` // Nome opcional do arquivo
}

// ZenviaEmailContent representa o conteúdo da mensagem (baseado no código WinDev)
type ZenviaEmailContent struct {
	Type        string             `json:"type"`        // Sempre "email"
	Subject     string             `json:"subject"`     // Assunto vai DENTRO do contents
	HTML        string             `json:"html"`        // Corpo HTML
	Attachments []ZenviaAttachment `json:"attachments,omitempty"`
}

// ZenviaEmailRequest representa a requisição para a API Zenvia
type ZenviaEmailRequest struct {
	From     string               `json:"from"`
	To       string               `json:"to"`
	Contents []ZenviaEmailContent `json:"contents"`
}

// ZenviaEmailResponse representa a resposta da API Zenvia
type ZenviaEmailResponse struct {
	ID        string `json:"id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Direction string `json:"direction"`
	Channel   string `json:"channel"`
	Timestamp string `json:"timestamp"`
}

// ZenviaErrorResponse representa uma resposta de erro da API Zenvia
type ZenviaErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details []struct {
		Code    string `json:"code"`
		Path    string `json:"path"`
		Message string `json:"message"`
	} `json:"details,omitempty"`
}

// NewZenviaProvider cria um novo provider Zenvia
func NewZenviaProvider(apiToken string, logger *zap.Logger) *ZenviaProvider {
	return &ZenviaProvider{
		apiToken: apiToken,
		logger:   logger,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Send envia um email via Zenvia API
func (z *ZenviaProvider) Send(ctx context.Context, email EmailData) (SendResult, error) {
	z.logger.Debug("Enviando email via Zenvia",
		zap.String("to", email.To),
		zap.String("from", email.From))

	// Preparar conteúdo (conforme código WinDev)
	// O tipo é SEMPRE "email", e o subject vai DENTRO do contents
	content := ZenviaEmailContent{
		Type:    "email",
		Subject: email.Subject,
		HTML:    email.Body, // Sempre usar HTML (mesmo que seja texto simples)
	}

	// IMPORTANTE: Zenvia só aceita anexos via URL pública!
	if email.Attachment != nil {
		if email.Attachment.URL != "" {
			// Anexo via URL (correto para Zenvia)
			content.Attachments = []ZenviaAttachment{
				{
					FileURL:  email.Attachment.URL,
					FileName: email.Attachment.Filename,
				},
			}
			z.logger.Debug("Anexo via URL adicionado ao email Zenvia",
				zap.String("url", email.Attachment.URL),
				zap.String("filename", email.Attachment.Filename))
		} else if email.Attachment.Data != nil {
			// Anexo em base64 - NÃO suportado pela Zenvia
			z.logger.Warn("⚠️  AVISO: Zenvia não suporta anexos em base64",
				zap.String("filename", email.Attachment.Filename),
				zap.String("info", "Zenvia só aceita anexos via URL pública (fileUrl). O anexo será ignorado."))
			// Anexo será ignorado - continua sem anexo
		}
	}

	// Preparar requisição
	req := ZenviaEmailRequest{
		From:     email.From,
		To:       email.To,
		Contents: []ZenviaEmailContent{content},
	}

	// Serializar para JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		z.logger.Error("Erro ao serializar requisição Zenvia",
			zap.Error(err))
		return SendResult{
			Success: false,
			Error:   fmt.Errorf("erro ao serializar requisição: %w", err),
		}, err
	}

	// Logar JSON enviado para debug
	z.logger.Info("📤 JSON enviado para Zenvia",
		zap.String("json", string(jsonData)))

	// Criar requisição HTTP
	httpReq, err := http.NewRequestWithContext(ctx, "POST", zenviaEmailAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		z.logger.Error("Erro ao criar requisição HTTP",
			zap.Error(err))
		return SendResult{
			Success: false,
			Error:   fmt.Errorf("erro ao criar requisição HTTP: %w", err),
		}, err
	}

	// Configurar headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-TOKEN", z.apiToken)

	// Enviar requisição
	resp, err := z.httpClient.Do(httpReq)
	if err != nil {
		z.logger.Error("Erro ao enviar requisição para Zenvia",
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
		z.logger.Error("Erro ao ler resposta da Zenvia",
			zap.Error(err),
			zap.Int("status_code", resp.StatusCode))
		return SendResult{
			Success: false,
			Error:   fmt.Errorf("erro ao ler resposta: %w", err),
		}, err
	}

	// Verificar status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		// Tentar parsear erro
		var errorResp ZenviaErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil {
			z.logger.Error("Erro na API Zenvia",
				zap.Int("status_code", resp.StatusCode),
				zap.String("error_code", errorResp.Code),
				zap.String("error_message", errorResp.Message))

			errorMsg := fmt.Sprintf("Zenvia API error: %s - %s", errorResp.Code, errorResp.Message)

			// Verificar se é erro de email inválido
			if errorResp.Code == "VALIDATION_ERROR" {
				for _, detail := range errorResp.Details {
					if detail.Path == "to" {
						errorMsg = "email inválido: " + detail.Message
						break
					}
				}
			}

			return SendResult{
				Success: false,
				Error:   fmt.Errorf(errorMsg),
			}, fmt.Errorf(errorMsg)
		}

		// Erro genérico
		errorMsg := fmt.Sprintf("Zenvia API returned status %d: %s", resp.StatusCode, string(body))
		z.logger.Error("Erro na resposta da Zenvia",
			zap.Int("status_code", resp.StatusCode),
			zap.String("body", string(body)))

		return SendResult{
			Success: false,
			Error:   fmt.Errorf(errorMsg),
		}, fmt.Errorf(errorMsg)
	}

	// Parsear resposta de sucesso
	var zenviaResp ZenviaEmailResponse
	if err := json.Unmarshal(body, &zenviaResp); err != nil {
		z.logger.Error("Erro ao parsear resposta de sucesso da Zenvia",
			zap.Error(err),
			zap.String("body", string(body)))
		return SendResult{
			Success: false,
			Error:   fmt.Errorf("erro ao parsear resposta: %w", err),
		}, err
	}

	z.logger.Info("Email enviado via Zenvia com sucesso",
		zap.String("message_id", zenviaResp.ID),
		zap.String("to", email.To),
		zap.String("timestamp", zenviaResp.Timestamp))

	return SendResult{
		Success:    true,
		ProviderID: zenviaResp.ID,
		Error:      nil,
	}, nil
}

// GetName retorna o nome do provider
func (z *ZenviaProvider) GetName() string {
	return "Zenvia"
}

// ValidateEmail valida o formato do email
func (z *ZenviaProvider) ValidateEmail(email string) error {
	return ValidateEmail(email)
}
