package message

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// EmailStatus representa o status de envio de um email
type EmailStatus int

const (
	StatusPending          EmailStatus = 0   // Pendente
	StatusSent             EmailStatus = 2   // Enviado com sucesso
	StatusError            EmailStatus = 3   // Erro temporário (retentar)
	StatusPermanentFailure EmailStatus = 4   // Falha permanente
	StatusInvalidEmail     EmailStatus = 125 // E-mail inválido
)

// Email representa uma mensagem de email
type Email struct {
	ID               int64
	CliCodigo        sql.NullInt64
	Remetente        string
	Destinatario     string
	Assunto          string
	Corpo            string
	TipoCorpo        string // "text/plain" ou "text/html"
	StatusEnvio      EmailStatus
	DataCadastro     time.Time
	DataAgendamento  sql.NullTime
	DataEnvio        sql.NullTime
	QTDTentativas    int
	DetalhesErro     sql.NullString
	IDProvider       sql.NullString
	MetodoEnvio      sql.NullInt64
	Prioridade       int
	AnexoReferencia  sql.NullString
	AnexoNome        sql.NullString
	AnexoTipo        sql.NullString
	IPOrigem         sql.NullString
}

// Priority constants
const (
	PriorityHigh   = 1
	PriorityNormal = 2
	PriorityLow    = 3
)

// ProviderStringToCode converte nome do provider para código numérico
func ProviderStringToCode(provider string) int {
	switch strings.ToLower(provider) {
	case "mock":
		return 0
	case "smtp":
		return 1024
	case "sendgrid":
		return 2048
	case "zenvia":
		return 4096
	case "pontaltech":
		return 8192
	default:
		return 0
	}
}

// ProviderCodeToString converte código numérico para nome do provider
func ProviderCodeToString(code int) string {
	switch code {
	case 0:
		return "mock"
	case 1024:
		return "smtp"
	case 2048:
		return "sendgrid"
	case 4096:
		return "zenvia"
	case 8192:
		return "pontaltech"
	default:
		return "unknown"
	}
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

// TruncateSubject trunca o assunto se necessário
func TruncateSubject(subject string) (string, bool) {
	const maxLength = 500
	if len(subject) <= maxLength {
		return subject, false
	}
	return subject[:maxLength], true
}

// TruncateBody trunca o corpo se necessário (limite bem alto para CLOBs)
func TruncateBody(body string) (string, bool) {
	const maxLength = 1000000 // 1MB de texto
	if len(body) <= maxLength {
		return body, false
	}
	return body[:maxLength], true
}

// GetStatusDescription retorna descrição textual do status
func GetStatusDescription(status EmailStatus) string {
	switch status {
	case StatusPending:
		return "Pendente"
	case StatusSent:
		return "Enviado com sucesso"
	case StatusInvalidEmail:
		return "E-mail inválido"
	case StatusError:
		return "Erro no envio"
	case StatusPermanentFailure:
		return "Falha permanente"
	default:
		return "Desconhecido"
	}
}
