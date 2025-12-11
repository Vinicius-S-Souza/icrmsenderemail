package email

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// MockProvider implementa Provider para testes (simula envio sem realmente enviar)
type MockProvider struct {
	logger *zap.Logger
	delay  time.Duration
}

// NewMockProvider cria um novo provider mock
func NewMockProvider(logger *zap.Logger) *MockProvider {
	return &MockProvider{
		logger: logger,
		delay:  100 * time.Millisecond, // Simula delay de envio
	}
}

// Send simula o envio de um email
func (m *MockProvider) Send(ctx context.Context, email EmailData) (SendResult, error) {
	m.logger.Info("📧 [MOCK] Simulando envio de email",
		zap.String("to", email.To),
		zap.String("from", email.From),
		zap.String("subject", email.Subject),
		zap.Int("body_length", len(email.Body)),
		zap.String("content_type", email.ContentType))

	// Simular delay de envio
	select {
	case <-ctx.Done():
		return SendResult{
			Success: false,
			Error:   ctx.Err(),
		}, ctx.Err()
	case <-time.After(m.delay):
		// Continuar
	}

	// Gerar ID fictício
	mockID := fmt.Sprintf("mock-%d-%d", email.ID, time.Now().Unix())

	m.logger.Info("✅ [MOCK] Email enviado com sucesso",
		zap.String("mock_id", mockID),
		zap.String("to", email.To))

	return SendResult{
		Success:    true,
		ProviderID: mockID,
		Error:      nil,
	}, nil
}

// GetName retorna o nome do provider
func (m *MockProvider) GetName() string {
	return "Mock"
}

// ValidateEmail valida o formato do email
func (m *MockProvider) ValidateEmail(email string) error {
	return ValidateEmail(email)
}
