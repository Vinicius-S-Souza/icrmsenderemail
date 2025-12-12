package template

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/cliente"
	"go.uber.org/zap"
)

// MacroProcessor processa substituição de macros em templates
type MacroProcessor struct {
	clienteRepo *cliente.Repository
	logger      *zap.Logger
	empresaNome string // Nome da empresa (configurável)
}

// NewMacroProcessor cria um novo processador de macros
func NewMacroProcessor(clienteRepo *cliente.Repository, empresaNome string, logger *zap.Logger) *MacroProcessor {
	return &MacroProcessor{
		clienteRepo: clienteRepo,
		empresaNome: empresaNome,
		logger:      logger,
	}
}

// ReplaceMacros substitui todas as macros no conteúdo com os dados fornecidos
func (mp *MacroProcessor) ReplaceMacros(content string, data MacroData) string {
	result := content

	// Substituir cada macro
	result = strings.ReplaceAll(result, "{{nome}}", data.Nome)
	result = strings.ReplaceAll(result, "{{email}}", data.Email)
	result = strings.ReplaceAll(result, "{{cpf_cnpj}}", data.CpfCnpj)
	result = strings.ReplaceAll(result, "{{codigo}}", data.Codigo)
	result = strings.ReplaceAll(result, "{{data}}", data.Data)
	result = strings.ReplaceAll(result, "{{hora}}", data.Hora)
	result = strings.ReplaceAll(result, "{{data_hora}}", data.DataHora)
	result = strings.ReplaceAll(result, "{{empresa}}", data.Empresa)
	result = strings.ReplaceAll(result, "{{ano}}", data.Ano)

	// Substituir macros personalizadas se houver
	if data.CustomData != nil {
		for key, value := range data.CustomData {
			placeholder := fmt.Sprintf("{{%s}}", key)
			result = strings.ReplaceAll(result, placeholder, value)
		}
	}

	return result
}

// GetMacroDataFromCliente busca dados do cliente e popula MacroData
func (mp *MacroProcessor) GetMacroDataFromCliente(ctx context.Context, cliCodigo int) (*MacroData, error) {
	// Buscar cliente no banco
	cli, err := mp.clienteRepo.FindByCodigo(ctx, cliCodigo)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar cliente para macros: %w", err)
	}

	// Obter data/hora atual
	now := time.Now()

	// Popular MacroData
	data := &MacroData{
		Nome:     cli.CliNome,
		Email:    cli.Email,
		CpfCnpj:  cli.CliCpfCnpj,
		Codigo:   fmt.Sprintf("%d", cli.CliCodigo),
		Data:     now.Format("02/01/2006"),
		Hora:     now.Format("15:04"),
		DataHora: now.Format("02/01/2006 15:04"),
		Empresa:  mp.empresaNome,
		Ano:      now.Format("2006"),
	}

	mp.logger.Debug("MacroData gerado para cliente",
		zap.Int("cliCodigo", cliCodigo),
		zap.String("nome", data.Nome),
		zap.String("email", data.Email))

	return data, nil
}

// GetDefaultMacroData retorna MacroData com valores padrão (quando não há cliente)
func (mp *MacroProcessor) GetDefaultMacroData() *MacroData {
	now := time.Now()

	return &MacroData{
		Nome:     "",
		Email:    "",
		CpfCnpj:  "",
		Codigo:   "",
		Data:     now.Format("02/01/2006"),
		Hora:     now.Format("15:04"),
		DataHora: now.Format("02/01/2006 15:04"),
		Empresa:  mp.empresaNome,
		Ano:      now.Format("2006"),
	}
}

// ExtractUsedMacros extrai todas as macros usadas no conteúdo
func ExtractUsedMacros(content string) []string {
	// Regex para encontrar padrões {{palavra}}
	re := regexp.MustCompile(`\{\{([^}]+)\}\}`)
	matches := re.FindAllStringSubmatch(content, -1)

	seen := make(map[string]bool)
	var macros []string

	for _, match := range matches {
		if len(match) > 1 {
			macro := "{{" + match[1] + "}}"
			if !seen[macro] {
				seen[macro] = true
				macros = append(macros, macro)
			}
		}
	}

	return macros
}

// ValidateMacros verifica se todas as macros usadas são válidas
func ValidateMacros(content string) (bool, []string) {
	usedMacros := ExtractUsedMacros(content)

	// Criar mapa de macros válidas
	validMacros := make(map[string]bool)
	for _, macro := range AvailableMacros {
		validMacros[macro.Key] = true
	}

	// Verificar macros inválidas
	var invalidMacros []string
	for _, macro := range usedMacros {
		if !validMacros[macro] {
			invalidMacros = append(invalidMacros, macro)
		}
	}

	return len(invalidMacros) == 0, invalidMacros
}

// ProcessTemplate processa um template completo substituindo macros
func (mp *MacroProcessor) ProcessTemplate(ctx context.Context, template *Template, cliCodigo sql.NullInt64) (string, string, error) {
	var macroData *MacroData
	var err error

	// Obter dados do cliente se fornecido
	if cliCodigo.Valid {
		macroData, err = mp.GetMacroDataFromCliente(ctx, int(cliCodigo.Int64))
		if err != nil {
			mp.logger.Warn("Erro ao buscar dados do cliente, usando valores padrão",
				zap.Error(err),
				zap.Int64("cliCodigo", cliCodigo.Int64))
			macroData = mp.GetDefaultMacroData()
		}
	} else {
		macroData = mp.GetDefaultMacroData()
	}

	// Processar cada seção do template
	var fullHTML strings.Builder

	// Header
	if template.HeaderHTML.Valid && template.HeaderHTML.String != "" {
		headerProcessed := mp.ReplaceMacros(template.HeaderHTML.String, *macroData)
		fullHTML.WriteString(headerProcessed)
	}

	// Body
	bodyProcessed := mp.ReplaceMacros(template.BodyHTML, *macroData)
	fullHTML.WriteString(bodyProcessed)

	// Footer
	if template.FooterHTML.Valid && template.FooterHTML.String != "" {
		footerProcessed := mp.ReplaceMacros(template.FooterHTML.String, *macroData)
		fullHTML.WriteString(footerProcessed)
	}

	// Processar assunto
	var assunto string
	if template.AssuntoPadrao.Valid && template.AssuntoPadrao.String != "" {
		assunto = mp.ReplaceMacros(template.AssuntoPadrao.String, *macroData)
	}

	mp.logger.Debug("Template processado com macros substituídas",
		zap.Int64("templateId", template.ID),
		zap.String("templateNome", template.Nome),
		zap.Bool("temCliente", cliCodigo.Valid))

	// Retornar assunto primeiro, depois corpo (ordem esperada pelos handlers)
	return assunto, fullHTML.String(), nil
}

// GetMacroPreviewData retorna dados de exemplo para preview de templates
func GetMacroPreviewData() MacroData {
	now := time.Now()

	return MacroData{
		Nome:     "João da Silva",
		Email:    "joao@exemplo.com",
		CpfCnpj:  "123.456.789-00",
		Codigo:   "12345",
		Data:     now.Format("02/01/2006"),
		Hora:     now.Format("15:04"),
		DataHora: now.Format("02/01/2006 15:04"),
		Empresa:  "Minha Empresa Ltda",
		Ano:      now.Format("2006"),
	}
}
