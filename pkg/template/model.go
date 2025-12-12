package template

import (
	"database/sql"
	"time"
)

// Template representa um template de e-mail HTML
type Template struct {
	ID              int64
	Nome            string
	Descricao       sql.NullString
	HeaderHTML      sql.NullString
	BodyHTML        string
	FooterHTML      sql.NullString
	AssuntoPadrao   sql.NullString
	Ativo           bool
	DataCriacao     time.Time
	DataAtualizacao time.Time
	CriadoPor       sql.NullString
}

// Macro representa um placeholder substituível no template
type Macro struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	Example     string `json:"example"`
}

// AvailableMacros lista todas as macros disponíveis para uso em templates
var AvailableMacros = []Macro{
	{
		Key:         "{{nome}}",
		Description: "Nome completo do cliente",
		Example:     "João da Silva",
	},
	{
		Key:         "{{email}}",
		Description: "E-mail do cliente",
		Example:     "joao@exemplo.com",
	},
	{
		Key:         "{{cpf_cnpj}}",
		Description: "CPF ou CNPJ do cliente",
		Example:     "123.456.789-00",
	},
	{
		Key:         "{{codigo}}",
		Description: "Código do cliente no sistema",
		Example:     "12345",
	},
	{
		Key:         "{{data}}",
		Description: "Data atual no formato DD/MM/YYYY",
		Example:     "12/12/2025",
	},
	{
		Key:         "{{hora}}",
		Description: "Hora atual no formato HH:MM",
		Example:     "14:30",
	},
	{
		Key:         "{{data_hora}}",
		Description: "Data e hora atual no formato DD/MM/YYYY HH:MM",
		Example:     "12/12/2025 14:30",
	},
	{
		Key:         "{{empresa}}",
		Description: "Nome da empresa remetente",
		Example:     "Minha Empresa Ltda",
	},
	{
		Key:         "{{ano}}",
		Description: "Ano atual",
		Example:     "2025",
	},
}

// MacroData contém os dados para substituição de macros
type MacroData struct {
	Nome        string
	Email       string
	CpfCnpj     string
	Codigo      string
	Data        string
	Hora        string
	DataHora    string
	Empresa     string
	Ano         string
	CustomData  map[string]string // Campos personalizados adicionais
}

// TemplateDTO representa o template para transferência de dados (API)
type TemplateDTO struct {
	ID              int64  `json:"id"`
	Nome            string `json:"nome"`
	Descricao       string `json:"descricao"`
	HeaderHTML      string `json:"headerHtml"`
	BodyHTML        string `json:"bodyHtml"`
	FooterHTML      string `json:"footerHtml"`
	AssuntoPadrao   string `json:"assuntoPadrao"`
	Ativo           bool   `json:"ativo"`
	DataCriacao     string `json:"dataCriacao"`
	DataAtualizacao string `json:"dataAtualizacao"`
	CriadoPor       string `json:"criadoPor"`
}

// ToDTO converte Template para TemplateDTO
func (t *Template) ToDTO() TemplateDTO {
	return TemplateDTO{
		ID:              t.ID,
		Nome:            t.Nome,
		Descricao:       t.Descricao.String,
		HeaderHTML:      t.HeaderHTML.String,
		BodyHTML:        t.BodyHTML,
		FooterHTML:      t.FooterHTML.String,
		AssuntoPadrao:   t.AssuntoPadrao.String,
		Ativo:           t.Ativo,
		DataCriacao:     t.DataCriacao.Format("02/01/2006 15:04:05"),
		DataAtualizacao: t.DataAtualizacao.Format("02/01/2006 15:04:05"),
		CriadoPor:       t.CriadoPor.String,
	}
}

// GetFullHTML retorna o HTML completo concatenando header + body + footer
func (t *Template) GetFullHTML() string {
	var html string

	if t.HeaderHTML.Valid && t.HeaderHTML.String != "" {
		html += t.HeaderHTML.String
	}

	html += t.BodyHTML

	if t.FooterHTML.Valid && t.FooterHTML.String != "" {
		html += t.FooterHTML.String
	}

	return html
}

// Validate valida os campos do template
func (t *Template) Validate() error {
	if t.Nome == "" {
		return ErrNomeObrigatorio
	}
	if len(t.Nome) > 100 {
		return ErrNomeMuitoLongo
	}
	if t.BodyHTML == "" {
		return ErrBodyObrigatorio
	}
	return nil
}
