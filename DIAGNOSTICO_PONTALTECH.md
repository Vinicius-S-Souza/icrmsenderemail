# Diagnóstico do Erro Pontaltech

## Erro Identificado
```
"Post \"https://api.pontaltech.com.br/v1/email/send\": dial tcp: lookup api.pontaltech.com.br on 10.255.255.254:53: no such host"
```

## Causa Raiz
O domínio `api.pontaltech.com.br` **NÃO EXISTE** (NXDOMAIN).

### Testes DNS Realizados
```bash
$ nslookup api.pontaltech.com.br
Server:         10.255.255.254
Address:        10.255.255.254#53

** server can't find api.pontaltech.com.br: NXDOMAIN

$ nslookup pontaltech.com.br
Server:         10.255.255.254
Address:        10.255.255.254#53

Name:   pontaltech.com.br
Address: 50.62.140.1  ✅ EXISTE
```

## Possíveis Causas

1. **URL da API Incorreta** - O subdomínio `api.pontaltech.com.br` pode não ser o correto
2. **Serviço Descontinuado** - A API Pontaltech pode ter sido descontinuada ou mudou de URL
3. **Documentação Desatualizada** - A URL pode ter sido copiada de documentação antiga
4. **API em Outro Domínio** - A API pode estar em outro domínio (ex: `smtp.pontaltech.com.br`, `api.pontaltech.net`, etc)

## Soluções Implementadas

### 1. URL Configurável
Adicionado suporte para configurar a URL da API via arquivo `dbinit.ini`:

```ini
[email]
# ===== Pontaltech (provider=pontaltech) =====
pontaltech_username=intellisyspremium
pontaltech_password=kCHykmXl
pontaltech_account_id=7818
# Configure a URL correta da API Pontaltech aqui:
pontaltech_api_url=https://api.exemplo.com/v1/email/send
```

### 2. Melhor Detecção de Erros DNS
O código agora detecta especificamente erros de DNS e fornece mensagens claras:

```
❌ Erro de DNS ao acessar API Pontaltech
URL: https://api.pontaltech.com.br/v1/email/send
Solução: Verifique se a URL da API está correta. Configure 'pontaltech_api_url' no dbinit.ini
```

### 3. Aviso na Inicialização
Quando a URL customizada não é configurada, um aviso é exibido:

```
⚠️  URL da API Pontaltech não configurada, usando padrão
URL: https://api.pontaltech.com.br/v1/email/send
ATENÇÃO: O domínio api.pontaltech.com.br pode não existir. Configure 'pontaltech_api_url' no dbinit.ini se necessário
```

## Próximos Passos

### Opção 1: Usar Outro Provider
Se o serviço Pontaltech não estiver disponível, considere usar outros providers:
- **SendGrid** (já configurado e funcionando)
- **SMTP** (genérico, funciona com qualquer servidor SMTP)
- **Zenvia** (se disponível)
- **Mock** (para testes)

Para alterar, edite o `dbinit.ini`:
```ini
[email]
provider=sendgrid  # ou smtp, zenvia, mock
```

### Opção 2: Descobrir a URL Correta
1. **Contate o suporte Pontaltech** para obter a URL correta da API
2. **Verifique a documentação oficial** do Pontaltech
3. **Verifique contratos ou emails** da Pontaltech que possam conter a URL

### Opção 3: Usar Mock para Testes
Temporariamente, use o provider Mock para continuar os testes:
```ini
[email]
provider=mock
```

## Resumo Técnico das Alterações

### Arquivos Modificados:

1. **pkg/config/config.go**
   - Adicionado campo `PontaltechAPIURL` na struct `EmailConfig`
   - Carrega `pontaltech_api_url` do arquivo INI

2. **pkg/email/pontaltech_provider.go**
   - Adicionado campo `apiURL` na struct `PontaltechProvider`
   - Modificado `NewPontaltechProvider` para aceitar URL customizada
   - Melhorado tratamento de erros DNS
   - Adicionado aviso quando URL não está configurada

3. **cmd/icrmsenderemail/main.go**
   - Atualizado chamada de `NewPontaltechProvider` para passar URL customizada

4. **dbinit.ini e dbinit.ini.example**
   - Adicionado comentários sobre `pontaltech_api_url`
   - Documentação do problema de DNS

## Como Testar

```bash
# 1. Configure a URL correta no dbinit.ini (se disponível)
nano dbinit.ini

# 2. Ou mude para outro provider
# Edite [email] provider=sendgrid

# 3. Execute a aplicação
go run ./cmd/icrmsenderemail/main.go
```
