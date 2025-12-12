package message

import (
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/zap"
)

// Repository gerencia operações de banco de dados para emails
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

// GetPendingEmails busca emails pendentes para envio
func (r *Repository) GetPendingEmails(ctx context.Context, limit, daysOffset, maxTentativas int) ([]Email, error) {
	query := `
		SELECT
			ID, CLICODIGO, REMETENTE, DESTINATARIO, ASSUNTO,
			CORPO, TIPO_CORPO, STATUS_ENVIO, DATA_CADASTRO,
			DATA_AGENDAMENTO, DATA_ENVIO, QTD_TENTATIVAS,
			DETALHES_ERRO, ID_PROVIDER, METODO_ENVIO, PRIORIDADE,
			ANEXO_REFERENCIA, ANEXO_NOME, ANEXO_TIPO, IP_ORIGEM, TEMPLATE_ID
		FROM MENSAGEMEMAIL
		WHERE STATUS_ENVIO = 0
		  AND QTD_TENTATIVAS < :1
		  AND (DATA_AGENDAMENTO IS NULL OR DATA_AGENDAMENTO <= SYSDATE + :2)
		ORDER BY PRIORIDADE ASC, DATA_CADASTRO ASC
		FETCH FIRST :3 ROWS ONLY`

	rows, err := r.db.QueryContext(ctx, query, maxTentativas, daysOffset, limit)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar emails pendentes: %w", err)
	}
	defer rows.Close()

	var emails []Email
	for rows.Next() {
		var e Email
		err := rows.Scan(
			&e.ID, &e.CliCodigo, &e.Remetente, &e.Destinatario, &e.Assunto,
			&e.Corpo, &e.TipoCorpo, &e.StatusEnvio, &e.DataCadastro,
			&e.DataAgendamento, &e.DataEnvio, &e.QTDTentativas,
			&e.DetalhesErro, &e.IDProvider, &e.MetodoEnvio, &e.Prioridade,
			&e.AnexoReferencia, &e.AnexoNome, &e.AnexoTipo, &e.IPOrigem, &e.TemplateID,
		)
		if err != nil {
			r.logger.Error("Erro ao escanear email", zap.Error(err))
			continue
		}
		emails = append(emails, e)
	}

	return emails, rows.Err()
}

// MarkAsSent marca email como enviado com sucesso
func (r *Repository) MarkAsSent(ctx context.Context, id int64, providerID string, metodo int) error {
	query := `
		UPDATE MENSAGEMEMAIL 
		SET STATUS_ENVIO = 2,
			DATA_ENVIO = SYSDATE,
			ID_PROVIDER = :1,
			METODO_ENVIO = :2,
			QTD_TENTATIVAS = QTD_TENTATIVAS + 1,
			DETALHES_ERRO = NULL
		WHERE ID = :3`

	_, err := r.db.ExecContext(ctx, query, providerID, metodo, id)
	if err != nil {
		return fmt.Errorf("erro ao marcar email como enviado: %w", err)
	}

	r.logger.Debug("Email marcado como enviado", zap.Int64("id", id))
	return nil
}

// MarkAsError marca email com erro temporário
func (r *Repository) MarkAsError(ctx context.Context, id int64, errorMsg string, metodo int) error {
	query := `
		UPDATE MENSAGEMEMAIL 
		SET STATUS_ENVIO = 3,
			QTD_TENTATIVAS = QTD_TENTATIVAS + 1,
			DETALHES_ERRO = :1,
			METODO_ENVIO = :2
		WHERE ID = :3`

	_, err := r.db.ExecContext(ctx, query, errorMsg, metodo, id)
	if err != nil {
		return fmt.Errorf("erro ao marcar email com erro: %w", err)
	}

	r.logger.Debug("Email marcado com erro", zap.Int64("id", id))
	return nil
}

// MarkAsInvalid marca email como inválido
func (r *Repository) MarkAsInvalid(ctx context.Context, id int64, errorMsg string, metodo int) error {
	query := `
		UPDATE MENSAGEMEMAIL 
		SET STATUS_ENVIO = 125,
			QTD_TENTATIVAS = QTD_TENTATIVAS + 1,
			DETALHES_ERRO = :1,
			METODO_ENVIO = :2
		WHERE ID = :3`

	_, err := r.db.ExecContext(ctx, query, errorMsg, metodo, id)
	if err != nil {
		return fmt.Errorf("erro ao marcar email como inválido: %w", err)
	}

	r.logger.Debug("Email marcado como inválido", zap.Int64("id", id))
	return nil
}

