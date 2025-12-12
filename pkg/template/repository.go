package template

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

// Repository gerencia operações de banco de dados para templates
type Repository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewRepository cria um novo repository
func NewRepository(db *sql.DB, logger *zap.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Create insere um novo template no banco
func (r *Repository) Create(ctx context.Context, template *Template) (int64, error) {
	query := `
		INSERT INTO TEMPLATEEMAIL (
			ID, NOME, DESCRICAO, HEADER_HTML, BODY_HTML,
			FOOTER_HTML, ASSUNTO_PADRAO, ATIVO,
			DATA_CRIACAO, DATA_ATUALIZACAO, CRIADO_POR
		) VALUES (
			SEQ_TEMPLATEEMAIL.NEXTVAL, :1, :2, :3, :4,
			:5, :6, :7,
			SYSDATE, SYSDATE, :8
		) RETURNING ID INTO :9`

	var id int64
	_, err := r.db.ExecContext(ctx, query,
		template.Nome,
		template.Descricao,
		template.HeaderHTML,
		template.BodyHTML,
		template.FooterHTML,
		template.AssuntoPadrao,
		boolToInt(template.Ativo),
		template.CriadoPor,
		sql.Out{Dest: &id},
	)

	if err != nil {
		// Verificar se é erro de nome duplicado
		if strings.Contains(err.Error(), "UK_TEMPLATEEMAIL_NOME") {
			return 0, ErrNomeDuplicado
		}
		return 0, fmt.Errorf("erro ao criar template: %w", err)
	}

	r.logger.Info("Template criado com sucesso",
		zap.Int64("id", id),
		zap.String("nome", template.Nome))

	return id, nil
}

// Update atualiza um template existente
func (r *Repository) Update(ctx context.Context, template *Template) error {
	query := `
		UPDATE TEMPLATEEMAIL
		SET NOME = :1,
			DESCRICAO = :2,
			HEADER_HTML = :3,
			BODY_HTML = :4,
			FOOTER_HTML = :5,
			ASSUNTO_PADRAO = :6,
			ATIVO = :7,
			DATA_ATUALIZACAO = SYSDATE,
			CRIADO_POR = :8
		WHERE ID = :9`

	result, err := r.db.ExecContext(ctx, query,
		template.Nome,
		template.Descricao,
		template.HeaderHTML,
		template.BodyHTML,
		template.FooterHTML,
		template.AssuntoPadrao,
		boolToInt(template.Ativo),
		template.CriadoPor,
		template.ID,
	)

	if err != nil {
		// Verificar se é erro de nome duplicado
		if strings.Contains(err.Error(), "UK_TEMPLATEEMAIL_NOME") {
			return ErrNomeDuplicado
		}
		return fmt.Errorf("erro ao atualizar template: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar atualização: %w", err)
	}

	if rowsAffected == 0 {
		return ErrTemplateNaoEncontrado
	}

	r.logger.Info("Template atualizado com sucesso",
		zap.Int64("id", template.ID),
		zap.String("nome", template.Nome))

	return nil
}

// Delete realiza soft delete do template (marca como inativo)
func (r *Repository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE TEMPLATEEMAIL
		SET ATIVO = 0,
			DATA_ATUALIZACAO = SYSDATE
		WHERE ID = :1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir template: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar exclusão: %w", err)
	}

	if rowsAffected == 0 {
		return ErrTemplateNaoEncontrado
	}

	r.logger.Info("Template excluído (soft delete)",
		zap.Int64("id", id))

	return nil
}

// GetByID busca um template por ID
func (r *Repository) GetByID(ctx context.Context, id int64) (*Template, error) {
	query := `
		SELECT
			ID, NOME, DESCRICAO, HEADER_HTML, BODY_HTML,
			FOOTER_HTML, ASSUNTO_PADRAO, ATIVO,
			DATA_CRIACAO, DATA_ATUALIZACAO, CRIADO_POR
		FROM TEMPLATEEMAIL
		WHERE ID = :1`

	var t Template
	var ativo int

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.Nome, &t.Descricao, &t.HeaderHTML, &t.BodyHTML,
		&t.FooterHTML, &t.AssuntoPadrao, &ativo,
		&t.DataCriacao, &t.DataAtualizacao, &t.CriadoPor,
	)

	if err == sql.ErrNoRows {
		return nil, ErrTemplateNaoEncontrado
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar template: %w", err)
	}

	t.Ativo = ativo == 1

	return &t, nil
}

// GetByNome busca um template por nome
func (r *Repository) GetByNome(ctx context.Context, nome string) (*Template, error) {
	query := `
		SELECT
			ID, NOME, DESCRICAO, HEADER_HTML, BODY_HTML,
			FOOTER_HTML, ASSUNTO_PADRAO, ATIVO,
			DATA_CRIACAO, DATA_ATUALIZACAO, CRIADO_POR
		FROM TEMPLATEEMAIL
		WHERE NOME = :1`

	var t Template
	var ativo int

	err := r.db.QueryRowContext(ctx, query, nome).Scan(
		&t.ID, &t.Nome, &t.Descricao, &t.HeaderHTML, &t.BodyHTML,
		&t.FooterHTML, &t.AssuntoPadrao, &ativo,
		&t.DataCriacao, &t.DataAtualizacao, &t.CriadoPor,
	)

	if err == sql.ErrNoRows {
		return nil, ErrTemplateNaoEncontrado
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar template por nome: %w", err)
	}

	t.Ativo = ativo == 1

	return &t, nil
}

