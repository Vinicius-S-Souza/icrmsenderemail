# Changelog

Todas as mudanças notáveis neste projeto serão documentadas neste arquivo.

O formato é baseado em [Keep a Changelog](https://keepachangelog.com/pt-BR/1.0.0/),
e este projeto adere ao [Versionamento Semântico](https://semver.org/lang/pt-BR/).

## [1.1.0] - 11/12/2025 16:50

### ✨ Adicionado
- **Suporte completo a anexos** via base64 (SendGrid, Pontaltech) e URL pública (Zenvia)
- Campo `URL` na estrutura `Attachment` para anexos via URL
- Campo `AttachmentURL` na API de disparo manual
- Validação de URL de anexos no backend e frontend
- Detecção automática de tipo de anexo (URL vs base64) no processor
- Interface web inteligente que adapta campos de anexo baseado no provider ativo
- Função `toggleAnexoFields()` para mostrar/ocultar campos apropriados
- Função `handleUrlInput()` para validação em tempo real de URLs
- Documentação completa sobre limitações de anexos no Zenvia ([ZENVIA_ANEXOS.md](ZENVIA_ANEXOS.md))
- Instruções detalhadas para implementação frontend ([INSTRUCOES_FRONTEND_ANEXO_URL.md](INSTRUCOES_FRONTEND_ANEXO_URL.md))

### 🔧 Corrigido
- **E-mail "from" agora usa `default_from` do dbinit.ini** em todos os providers
- Provider Pontaltech agora extrai corretamente o ID da mensagem de `messages[0].id`
- Provider Pontaltech valida `invalidMessages[]` e trata como erro
- Provider Zenvia usa estrutura JSON correta com `type: "email"` e `subject` dentro de `contents`
- Provider Zenvia aceita anexos via `fileUrl` quando disponível
- Migração do banco de dados para alterar `ANEXO_REFERENCIA` de VARCHAR2(500) para CLOB

### 🗄️ Banco de Dados
- Alterado campo `ANEXO_REFERENCIA` de `VARCHAR2(500)` para `CLOB` (suporta arquivos grandes)
- Adicionado suporte a campo `ANEXO_TIPO` com valor "url" para diferenciar anexos por URL
- Script de migração: [sql/alter_anexo_referencia_to_clob.sql](sql/alter_anexo_referencia_to_clob.sql)
- Utilitário de migração: [cmd/migrate/main.go](cmd/migrate/main.go)

### 📝 Mudanças
- Estrutura `PontaltechEmailResponse` atualizada para formato correto da API
- Estrutura `ZenviaEmailContent` atualizada com campos corretos
- Provider Zenvia ignora anexos base64 com aviso no log
- Campo `Attachment.Data` agora exclusivo para base64
- Campo `Attachment.URL` para anexos via URL pública

### 📦 Arquivos Modificados
- `pkg/email/sender.go` - Adicionado campo URL em Attachment
- `pkg/manual/handler.go` - Suporte a AttachmentURL com validação
- `pkg/message/processor.go` - Detecção automática de tipo de anexo
- `pkg/message/repository.go` - InsertEmail agora salva campos de anexo
- `pkg/email/sendgrid_provider.go` - Usa default_from configurado
- `pkg/email/pontaltech_provider.go` - Estrutura de resposta corrigida e extração de ID
- `pkg/email/zenvia_provider.go` - Estrutura JSON correta e suporte a fileUrl
- `pkg/manual/html.go` - Interface adaptativa baseada em provider
- `sql/create_table_mensagememail.sql` - ANEXO_REFERENCIA como CLOB

### 📚 Documentação
- Criado [ZENVIA_ANEXOS.md](ZENVIA_ANEXOS.md) - Limitações e comparação de providers
- Criado [INSTRUCOES_FRONTEND_ANEXO_URL.md](INSTRUCOES_FRONTEND_ANEXO_URL.md) - Guia de implementação
- Atualizado [README.md](README.md) - Documentação de anexos e versão

### 🐛 Bugs Conhecidos
Nenhum no momento.

### ⚠️ Notas de Migração

#### De 1.0.0 para 1.1.0

1. **Banco de Dados:**
   ```bash
   # Execute a migração do banco
   go run ./cmd/migrate/main.go
   ```
   Ou execute manualmente:
   ```sql
   ALTER TABLE MENSAGEMEMAIL MODIFY (ANEXO_REFERENCIA CLOB);
   ```

2. **Configuração:**
   - Verifique que `default_from` está configurado no `dbinit.ini`
   - Exemplo: `default_from=noreply@seudominio.com.br`

3. **Zenvia:**
   - Se usar Zenvia com anexos, forneça URLs públicas
   - Anexos base64 serão ignorados com aviso no log

---

## [1.0.0] - 11/12/2025

### ✨ Versão Inicial
- Sistema completo de envio de e-mails
- Suporte a múltiplos providers (Mock, SMTP, SendGrid, Zenvia, Pontaltech)
- Dashboard web em tempo real
- Interface de disparo manual
- Processamento paralelo com workers
- Retry automático e circuit breaker
- Logs estruturados e métricas
- Health check HTTP
- Graceful shutdown
- Suporte a Windows Service