// MarkAsPermanentFailure marca email com falha permanente
func (r *Repository) MarkAsPermanentFailure(ctx context.Context, id int64, errorMsg string, metodo int) error {
	query := `
		UPDATE MENSAGEMEMAIL 
		SET STATUS_ENVIO = 4,
			QTD_TENTATIVAS = QTD_TENTATIVAS + 1,
			DETALHES_ERRO = :1,
			METODO_ENVIO = :2
		WHERE ID = :3`

	_, err := r.db.ExecContext(ctx, query, errorMsg, metodo, id)
	if err != nil {
		return fmt.Errorf("erro ao marcar email como falha permanente: %w", err)
	}

	r.logger.Debug("Email marcado como falha permanente", zap.Int64("id", id))
	return nil
}

// GetByID busca um email por ID
func (r *Repository) GetByID(ctx context.Context, id int64) (*Email, error) {
	query := `
		SELECT
			ID, CLICODIGO, REMETENTE, DESTINATARIO, ASSUNTO,
			CORPO, TIPO_CORPO, STATUS_ENVIO, DATA_CADASTRO,
			DATA_AGENDAMENTO, DATA_ENVIO, QTD_TENTATIVAS,
			DETALHES_ERRO, ID_PROVIDER, METODO_ENVIO, PRIORIDADE,
			ANEXO_REFERENCIA, ANEXO_NOME, ANEXO_TIPO, IP_ORIGEM, TEMPLATE_ID
		FROM MENSAGEMEMAIL
		WHERE ID = :1`

	var e Email
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&e.ID, &e.CliCodigo, &e.Remetente, &e.Destinatario, &e.Assunto,
		&e.Corpo, &e.TipoCorpo, &e.StatusEnvio, &e.DataCadastro,
		&e.DataAgendamento, &e.DataEnvio, &e.QTDTentativas,
		&e.DetalhesErro, &e.IDProvider, &e.MetodoEnvio, &e.Prioridade,
		&e.AnexoReferencia, &e.AnexoNome, &e.AnexoTipo, &e.IPOrigem, &e.TemplateID,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("email não encontrado: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar email: %w", err)
	}

	return &e, nil
}

// InsertEmail insere um novo email no banco
func (r *Repository) InsertEmail(ctx context.Context, email *Email) (int64, error) {
	query := `
		INSERT INTO MENSAGEMEMAIL (
			ID, CLICODIGO, REMETENTE, DESTINATARIO, ASSUNTO,
			CORPO, TIPO_CORPO, STATUS_ENVIO, DATA_CADASTRO,
			DATA_AGENDAMENTO, PRIORIDADE, IP_ORIGEM,
			ANEXO_REFERENCIA, ANEXO_NOME, ANEXO_TIPO, TEMPLATE_ID
		) VALUES (
			SEQ_MENSAGEMEMAIL.NEXTVAL, :1, :2, :3, :4,
			:5, :6, :7, SYSDATE,
			:8, :9, :10,
			:11, :12, :13, :14
		) RETURNING ID INTO :15`

	var id int64
	_, err := r.db.ExecContext(ctx, query,
		email.CliCodigo, email.Remetente, email.Destinatario, email.Assunto,
		email.Corpo, email.TipoCorpo, int(email.StatusEnvio), // Convert EmailStatus to int
		email.DataAgendamento, email.Prioridade, email.IPOrigem,
		email.AnexoReferencia, email.AnexoNome, email.AnexoTipo, email.TemplateID,
		sql.Out{Dest: &id},
	)

	if err != nil {
		return 0, fmt.Errorf("erro ao inserir email: %w", err)
	}

	r.logger.Info("Email inserido com sucesso", zap.Int64("id", id))
	return id, nil
}

// CountPendingEmails conta emails pendentes
func (r *Repository) CountPendingEmails(ctx context.Context, daysOffset, maxTentativas int) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM MENSAGEMEMAIL
		WHERE STATUS_ENVIO = 0
		  AND QTD_TENTATIVAS < :1
		  AND (DATA_AGENDAMENTO IS NULL OR DATA_AGENDAMENTO <= SYSDATE + :2)`

	var count int64
	err := r.db.QueryRowContext(ctx, query, maxTentativas, daysOffset).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("erro ao contar emails pendentes: %w", err)
	}

	return count, nil
}

// GetStats retorna estatísticas de emails
func (r *Repository) GetStats(ctx context.Context) (map[string]int64, error) {
	query := `
		SELECT 
			STATUS_ENVIO,
			COUNT(*) as total
		FROM MENSAGEMEMAIL
		WHERE DATA_CADASTRO >= TRUNC(SYSDATE)
		GROUP BY STATUS_ENVIO`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar estatísticas: %w", err)
	}
	defer rows.Close()

	stats := make(map[string]int64)
	for rows.Next() {
		var status int
		var total int64
		if err := rows.Scan(&status, &total); err != nil {
			continue
		}
		statusKey := fmt.Sprintf("status_%d", status)
		stats[statusKey] = total
	}

	return stats, rows.Err()
}
