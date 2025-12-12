# ICRMSenderEmail

**Data de criação:** 11/12/2025
**Última atualização:** 12/12/2025 23:45
**Versão:** 1.3.2

Serviço em Golang para envio automatizado de e-mails através de múltiplos provedores (SMTP, SendGrid, Zenvia, Pontaltech), com suporte a dashboard web e disparo manual.

## 📋 Características

- ✅ Envio de e-mail via múltiplos provedores
- ✅ **Suporte a anexos** (base64 e URL pública)
- ✅ **Detecção automática de tipo de anexo por provider**
- ✅ **Sistema completo de templates HTML** com editor WYSIWYG
- ✅ **Macros/placeholders** para personalização de e-mails
- ✅ **Preview de templates** em tempo real
- ✅ **Monitoramento de tamanho de templates** em tempo real
- ✅ **Validação automática de limites da API Zenvia** (65KB)
- ✅ **Remoção automática de imagens base64** quando exceder limite
- ✅ **Handler customizado para inserção de imagens** (URL ou upload)
- ✅ Processamento paralelo com workers
- ✅ Retry automático com exponential backoff
- ✅ Circuit breaker para proteção contra falhas
- ✅ Dashboard web em tempo real (Server-Sent Events)
- ✅ **Interface web inteligente para disparo manual** (adapta-se ao provider)
- ✅ **UI moderna e responsiva** com feedback visual
- ✅ Health check HTTP
- ✅ Logs estruturados com rotação diária
- ✅ Métricas de performance
- ✅ Graceful shutdown
- ✅ Suporte a Windows Service
- ✅ Validação de clientes via CLIENTES + CLIENTESEXTENSAO
- ✅ Suporte a HTML e texto plano
- ✅ **E-mail "from" configurável via default_from**

## 🚀 Provedores Suportados

| Provider | Descrição | Autenticação | Anexos |
|----------|-----------|--------------|--------|
| `mock` | Simulação para testes | Nenhuma | ❌ Não |
| `smtp` | SMTP genérico | Usuário/Senha | ❌ Não |
| `sendgrid` | SendGrid API v3 | API Key | ✅ Base64 |
| `zenvia` | Zenvia Email API | Token | ✅ URL Pública |
| `pontaltech` | Pontaltech Email API | Basic Auth | ✅ Base64 |

### 📎 Suporte a Anexos

#### SendGrid e Pontaltech
- ✅ Anexos enviados em **base64** diretamente no JSON
- ✅ Upload de arquivo pela interface web
- ✅ Tamanho máximo: 10MB
- ✅ Todos os tipos de arquivo suportados

#### Zenvia
- ⚠️ Anexos **apenas via URL pública**
- ✅ Campo de URL na interface web
- ✅ Validação automática de URL
- ❌ **NÃO aceita** base64
- 📝 Ver documentação: [ZENVIA_ANEXOS.md](ZENVIA_ANEXOS.md)

## 📦 Instalação

### Pré-requisitos

- Go 1.23.0+
- Oracle Database 11g+ (com driver godror)
- Acesso à tabela `MENSAGEMEMAIL` (ver SQL abaixo)

### Build

```bash
# Clone o repositório
cd /caminho/para/icrmsenderemail

# Download de dependências
go mod download

# Compilar
go build -o build/icrmsenderemail.exe ./cmd/icrmsenderemail
```

Ou use o Makefile:

```bash
make build
```

## ⚙️ Configuração

1. Copie o arquivo de configuração exemplo:

```bash
cp dbinit.ini.example dbinit.ini
```

2. Edite `dbinit.ini` com suas credenciais:

```ini
[oracle]
username=seu_usuario
password=sua_senha
tns=seu_tns

[email]
provider=sendgrid
sendgrid_api_key=SG.xxxxxxxxxxxxx
default_from=noreply@suaempresa.com
```

## 🗄️ Banco de Dados

Execute o script SQL para criar a tabela:

```bash
sqlplus usuario/senha@tns @sql/create_table_mensagememail.sql
```

A tabela `MENSAGEMEMAIL` contém:

