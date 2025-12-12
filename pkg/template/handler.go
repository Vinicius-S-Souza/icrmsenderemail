package template

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Handler gerencia as requisições HTTP para templates
type Handler struct {
	repo           *Repository
	macroProcessor *MacroProcessor
	logger         *zap.Logger
}

// NewHandler cria uma nova instância do handler
func NewHandler(repo *Repository, macroProcessor *MacroProcessor, logger *zap.Logger) *Handler {
	return &Handler{
		repo:           repo,
		macroProcessor: macroProcessor,
		logger:         logger,
	}
}

// Request/Response Structures

// CreateTemplateRequest representa a requisição para criar template
type CreateTemplateRequest struct {
	Nome          string `json:"nome"`
	Descricao     string `json:"descricao"`
	HeaderHTML    string `json:"headerHtml"`
	BodyHTML      string `json:"bodyHtml"`
	FooterHTML    string `json:"footerHtml"`
	AssuntoPadrao string `json:"assuntoPadrao"`
	Ativo         bool   `json:"ativo"`
	CriadoPor     string `json:"criadoPor"`
}

// UpdateTemplateRequest representa a requisição para atualizar template
type UpdateTemplateRequest struct {
	Nome          string `json:"nome"`
	Descricao     string `json:"descricao"`
	HeaderHTML    string `json:"headerHtml"`
	BodyHTML      string `json:"bodyHtml"`
	FooterHTML    string `json:"footerHtml"`
	AssuntoPadrao string `json:"assuntoPadrao"`
	Ativo         bool   `json:"ativo"`
	CriadoPor     string `json:"criadoPor"`
}

// TemplateResponse representa a resposta com dados do template
type TemplateResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Data    TemplateDTO `json:"data,omitempty"`
}

