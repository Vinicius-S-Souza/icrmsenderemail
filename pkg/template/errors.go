package template

import "errors"

// Erros do domínio de templates
var (
	ErrNomeObrigatorio     = errors.New("nome do template é obrigatório")
	ErrNomeMuitoLongo      = errors.New("nome do template muito longo (máximo 100 caracteres)")
	ErrBodyObrigatorio     = errors.New("corpo do template é obrigatório")
	ErrTemplateNaoEncontrado = errors.New("template não encontrado")
	ErrNomeDuplicado       = errors.New("já existe um template com este nome")
	ErrTemplateEmUso       = errors.New("template está em uso e não pode ser excluído")
)
