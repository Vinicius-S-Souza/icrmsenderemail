package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/smtp"
	"strings"
	"time"

	"go.uber.org/zap"
)

// SMTPProvider implementa Provider para SMTP genérico
type SMTPProvider struct {
	host     string
	port     int
	username string
	password string
	useTLS   bool
	logger   *zap.Logger
}

// NewSMTPProvider cria um novo provider SMTP
func NewSMTPProvider(host string, port int, username, password string, useTLS bool, logger *zap.Logger) *SMTPProvider {
	return &SMTPProvider{
		host:     host,
		port:     port,
		username: username,
		password: password,
		useTLS:   useTLS,
		logger:   logger,
	}
}

// Send envia um email via SMTP
func (s *SMTPProvider) Send(ctx context.Context, email EmailData) (SendResult, error) {
	s.logger.Debug("Enviando email via SMTP",
		zap.String("host", s.host),
		zap.Int("port", s.port),
		zap.String("to", email.To))

	// Montar mensagem no formato RFC 822
	msg := s.buildMessage(email)

	// Conectar ao servidor SMTP
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	var err error
	if s.useTLS {
		err = s.sendWithTLS(addr, email.From, []string{email.To}, []byte(msg))
	} else {
		err = smtp.SendMail(addr, s.getAuth(), email.From, []string{email.To}, []byte(msg))
	}

	if err != nil {
		s.logger.Error("Erro ao enviar email via SMTP",
			zap.Error(err),
			zap.String("to", email.To))
		return SendResult{
			Success: false,
			Error:   fmt.Errorf("erro SMTP: %w", err),
		}, err
	}

	// SMTP não retorna ID de mensagem, gerar um baseado no timestamp
	providerID := fmt.Sprintf("smtp-%d-%d", email.ID, time.Now().Unix())

	s.logger.Info("Email enviado via SMTP com sucesso",
		zap.String("to", email.To),
		zap.String("provider_id", providerID))

	return SendResult{
		Success:    true,
		ProviderID: providerID,
		Error:      nil,
	}, nil
}

// sendWithTLS envia email usando TLS
func (s *SMTPProvider) sendWithTLS(addr, from string, to []string, msg []byte) error {
	// Configurar TLS
	tlsConfig := &tls.Config{
		ServerName:         s.host,
		InsecureSkipVerify: false,
	}

	// Conectar com TLS
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("erro ao conectar com TLS: %w", err)
	}
	defer conn.Close()

	// Criar cliente SMTP
	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("erro ao criar cliente SMTP: %w", err)
	}
	defer client.Quit()

	// Autenticar se credenciais fornecidas
	if s.username != "" && s.password != "" {
		auth := s.getAuth()
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("erro de autenticação: %w", err)
		}
	}

	// Enviar email
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("erro no MAIL FROM: %w", err)
	}

	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("erro no RCPT TO: %w", err)
		}
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("erro no DATA: %w", err)
	}

	_, err = writer.Write(msg)
	if err != nil {
		return fmt.Errorf("erro ao escrever mensagem: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("erro ao fechar writer: %w", err)
	}

	return nil
}

// getAuth retorna o mecanismo de autenticação
func (s *SMTPProvider) getAuth() smtp.Auth {
	if s.username == "" || s.password == "" {
		return nil
	}
	return smtp.PlainAuth("", s.username, s.password, s.host)
}

// buildMessage constrói a mensagem no formato RFC 822
func (s *SMTPProvider) buildMessage(email EmailData) string {
	var msg strings.Builder

	// Headers obrigatórios
	msg.WriteString(fmt.Sprintf("From: %s\r\n", email.From))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", email.To))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))
	msg.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
	msg.WriteString("MIME-Version: 1.0\r\n")

	// Se houver anexo, usar multipart
	if email.Attachment != nil {
		boundary := fmt.Sprintf("boundary-%d", time.Now().Unix())
		msg.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", boundary))
		msg.WriteString("\r\n")

		// Parte do corpo
		msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		msg.WriteString(fmt.Sprintf("Content-Type: %s; charset=\"UTF-8\"\r\n", email.ContentType))
		msg.WriteString("\r\n")
		msg.WriteString(email.Body)
		msg.WriteString("\r\n")

		// Parte do anexo
		msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		msg.WriteString(fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n", email.Attachment.ContentType, email.Attachment.Filename))
		msg.WriteString("Content-Transfer-Encoding: base64\r\n")
		msg.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", email.Attachment.Filename))
		msg.WriteString("\r\n")

		// Ler e codificar anexo (simplificado - em produção use base64)
		if email.Attachment.Data != nil {
			data, _ := io.ReadAll(email.Attachment.Data)
			msg.WriteString(string(data))
		}
		msg.WriteString("\r\n")

		msg.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	} else {
		// Sem anexo - mensagem simples
		msg.WriteString(fmt.Sprintf("Content-Type: %s; charset=\"UTF-8\"\r\n", email.ContentType))
		msg.WriteString("\r\n")
		msg.WriteString(email.Body)
	}

	return msg.String()
}

// GetName retorna o nome do provider
func (s *SMTPProvider) GetName() string {
	return "SMTP"
}

// ValidateEmail valida o formato do email
func (s *SMTPProvider) ValidateEmail(email string) error {
	return ValidateEmail(email)
}