- **ID**: Identificador único (NUMBER)
- **CLICODIGO**: Código do cliente (FK para CLIENTES)
- **REMETENTE**: E-mail do remetente (VARCHAR2)
- **DESTINATARIO**: E-mail do destinatário (VARCHAR2)
- **ASSUNTO**: Assunto do e-mail (VARCHAR2)
- **CORPO**: Corpo do e-mail (CLOB)
- **TIPO_CORPO**: Tipo de conteúdo: `text/plain` ou `text/html`
- **STATUS_ENVIO**: Status (0=Pendente, 2=Enviado, 3=Erro, 4=Falha, 125=Inválido)
- **DATA_CADASTRO**, **DATA_AGENDAMENTO**, **DATA_ENVIO**: Timestamps
- **QTD_TENTATIVAS**: Contador de tentativas
- **DETALHES_ERRO**: Mensagem de erro
- **ID_PROVIDER**: ID retornado pelo provider
- **METODO_ENVIO**: Código numérico do provider
- **PRIORIDADE**: Prioridade (1=Alta, 2=Normal, 3=Baixa)
- **ANEXO_REFERENCIA**, **ANEXO_NOME**, **ANEXO_TIPO**: Campos de anexo
- **IP_ORIGEM**: IP de origem (disparo manual)

## 🎯 Uso

### Modo Normal (Foreground)

```bash
./build/icrmsenderemail.exe
```

### Como Serviço do Windows

```bash
# Instalar serviço
./build/icrmsenderemail.exe install

# Iniciar serviço
./build/icrmsenderemail.exe start

# Parar serviço
./build/icrmsenderemail.exe stop

# Reiniciar serviço
./build/icrmsenderemail.exe restart

# Desinstalar serviço
./build/icrmsenderemail.exe uninstall
```

### Ver Versão

```bash
./build/icrmsenderemail.exe version
```

## 📊 Dashboard Web

Acesse o dashboard em tempo real:

```
http://localhost:3101
```

O dashboard mostra:

- Total de e-mails processados
- E-mails pendentes
- Taxa de sucesso/erro
- E-mails inválidos
- Tempos médios de processamento
- Gráficos em tempo real (atualização a cada 2 segundos)

## 📝 Templates de E-mail

Acesse o gerenciador de templates:

```
http://localhost:3101/templates
```

### Funcionalidades:

- **CRUD Completo**: Criar, editar, listar, duplicar e excluir templates
- **Editor WYSIWYG**: Interface rica com Quill.js para edição HTML
- **Seções Separadas**: Header, Body e Footer editáveis individualmente
- **Macros Disponíveis**:
  - `{{nome}}` - Nome do cliente
  - `{{email}}` - E-mail do cliente
  - `{{cpf_cnpj}}` - CPF/CNPJ do cliente
  - `{{codigo}}` - Código do cliente
  - `{{data}}` - Data atual (DD/MM/YYYY)
  - `{{hora}}` - Hora atual (HH:MM)
  - `{{data_hora}}` - Data e hora completa
  - `{{empresa}}` - Nome da empresa
  - `{{ano}}` - Ano atual
- **Preview em Tempo Real**: Visualize como o e-mail ficará antes de salvar
- **Busca e Paginação**: Encontre templates facilmente
- **Soft Delete**: Templates excluídos ficam inativos mas não são removidos

