package email

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

// SendResult representa o resultado de um envio de email
type SendResult struct {
	Success    bool
	ProviderID string // ID da mensagem no provider
	Error      error
}

// Provider interface que todos os provedores de email devem implementar
type Provider interface {
	Send(ctx context.Context, email EmailData) (SendResult, error)
	GetName() string
	ValidateEmail(email string) error
}

// EmailData contém os dados necessários para enviar um email
type EmailData struct {
	ID          int64
	From        string
	To          string
	Subject     string
	Body        string
	ContentType string // "text/plain" ou "text/html"
	Attachment  *Attachment
}

// Attachment representa um anexo de email
type Attachment struct {
	Filename    string
	ContentType string
	Data        io.Reader // Para anexos em base64 (SendGrid, Pontaltech)
	Size        int64
	URL         string // Para anexos via URL pública (Zenvia)
}

// Sender gerencia o envio de emails através de um provider
type Sender struct {
	provider Provider
	logger   *zap.Logger
}

// NewSender cria um novo sender com o provider especificado
func NewSender(provider Provider, logger *zap.Logger) *Sender {
	return &Sender{
		provider: provider,
		logger:   logger,
	}
}

// Send envia um email através do provider configurado
func (s *Sender) Send(ctx context.Context, email EmailData) SendResult {
	s.logger.Debug("Enviando email",
		zap.Int64("id", email.ID),
		zap.String("to", maskEmail(email.To)),
		zap.String("subject", email.Subject),
		zap.String("provider", s.provider.GetName()))

	// Validar email de destino
	if err := s.provider.ValidateEmail(email.To); err != nil {
		s.logger.Error("Email de destino inválido",
			zap.Int64("id", email.ID),
			zap.String("to", email.To),
			zap.Error(err))
		return SendResult{
			Success: false,
			Error:   fmt.Errorf("email inválido: %w", err),
		}
	}

	// Enviar através do provider
	result, err := s.provider.Send(ctx, email)
	if err != nil {
		s.logger.Error("Erro ao enviar email",
			zap.Int64("id", email.ID),
			zap.String("provider", s.provider.GetName()),
			zap.Error(err))
		return SendResult{
			Success: false,
			Error:   err,
		}
	}

	if result.Success {
		s.logger.Info("Email enviado com sucesso",
			zap.Int64("id", email.ID),
			zap.String("provider_id", result.ProviderID),
			zap.String("provider", s.provider.GetName()))
	}

	return result
}

// GetProvider retorna o provider atual
func (s *Sender) GetProvider() Provider {
	return s.provider
}

// maskEmail mascara parte do email para logs
func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "***"
	}

	localPart := parts[0]
	domain := parts[1]

	if len(localPart) <= 2 {
		return "**@" + domain
	}

	return localPart[:2] + "***@" + domain
}

// ValidateEmail valida formato de e-mail
func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)

	if email == "" {
		return fmt.Errorf("e-mail não pode ser vazio")
	}

	// Regex simples para validação de e-mail
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	if !emailRegex.MatchString(email) {
		return fmt.Errorf("formato de e-mail inválido: %s", email)
	}

	// Validar comprimento máximo
	if len(email) > 254 {
		return fmt.Errorf("e-mail muito longo: %d caracteres (máximo 254)", len(email))
	}

	return nil
}
