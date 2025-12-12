# Changelog

Todas as mudanças notáveis neste projeto serão documentadas neste arquivo.

O formato é baseado em [Keep a Changelog](https://keepachangelog.com/pt-BR/1.0.0/),
e este projeto adere ao [Versionamento Semântico](https://semver.org/lang/pt-BR/).

## [1.3.2] - 12/12/2025 23:45

### 🎨 Melhorado
- **Campos de anexo na página de disparo manual**
  - Aumentado o tamanho dos campos de anexo (file e URL)
  - Adicionado padding de 16px para melhor usabilidade
  - Borda tracejada colorida (dashed) para destacar área de upload
  - Background azul claro (#f8f9ff) para melhor visibilidade
  - Efeitos de hover e focus aprimorados
  - Container do campo com background e padding destacados
  - Label maior (1.1rem) e em cor roxa (#667eea)

- **Botão fechar no preview de e-mail**
  - Aumentado tamanho do botão (40x40px)
  - Adicionado background semitransparente branco
  - Borda arredondada (8px) para melhor aparência
  - Efeito de hover com escala e mudança de background
  - Efeito de click com animação (scale 0.95)
  - Melhor contraste e visibilidade
  - Posicionamento mantido no canto superior direito

### 📝 Detalhes Técnicos
- CSS para `input[type="file"]` e `input[type="url"]` com padding 16px
- Estilos para `#anexoFileGroup` e `#anexoUrlGroup` com background destacado
- Botão `.modal-close` com background rgba(255,255,255,0.2)
- Transições suaves e animações de hover/active
- Fix no arquivo `pkg/manual/html.go`

## [1.3.1] - 12/12/2025 23:30

### ✨ Adicionado
- **Handler customizado para inserção de imagens**
  - Novo diálogo interativo ao clicar no botão de imagem
  - Opção 1: Inserir URL da imagem (recomendado)
    - Não aumenta o tamanho do template
    - Validação automática de URL (deve começar com http:// ou https://)
    - Mensagem de sucesso ao inserir
  - Opção 2: Fazer upload de arquivo (base64)
    - Limite de 2MB por imagem
    - Aviso sobre aumento de tamanho
    - Validação de tamanho do arquivo
  - Atualização automática das estatísticas após inserir imagem
  - Mensagens claras e orientativas em cada etapa

### 📝 Detalhes Técnicos
- Função `imageHandler()` customizada para os 3 editores (header, body, footer)
- Validação de URL com regex
- Limite de upload: 2MB
- FileReader API para conversão base64
- Integração com sistema de estatísticas
- Fix no arquivo `pkg/template/html.go`

## [1.3.0] - 12/12/2025 23:15

### 🐛 Corrigido
- **Contador de tamanho não aparecia ao abrir página de novo template**
  - Adicionada chamada inicial de `updateSizeStats()` após inicialização dos editores
  - Agora mostra "0 KB" imediatamente ao abrir a página
  - Fix na função `DOMContentLoaded`

### ✨ Adicionado
- **Sistema de monitoramento de tamanho de templates em tempo real**
  - Novo painel na sidebar do editor mostrando estatísticas de tamanho
  - Indicador visual do tamanho total do HTML em KB
  - Indicador específico para tamanho de imagens base64
  - Barras de progresso com cores (verde/amarelo/vermelho) baseadas no uso
  - Contador de imagens base64 no template
  - Alertas automáticos quando próximo do limite (80%)
  - Avisos críticos quando excede o limite da Zenvia (65 KB)
  - Validação ao salvar com confirmação do usuário
  - Atualização em tempo real conforme o usuário digita

### 🎨 Interface
- Seção "📊 Tamanho do Template" na sidebar do editor
- Indicadores visuais com cores:
  - Verde: tamanho OK (< 80% do limite)
  - Amarelo: próximo do limite (80-100%)
  - Vermelho: excede o limite (> 100%)
- Mensagens contextuais:
  - ⚠️ Aviso quando imagens serão removidas automaticamente
  - ❌ Erro quando conteúdo excede limite mesmo sem imagens
- Confirmação interativa antes de salvar templates grandes

### 📝 Detalhes Técnicos
- Função `updateSizeStats()` para calcular tamanho em tempo real
- Função `countBase64Images()` para identificar e medir imagens
- Função `getByteSize()` para cálculo preciso em bytes
- Monitoramento via evento `text-change` dos editores Quill
- Validação integrada na função `saveTemplate()`
- Fix no arquivo `pkg/template/html.go`

## [1.2.2] - 12/12/2025 22:30

### 🐛 Corrigido
- **Correção CRÍTICA no envio de e-mails via Zenvia com imagens base64**
  - Adicionada validação de tamanho do HTML antes do envio
  - Implementada remoção automática de imagens base64 quando o HTML exceder 65KB
  - Mensagens com imagens base64 agora são processadas corretamente
  - Imagens removidas são substituídas por placeholder informativo
  - Limite de tamanho: 65.000 bytes (limite da API Zenvia)
  - Fix no arquivo `pkg/email/zenvia_provider.go`
  - Adicionados logs detalhados para diagnóstico de tamanho do HTML

### 📝 Detalhes Técnicos
- Constante `zenviaMaxHTMLLength = 65000` para controle do limite
- Função `removeBase64Images()` para remover imagens base64 via regex
- Validação automática no método `Send()` do ZenviaProvider
- Erro descritivo quando HTML excede limite mesmo após processamento

## [1.2.1] - 12/12/2025 21:00

### 🐛 Corrigido
- **Correção CRÍTICA no roteamento da API de templates**
  - Corrigido erro 404 ao acessar `/api/templates/:id`
  - Adicionado handler específico `handleTemplatesAPIWithID` para rotas com ID
  - Reorganizada ordem de registro de rotas (rotas mais específicas primeiro)
  - Corrigido middleware CORS para permitir PUT e DELETE
  - Fix no arquivo `pkg/dashboard/dashboard.go`
- **Correção no carregamento de templates na página de edição**
  - Adicionados logs de debug para diagnóstico
  - Implementada limpeza dos editores antes de carregar conteúdo
  - Melhorado tratamento de erros no carregamento
  - Fix no arquivo `pkg/template/html.go:987-996`
- **Correção CRÍTICA no preview de templates (disparo manual)**
  - Corrigida inversão de assunto e corpo no retorno de `ProcessTemplate()`
  - Função agora retorna corretamente: (assunto, corpo, error)
  - Preview agora exibe assunto e corpo nas posições corretas
  - Fix no arquivo `pkg/template/macro.go:198`
  - Adicionados logs de debug detalhados na função de preview

## [1.2.0] - 12/12/2025 19:30

### ✨ Adicionado
- **Sistema completo de gerenciamento de templates de e-mail HTML**
- Tabela `TEMPLATEEMAIL` no banco de dados com suporte a seções (header, body, footer)
- Campo `TEMPLATE_ID` na tabela `MENSAGEMEMAIL` para vincular e-mails aos templates
- **Interface web de listagem de templates** com:
  - Tabela paginada com busca em tempo real
  - Filtros e ordenação
  - Ações: Editar, Duplicar, Excluir
  - Design responsivo com gradiente moderno
- **Editor WYSIWYG completo** com Quill.js:
  - 3 editores separados (Header, Body, Footer)
  - Toolbar rica com formatação completa
  - Sistema de tabs para alternar seções
  - Inserção de macros via clique
  - Preview em tempo real em nova janela
  - Validação de formulário
- **Sistema de macros/placeholders** com 9 macros disponíveis:
  - `{{nome}}` - Nome do cliente
  - `{{email}}` - E-mail do cliente
  - `{{cpf_cnpj}}` - CPF/CNPJ do cliente
  - `{{codigo}}` - Código do cliente
  - `{{data}}` - Data atual (DD/MM/YYYY)
  - `{{hora}}` - Hora atual (HH:MM)
  - `{{data_hora}}` - Data e hora (DD/MM/YYYY HH:MM)
  - `{{empresa}}` - Nome da empresa
  - `{{ano}}` - Ano atual
- **REST API completa para templates** com 10 endpoints:
  - `GET /api/templates` - Listar (paginado)
  - `GET /api/templates/:id` - Buscar por ID
  - `POST /api/templates` - Criar
  - `PUT /api/templates/:id` - Atualizar
  - `DELETE /api/templates/:id` - Excluir (soft delete)
  - `GET /api/templates/macros` - Listar macros
  - `POST /api/templates/preview` - Preview com dados de exemplo
  - `POST /api/templates/:id/duplicate` - Duplicar template
- Botão "📝 Templates" no dashboard principal para acesso rápido
- Substituição automática de macros usando dados do cliente
- Validação de macros inválidas

### 🗄️ Banco de Dados
- Criada tabela `TEMPLATEEMAIL` com campos:
  - ID, NOME (único), DESCRICAO, HEADER_HTML, BODY_HTML, FOOTER_HTML
  - ASSUNTO_PADRAO, ATIVO, DATA_CRIACAO, DATA_ATUALIZACAO, CRIADO_POR
- Adicionado campo `TEMPLATE_ID` em `MENSAGEMEMAIL`
- Foreign key constraint entre MENSAGEMEMAIL e TEMPLATEEMAIL
- Índices para otimização de consultas
- Scripts SQL: `sql/create_table_templateemail.sql`

### 📦 Arquivos Criados
- `pkg/template/model.go` - Estruturas de dados e DTOs
- `pkg/template/errors.go` - Erros do domínio
- `pkg/template/repository.go` - CRUD completo (270 linhas)
- `pkg/template/macro.go` - Sistema de substituição de macros (150 linhas)
- `pkg/template/handler.go` - REST API handlers (600 linhas)
- `pkg/template/html.go` - Interface web (1110 linhas)
- `sql/create_table_templateemail.sql` - Script de criação
- `sql/alter_mensagememail_template.sql` - Alteração da tabela existente

### 📝 Mudanças
- Versão atualizada para **1.2.0**
- `pkg/message/email.go` - Adicionado campo `TemplateID`
- `pkg/message/repository.go` - Queries atualizadas com TEMPLATE_ID
- `pkg/dashboard/dashboard.go` - Interface e rotas para templates
- `cmd/icrmsenderemail/main.go` - Registro do módulo de templates
- Dashboard principal agora tem botão laranja para "Templates"

### 🎯 Funcionalidades do Template
- **CRUD completo**: Criar, editar, listar, duplicar, excluir templates
- **Soft delete**: Templates excluídos ficam inativos mas não são removidos
- **Paginação**: Lista de templates com paginação e busca
- **Validação**: Nome único, corpo obrigatório
- **Preview**: Visualização com dados de exemplo antes de salvar
- **Macros**: Substituição automática com dados do cliente
- **Seções**: Header, Body e Footer editáveis separadamente
- **Versionamento**: Data de criação e última atualização

### 🐛 Bugs Conhecidos
Nenhum no momento.

### ⚠️ Notas de Migração

#### De 1.1.0 para 1.2.0

1. **Banco de Dados:**
   ```bash
   # Execute os scripts SQL
   sqlplus usuario/senha@tns @sql/create_table_templateemail.sql
   sqlplus usuario/senha@tns @sql/alter_mensagememail_template.sql
   ```

2. **Acesso às Funcionalidades:**
   - Templates: `http://localhost:3101/templates`
   - API: `http://localhost:3101/api/templates`
   - Dashboard atualizado com botão "📝 Templates"

3. **Compatibilidade:**
   - Totalmente compatível com versão anterior
   - Módulo de templates é opcional (não quebra funcionalidade existente)

---

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