### API REST de Templates:

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/api/templates` | Listar templates (paginado) |
| GET | `/api/templates/:id` | Buscar por ID |
| POST | `/api/templates` | Criar novo template |
| PUT | `/api/templates/:id` | Atualizar template |
| DELETE | `/api/templates/:id` | Excluir (soft delete) |
| GET | `/api/templates/macros` | Listar macros disponíveis |
| POST | `/api/templates/preview` | Preview com dados de exemplo |
| POST | `/api/templates/:id/duplicate` | Duplicar template |

## 📨 Disparo Manual

Acesse a interface de disparo manual:

```
http://localhost:3101/manual
```

Funcionalidades:

1. **Validar Cliente**: Por código ou CPF/CNPJ
   - Busca em `CLIENTES` + `CLIENTESEXTENSAO`
   - Retorna e-mail de `CLIEXTEMAIL2`
2. **Compor E-mail**: Destinatário, assunto, corpo
   - Suporte a texto plano ou HTML
   - Futuramente: Seleção de template com macros
3. **Acompanhamento**: Status em tempo real do envio

## 🔍 Health Check

Endpoint de saúde:

```
GET http://localhost:8081/health
```

Resposta:

```json
{
  "status": "ok",
  "timestamp": "2025-12-11T10:30:00Z",
  "database": "connected"
}
```

## 📈 Métricas

As métricas incluem:

- Total de mensagens processadas
- Taxa de sucesso/erro
- E-mails inválidos
- Tempo médio de processamento
- Tempo médio de envio
- Tempo médio de query
- Queries executadas

Logs a cada 60 segundos e no shutdown.

## 🏗️ Arquitetura

```
icrmsenderemail/
├── cmd/icrmsenderemail/
│   └── main.go                   # Ponto de entrada
├── pkg/
│   ├── config/                   # Configurações INI
│   ├── database/                 # Conexão Oracle
│   ├── logger/                   # Logger com rotação
│   ├── message/                  # Email model + repository + processor
│   ├── email/                    # Providers (SMTP, SendGrid, etc.)
│   ├── cliente/                  # Repository de clientes
│   ├── dashboard/                # Dashboard web
│   ├── manual/                   # Handler de disparo manual
│   ├── template/                 # Sistema de templates ⭐ NOVO
│   │   ├── model.go             # Estruturas de dados
│   │   ├── repository.go        # CRUD de templates
│   │   ├── handler.go           # REST API
│   │   ├── html.go              # Interface WYSIWYG
│   │   ├── macro.go             # Sistema de macros
│   │   └── errors.go            # Erros do domínio
│   ├── retry/                    # Retry com backoff
│   ├── control/                  # Graceful shutdown
│   ├── health/                   # Health check
│   ├── metrics/                  # Métricas de performance
│   ├── service/                  # Windows Service wrapper
│   └── version/                  # Informações de versão
├── sql/                          # Scripts SQL
│   ├── create_table_mensagememail.sql
│   ├── create_table_templateemail.sql ⭐ NOVO
│   └── alter_mensagememail_template.sql ⭐ NOVO
├── log/                          # Logs (criado automaticamente)
├── build/                        # Binários compilados
├── dbinit.ini.example            # Exemplo de configuração
├── go.mod                        # Dependências
├── Makefile                      # Build commands
└── README.md                     # Este arquivo
```

## 🔧 Desenvolvimento

### Comandos Makefile

```bash
make build          # Compilar aplicação
make run            # Executar em modo desenvolvimento
make clean          # Limpar build
make test           # Executar testes
```

### Provider Pattern

Criar um novo provider:

```go
type MyProvider struct {
    apiKey string
    logger *zap.Logger
}

func (p *MyProvider) Send(ctx context.Context, email EmailData) (SendResult, error) {
    // Implementar lógica de envio
}

func (p *MyProvider) GetName() string {
    return "MyProvider"
}

func (p *MyProvider) ValidateEmail(email string) error {
    return message.ValidateEmail(email)
}
```

Registrar em `main.go`:

```go
case "myprovider":
    provider = email.NewMyProvider(cfg.Email.MyProviderAPIKey, log)
```

## 🛡️ Segurança

- ✅ Validação de formato de e-mail
- ✅ Proteção contra SQL injection (prepared statements)
- ✅ Circuit breaker para proteção contra falhas
- ✅ Timeout em todas as operações de I/O
- ✅ Graceful shutdown para evitar perda de dados
- ✅ Logs estruturados (não expõem dados sensíveis)

## 📝 Logs

Os logs são gravados em:

```
log/icrmsenderemail_YYYYMMDD.log
```

Formato JSON estruturado:

```json
{
  "level": "info",
  "ts": "2025-12-11T10:30:00.123Z",
  "caller": "message/processor.go:289",
  "msg": "Email enviado com sucesso",
  "email_id": 12345,
  "provider_id": "abc123",
  "provider": "SendGrid",
  "duracao": "150ms"
}
```

## 🚦 Status de Envio

| Código | Descrição |
|--------|-----------|
| 0 | Pendente |
| 2 | Enviado com sucesso |
| 3 | Erro temporário (vai retentar) |
| 4 | Falha permanente |
| 125 | E-mail inválido |

## 🔗 Códigos de Provider

| Provider | Código |
|----------|--------|
| Mock | 0 |
| SMTP | 1024 |
| SendGrid | 2048 |
| Zenvia | 4096 |
| Pontaltech | 8192 |

## 📞 Suporte

Para dúvidas e problemas:

1. Verifique os logs em `log/icrmsenderemail_YYYYMMDD.log`
2. Consulte o dashboard em `http://localhost:3101`
3. Verifique o health check em `http://localhost:8081/health`

## 📄 Licença

Copyright © 2025 - Uso Interno

---

**Desenvolvido com** ❤️ **usando Golang 1.23.0**