// List retorna uma lista paginada de templates
func (r *Repository) List(ctx context.Context, page, limit int, searchTerm string) ([]Template, error) {
	offset := (page - 1) * limit

	query := `
		SELECT
			ID, NOME, DESCRICAO, HEADER_HTML, BODY_HTML,
			FOOTER_HTML, ASSUNTO_PADRAO, ATIVO,
			DATA_CRIACAO, DATA_ATUALIZACAO, CRIADO_POR
		FROM TEMPLATEEMAIL
		WHERE 1=1`

	var args []interface{}
	argPos := 1

	// Adicionar filtro de busca se fornecido
	if searchTerm != "" {
		query += fmt.Sprintf(" AND (UPPER(NOME) LIKE :%d OR UPPER(DESCRICAO) LIKE :%d)", argPos, argPos+1)
		searchPattern := "%" + strings.ToUpper(searchTerm) + "%"
		args = append(args, searchPattern, searchPattern)
		argPos += 2
	}

	query += " ORDER BY DATA_CRIACAO DESC"
	query += fmt.Sprintf(" OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", offset, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar templates: %w", err)
	}
	defer rows.Close()

	var templates []Template
	for rows.Next() {
		var t Template
		var ativo int

		err := rows.Scan(
			&t.ID, &t.Nome, &t.Descricao, &t.HeaderHTML, &t.BodyHTML,
			&t.FooterHTML, &t.AssuntoPadrao, &ativo,
			&t.DataCriacao, &t.DataAtualizacao, &t.CriadoPor,
		)
		if err != nil {
			r.logger.Error("Erro ao escanear template", zap.Error(err))
			continue
		}

		t.Ativo = ativo == 1
		templates = append(templates, t)
	}

	return templates, rows.Err()
}

// ListActive retorna apenas templates ativos
func (r *Repository) ListActive(ctx context.Context) ([]Template, error) {
	query := `
		SELECT
			ID, NOME, DESCRICAO, HEADER_HTML, BODY_HTML,
			FOOTER_HTML, ASSUNTO_PADRAO, ATIVO,
			DATA_CRIACAO, DATA_ATUALIZACAO, CRIADO_POR
		FROM TEMPLATEEMAIL
		WHERE ATIVO = 1
		ORDER BY NOME ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar templates ativos: %w", err)
	}
	defer rows.Close()

	var templates []Template
	for rows.Next() {
		var t Template
		var ativo int

		err := rows.Scan(
			&t.ID, &t.Nome, &t.Descricao, &t.HeaderHTML, &t.BodyHTML,
			&t.FooterHTML, &t.AssuntoPadrao, &ativo,
			&t.DataCriacao, &t.DataAtualizacao, &t.CriadoPor,
		)
		if err != nil {
			r.logger.Error("Erro ao escanear template ativo", zap.Error(err))
			continue
		}

		t.Ativo = ativo == 1
		templates = append(templates, t)
	}

	return templates, rows.Err()
}

// Count retorna o total de templates (para paginação)
func (r *Repository) Count(ctx context.Context, searchTerm string) (int64, error) {
	query := "SELECT COUNT(*) FROM TEMPLATEEMAIL WHERE 1=1"

	var args []interface{}
	argPos := 1

	// Adicionar filtro de busca se fornecido
	if searchTerm != "" {
		query += fmt.Sprintf(" AND (UPPER(NOME) LIKE :%d OR UPPER(DESCRICAO) LIKE :%d)", argPos, argPos+1)
		searchPattern := "%" + strings.ToUpper(searchTerm) + "%"
		args = append(args, searchPattern, searchPattern)
	}

	var count int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("erro ao contar templates: %w", err)
	}

	return count, nil
}

// Duplicate duplica um template existente
func (r *Repository) Duplicate(ctx context.Context, id int64, newName string) (int64, error) {
	// Buscar template original
	original, err := r.GetByID(ctx, id)
	if err != nil {
		return 0, err
	}

	// Criar novo template com dados do original
	newTemplate := &Template{
		Nome:          newName,
		Descricao:     original.Descricao,
		HeaderHTML:    original.HeaderHTML,
		BodyHTML:      original.BodyHTML,
		FooterHTML:    original.FooterHTML,
		AssuntoPadrao: original.AssuntoPadrao,
		Ativo:         original.Ativo,
		CriadoPor:     original.CriadoPor,
	}

	// Criar o novo template
	newID, err := r.Create(ctx, newTemplate)
	if err != nil {
		return 0, fmt.Errorf("erro ao duplicar template: %w", err)
	}

	r.logger.Info("Template duplicado com sucesso",
		zap.Int64("original_id", id),
		zap.Int64("novo_id", newID),
		zap.String("novo_nome", newName))

	return newID, nil
}

// CheckInUse verifica se um template está sendo usado por algum e-mail
func (r *Repository) CheckInUse(ctx context.Context, id int64) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM MENSAGEMEMAIL
		WHERE TEMPLATE_ID = :1`

	var count int64
	err := r.db.QueryRowContext(ctx, query, id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar uso do template: %w", err)
	}

	return count > 0, nil
}

// boolToInt converte bool para int (1 ou 0)
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