// TemplateListResponse representa a resposta com lista de templates
type TemplateListResponse struct {
	Success    bool          `json:"success"`
	Error      string        `json:"error,omitempty"`
	Data       []TemplateDTO `json:"data,omitempty"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int           `json:"totalPages"`
}

// MacrosResponse representa a resposta com lista de macros disponíveis
type MacrosResponse struct {
	Success bool    `json:"success"`
	Data    []Macro `json:"data"`
}

// PreviewRequest representa a requisição para preview
type PreviewRequest struct {
	HeaderHTML string `json:"headerHtml"`
	BodyHTML   string `json:"bodyHtml"`
	FooterHTML string `json:"footerHtml"`
	UseSampleData bool `json:"useSampleData"`
}

// PreviewResponse representa a resposta do preview
type PreviewResponse struct {
	Success bool   `json:"success"`
	HTML    string `json:"html"`
	Error   string `json:"error,omitempty"`
}

// DuplicateRequest representa a requisição para duplicar template
type DuplicateRequest struct {
	NewName string `json:"newName"`
}

// API Handlers

// ListTemplates retorna lista paginada de templates
func (h *Handler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Método não permitido"})
		return
	}

	// Parâmetros de paginação
	page := getQueryInt(r, "page", 1)
	limit := getQueryInt(r, "limit", 10)
	searchTerm := r.URL.Query().Get("search")
	activeOnly := r.URL.Query().Get("activeOnly") == "true"

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var templates []Template
	var err error
	var total int64

	if activeOnly {
		// Listar apenas ativos (sem paginação)
		templates, err = h.repo.ListActive(ctx)
		total = int64(len(templates))
	} else {
		// Listar com paginação
		templates, err = h.repo.List(ctx, page, limit, searchTerm)
		if err != nil {
			h.logger.Error("Erro ao listar templates", zap.Error(err))
			respondJSON(w, http.StatusInternalServerError, TemplateListResponse{
				Success: false,
				Error:   "Erro ao buscar templates",
			})
			return
		}

		// Contar total
		total, err = h.repo.Count(ctx, searchTerm)
		if err != nil {
			h.logger.Error("Erro ao contar templates", zap.Error(err))
		}
	}

	if err != nil {
		h.logger.Error("Erro ao listar templates", zap.Error(err))
		respondJSON(w, http.StatusInternalServerError, TemplateListResponse{
			Success: false,
			Error:   "Erro ao buscar templates",
		})
		return
	}

	// Converter para DTOs
	dtos := make([]TemplateDTO, len(templates))
	for i, t := range templates {
		dtos[i] = t.ToDTO()
	}

	// Calcular total de páginas
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	respondJSON(w, http.StatusOK, TemplateListResponse{
		Success:    true,
		Data:       dtos,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

// GetTemplate retorna um template específico por ID
func (h *Handler) GetTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Método não permitido"})
		return
	}

	// Extrair ID da URL
	id, err := extractIDFromPath(r.URL.Path, "/api/templates/")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, TemplateResponse{
			Success: false,
			Error:   "ID inválido",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	template, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrTemplateNaoEncontrado {
			respondJSON(w, http.StatusNotFound, TemplateResponse{
				Success: false,
				Error:   "Template não encontrado",
			})
			return
		}
		h.logger.Error("Erro ao buscar template", zap.Error(err), zap.Int64("id", id))
		respondJSON(w, http.StatusInternalServerError, TemplateResponse{
			Success: false,
			Error:   "Erro ao buscar template",
		})
		return
	}

	respondJSON(w, http.StatusOK, TemplateResponse{
		Success: true,
		Data:    template.ToDTO(),
	})
}

// CreateTemplate cria um novo template
func (h *Handler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Método não permitido"})
		return
	}

	var req CreateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Erro ao decodificar requisição", zap.Error(err))
		respondJSON(w, http.StatusBadRequest, TemplateResponse{
			Success: false,
			Error:   "Requisição inválida",
		})
		return
	}

	// Criar estrutura Template
	template := &Template{
		Nome:          strings.TrimSpace(req.Nome),
		Descricao:     sql.NullString{String: strings.TrimSpace(req.Descricao), Valid: req.Descricao != ""},
		HeaderHTML:    sql.NullString{String: req.HeaderHTML, Valid: req.HeaderHTML != ""},
		BodyHTML:      req.BodyHTML,
		FooterHTML:    sql.NullString{String: req.FooterHTML, Valid: req.FooterHTML != ""},
		AssuntoPadrao: sql.NullString{String: strings.TrimSpace(req.AssuntoPadrao), Valid: req.AssuntoPadrao != ""},
		Ativo:         req.Ativo,
		CriadoPor:     sql.NullString{String: strings.TrimSpace(req.CriadoPor), Valid: req.CriadoPor != ""},
	}

	// Validar
	if err := template.Validate(); err != nil {
		respondJSON(w, http.StatusBadRequest, TemplateResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Criar no banco
	id, err := h.repo.Create(ctx, template)
	if err != nil {
		if err == ErrNomeDuplicado {
			respondJSON(w, http.StatusConflict, TemplateResponse{
				Success: false,
				Error:   "Já existe um template com este nome",
			})
			return
		}
		h.logger.Error("Erro ao criar template", zap.Error(err))
		respondJSON(w, http.StatusInternalServerError, TemplateResponse{
			Success: false,
			Error:   "Erro ao criar template",
		})
		return
	}

	// Buscar template criado
	template.ID = id
	createdTemplate, _ := h.repo.GetByID(ctx, id)

	h.logger.Info("Template criado com sucesso",
		zap.Int64("id", id),
		zap.String("nome", template.Nome))

	respondJSON(w, http.StatusCreated, TemplateResponse{
		Success: true,
		Data:    createdTemplate.ToDTO(),
	})
}

// UpdateTemplate atualiza um template existente
func (h *Handler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Método não permitido"})
		return
	}

	// Extrair ID da URL
	id, err := extractIDFromPath(r.URL.Path, "/api/templates/")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, TemplateResponse{
			Success: false,
			Error:   "ID inválido",
		})
		return
	}

	var req UpdateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Erro ao decodificar requisição", zap.Error(err))
		respondJSON(w, http.StatusBadRequest, TemplateResponse{
			Success: false,
			Error:   "Requisição inválida",
		})
		return
	}

	// Criar estrutura Template
	template := &Template{
		ID:            id,
		Nome:          strings.TrimSpace(req.Nome),
		Descricao:     sql.NullString{String: strings.TrimSpace(req.Descricao), Valid: req.Descricao != ""},
		HeaderHTML:    sql.NullString{String: req.HeaderHTML, Valid: req.HeaderHTML != ""},
		BodyHTML:      req.BodyHTML,
		FooterHTML:    sql.NullString{String: req.FooterHTML, Valid: req.FooterHTML != ""},
		AssuntoPadrao: sql.NullString{String: strings.TrimSpace(req.AssuntoPadrao), Valid: req.AssuntoPadrao != ""},
		Ativo:         req.Ativo,
		CriadoPor:     sql.NullString{String: strings.TrimSpace(req.CriadoPor), Valid: req.CriadoPor != ""},
	}

	// Validar
	if err := template.Validate(); err != nil {
		respondJSON(w, http.StatusBadRequest, TemplateResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Atualizar no banco
	err = h.repo.Update(ctx, template)
	if err != nil {
		if err == ErrTemplateNaoEncontrado {
			respondJSON(w, http.StatusNotFound, TemplateResponse{
				Success: false,
				Error:   "Template não encontrado",
			})
			return
		}
		if err == ErrNomeDuplicado {
			respondJSON(w, http.StatusConflict, TemplateResponse{
				Success: false,
				Error:   "Já existe um template com este nome",
			})
			return
		}
		h.logger.Error("Erro ao atualizar template", zap.Error(err), zap.Int64("id", id))
		respondJSON(w, http.StatusInternalServerError, TemplateResponse{
			Success: false,
			Error:   "Erro ao atualizar template",
		})
		return
	}

	// Buscar template atualizado
	updatedTemplate, _ := h.repo.GetByID(ctx, id)

	h.logger.Info("Template atualizado com sucesso", zap.Int64("id", id))

	respondJSON(w, http.StatusOK, TemplateResponse{
		Success: true,
		Data:    updatedTemplate.ToDTO(),
	})
}

// DeleteTemplate exclui um template (soft delete)
func (h *Handler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Método não permitido"})
		return
	}

	// Extrair ID da URL
	id, err := extractIDFromPath(r.URL.Path, "/api/templates/")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "ID inválido",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Verificar se está em uso
	inUse, err := h.repo.CheckInUse(ctx, id)
	if err != nil {
		h.logger.Error("Erro ao verificar uso do template", zap.Error(err), zap.Int64("id", id))
	}

	if inUse {
		respondJSON(w, http.StatusConflict, map[string]interface{}{
			"success": false,
			"error":   "Template está em uso e não pode ser excluído",
		})
		return
	}

	// Excluir (soft delete)
	err = h.repo.Delete(ctx, id)
	if err != nil {
		if err == ErrTemplateNaoEncontrado {
			respondJSON(w, http.StatusNotFound, map[string]interface{}{
				"success": false,
				"error":   "Template não encontrado",
			})
			return
		}
		h.logger.Error("Erro ao excluir template", zap.Error(err), zap.Int64("id", id))
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Erro ao excluir template",
		})
		return
	}

	h.logger.Info("Template excluído com sucesso", zap.Int64("id", id))

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Template excluído com sucesso",
	})
}

// GetMacros retorna lista de macros disponíveis
func (h *Handler) GetMacros(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Método não permitido"})
		return
	}

	respondJSON(w, http.StatusOK, MacrosResponse{
		Success: true,
		Data:    AvailableMacros,
	})
}

// PreviewTemplate gera preview de um template
func (h *Handler) PreviewTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Método não permitido"})
		return
	}

	var req PreviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Erro ao decodificar requisição de preview", zap.Error(err))
		respondJSON(w, http.StatusBadRequest, PreviewResponse{
			Success: false,
			Error:   "Requisição inválida",
		})
		return
	}

	// Usar dados de exemplo para preview
	sampleData := GetMacroPreviewData()

	// Processar cada seção
	var html strings.Builder

	if req.HeaderHTML != "" {
		headerProcessed := h.macroProcessor.ReplaceMacros(req.HeaderHTML, sampleData)
		html.WriteString(headerProcessed)
	}

	bodyProcessed := h.macroProcessor.ReplaceMacros(req.BodyHTML, sampleData)
	html.WriteString(bodyProcessed)

	if req.FooterHTML != "" {
		footerProcessed := h.macroProcessor.ReplaceMacros(req.FooterHTML, sampleData)
		html.WriteString(footerProcessed)
	}

	respondJSON(w, http.StatusOK, PreviewResponse{
		Success: true,
		HTML:    html.String(),
	})
}

// DuplicateTemplate duplica um template existente
func (h *Handler) DuplicateTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Método não permitido"})
		return
	}

	// Extrair ID da URL
	id, err := extractIDFromPath(r.URL.Path, "/api/templates/")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, TemplateResponse{
			Success: false,
			Error:   "ID inválido",
		})
		return
	}

	var req DuplicateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Erro ao decodificar requisição de duplicação", zap.Error(err))
		respondJSON(w, http.StatusBadRequest, TemplateResponse{
			Success: false,
			Error:   "Requisição inválida",
		})
		return
	}

	if strings.TrimSpace(req.NewName) == "" {
		respondJSON(w, http.StatusBadRequest, TemplateResponse{
			Success: false,
			Error:   "Nome do novo template é obrigatório",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Duplicar template
	newID, err := h.repo.Duplicate(ctx, id, req.NewName)
	if err != nil {
		if err == ErrTemplateNaoEncontrado {
			respondJSON(w, http.StatusNotFound, TemplateResponse{
				Success: false,
				Error:   "Template original não encontrado",
			})
			return
		}
		if err == ErrNomeDuplicado {
			respondJSON(w, http.StatusConflict, TemplateResponse{
				Success: false,
				Error:   "Já existe um template com este nome",
			})
			return
		}
		h.logger.Error("Erro ao duplicar template", zap.Error(err), zap.Int64("id", id))
		respondJSON(w, http.StatusInternalServerError, TemplateResponse{
			Success: false,
			Error:   "Erro ao duplicar template",
		})
		return
	}

	// Buscar template duplicado
	duplicated, _ := h.repo.GetByID(ctx, newID)

	h.logger.Info("Template duplicado com sucesso",
		zap.Int64("original_id", id),
		zap.Int64("novo_id", newID))

	respondJSON(w, http.StatusCreated, TemplateResponse{
		Success: true,
		Data:    duplicated.ToDTO(),
	})
}

// Utility functions

// respondJSON envia uma resposta JSON
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// getQueryInt extrai um parâmetro inteiro da query string
func getQueryInt(r *http.Request, key string, defaultValue int) int {
	valueStr := r.URL.Query().Get(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// extractIDFromPath extrai o ID numérico da URL
func extractIDFromPath(path, prefix string) (int64, error) {
	// Remove o prefixo
	idStr := strings.TrimPrefix(path, prefix)
	// Remove possível sufixo (como /duplicate)
	parts := strings.Split(idStr, "/")
	if len(parts) == 0 {
		return 0, ErrTemplateNaoEncontrado
	}
	return strconv.ParseInt(parts[0], 10, 64)
}
