package cliente

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

// Cliente representa um cliente do sistema
type Cliente struct {
	CliCodigo  int
	CliNome    string
	CliCpfCnpj string
	Email      string // De CLIENTESEXTENSAO.CLIEXTEMAIL2
}

// Repository gerencia operações de clientes
type Repository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewRepository cria um novo repository de clientes
func NewRepository(db *sql.DB, logger *zap.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// FindByCodigo busca cliente por código
func (r *Repository) FindByCodigo(ctx context.Context, codigo int) (*Cliente, error) {
	query := `
		SELECT
			c.CLICODIGO,
			c.CLINOME,
			c.CLICPFCNPJ,
			NVL(ce.CLIEXTEMAIL2, '') as EMAIL
		FROM CLIENTES c
		LEFT JOIN CLIENTESEXTENSAO ce ON c.CLICODIGO = ce.CLICODIGO
		WHERE c.CLICODIGO = :1`

	var cliente Cliente
	err := r.db.QueryRowContext(ctx, query, codigo).Scan(
		&cliente.CliCodigo,
		&cliente.CliNome,
		&cliente.CliCpfCnpj,
		&cliente.Email,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("cliente não encontrado: código %d", codigo)
	}
	if err != nil {
		r.logger.Error("Erro ao buscar cliente por código",
			zap.Int("codigo", codigo),
			zap.Error(err))
		return nil, fmt.Errorf("erro ao buscar cliente: %w", err)
	}

	// Limpar e normalizar email
	cliente.Email = strings.TrimSpace(cliente.Email)

	r.logger.Debug("Cliente encontrado por código",
		zap.Int("codigo", codigo),
		zap.String("nome", cliente.CliNome),
		zap.Bool("tem_email", cliente.Email != ""))

	return &cliente, nil
}

// FindByCpfCnpj busca cliente por CPF/CNPJ
func (r *Repository) FindByCpfCnpj(ctx context.Context, cpfCnpj string) (*Cliente, error) {
	// Limpar CPF/CNPJ (remover caracteres não numéricos)
	cpfCnpjLimpo := LimparCpfCnpj(cpfCnpj)

	query := `
		SELECT
			c.CLICODIGO,
			c.CLINOME,
			c.CLICPFCNPJ,
			NVL(ce.CLIEXTEMAIL2, '') as EMAIL
		FROM CLIENTES c
		LEFT JOIN CLIENTESEXTENSAO ce ON c.CLICODIGO = ce.CLICODIGO
		WHERE REGEXP_REPLACE(c.CLICPFCNPJ, '[^0-9]', '') = :1`

	var cliente Cliente
	err := r.db.QueryRowContext(ctx, query, cpfCnpjLimpo).Scan(
		&cliente.CliCodigo,
		&cliente.CliNome,
		&cliente.CliCpfCnpj,
		&cliente.Email,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("cliente não encontrado: CPF/CNPJ %s", cpfCnpj)
	}
	if err != nil {
		r.logger.Error("Erro ao buscar cliente por CPF/CNPJ",
			zap.String("cpfCnpj", cpfCnpj),
			zap.Error(err))
		return nil, fmt.Errorf("erro ao buscar cliente: %w", err)
	}

	// Limpar e normalizar email
	cliente.Email = strings.TrimSpace(cliente.Email)

	r.logger.Debug("Cliente encontrado por CPF/CNPJ",
		zap.String("cpfCnpj", cpfCnpj),
		zap.Int("codigo", cliente.CliCodigo),
		zap.String("nome", cliente.CliNome),
		zap.Bool("tem_email", cliente.Email != ""))

	return &cliente, nil
}

// LimparCpfCnpj remove caracteres não numéricos de CPF/CNPJ
func LimparCpfCnpj(cpfCnpj string) string {
	reg := regexp.MustCompile(`[^0-9]`)
	return reg.ReplaceAllString(cpfCnpj, "")
}

// ValidarEmail valida se o email do cliente é válido
func ValidarEmail(email string) error {
	email = strings.TrimSpace(email)

	if email == "" {
		return fmt.Errorf("email não pode ser vazio")
	}

	// Regex simples para validação de email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	if !emailRegex.MatchString(email) {
		return fmt.Errorf("formato de email inválido")
	}

	return nil
}
