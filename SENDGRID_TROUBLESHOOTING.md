# SendGrid Troubleshooting - Email não chega ao destinatário
**Data de criação:** 11/12/2025
**Versão:** 1.0.0

## Problema
O SendGrid retorna status 202 (aceito) mas o email não chega ao destinatário.

## Principais Causas

### 1. ⚠️ Email/Domínio "From" Não Verificado
**Mais comum!** O SendGrid aceita o email mas só envia se o remetente estiver verificado.

**Solução:**
1. Acesse: https://app.sendgrid.com/settings/sender_auth/senders
2. Verifique se o email "from" está na lista de Single Sender Verification
3. Se não estiver, clique em "Create New Sender" e siga o processo de verificação
4. Você receberá um email de confirmação - clique no link para verificar

**Ou configure Domain Authentication (recomendado para produção):**
1. Acesse: https://app.sendgrid.com/settings/sender_auth/domain/create
2. Adicione seus registros DNS (CNAME)
3. Aguarde a verificação (pode levar até 48h)

### 2. 🔒 Conta em Sandbox Mode
Se sua conta SendGrid está em modo sandbox, emails só são enviados para endereços pré-aprovados.

**Solução:**
1. Acesse: https://app.sendgrid.com/settings/mail_settings
2. Verifique se "Sandbox Mode" está desabilitado
3. Se estiver habilitado, desabilite ou adicione o destinatário à lista de emails aprovados

### 3. 📊 Verificar Activity Feed
O SendGrid mantém um log detalhado de todas as tentativas de envio.

**Como verificar:**
1. Acesse: https://app.sendgrid.com/email_activity
2. Procure pelo Message ID que aparece nos logs do icrmsenderemail
3. Verifique o status:
   - **Processed**: Aceito pelo SendGrid
   - **Dropped**: Bloqueado (veja o motivo)
   - **Delivered**: Entregue com sucesso
   - **Bounce**: Rejeitado pelo servidor de destino
   - **Deferred**: Tentativa temporária de reenvio

### 4. 🚫 Lista de Supressão (Suppressions)
O destinatário pode estar em uma lista de bloqueio.

**Solução:**
1. Acesse: https://app.sendgrid.com/suppressions
2. Verifique as abas:
   - **Bounces**: Emails que retornaram erro permanente
   - **Blocks**: Bloqueados por IP ou outros motivos
   - **Spam Reports**: Marcados como spam
   - **Invalid Emails**: Endereços inválidos
   - **Unsubscribes**: Emails que cancelaram inscrição
3. Remova o destinatário da lista se necessário

### 5. 📧 Validar o JSON Enviado
Com a nova versão, o JSON completo está sendo logado.

**Como verificar:**
1. Olhe os logs: `log/icrmsenderemail_YYYYMMDD.log`
2. Procure por: `📤 JSON enviado para SendGrid`
3. Compare com o JSON que funciona no WinDev

**Estrutura esperada:**
```json
{
  "personalizations": [
    {
      "to": [
        {"email": "destinatario@exemplo.com"}
      ]
    }
  ],
  "from": {"email": "remetente@exemplo.com"},
  "subject": "Assunto do Email",
  "content": [
    {
      "type": "text/plain",
      "value": "Corpo do email"
    }
  ]
}
```

### 6. 🔑 Verificar API Key
Certifique-se de que a API Key tem permissões corretas.

**Solução:**
1. Acesse: https://app.sendgrid.com/settings/api_keys
2. Verifique se a chave tem permissão "Mail Send" ativada
3. Se necessário, crie uma nova API Key com permissões corretas

### 7. 📈 Limites de Envio
Contas gratuitas têm limite de 100 emails/dia.

**Solução:**
1. Acesse: https://app.sendgrid.com/account/billing
2. Verifique seu plano e limite de envios
3. Upgrade se necessário

## Como Comparar com WinDev

1. **Capture o JSON do WinDev:**
   - No código WinDev, adicione um log antes do `HTTPSend`:
   ```windev
   Info(jConteudo..JSONFormat())
   ```

2. **Compare com o JSON do Go:**
   - Veja o log: `📤 JSON enviado para SendGrid`

3. **Diferenças comuns:**
   - Formato de `content_type`: deve ser `text/plain` ou `text/html`
   - Anexos: verificar se base64 está correto
   - Estrutura de arrays: Go usa índice 0, WinDev usa índice 1

## Checklist de Verificação

- [ ] Email/domínio "from" está verificado no SendGrid
- [ ] Conta não está em Sandbox Mode
- [ ] Destinatário não está em lista de supressão
- [ ] API Key tem permissão "Mail Send"
- [ ] Não excedeu limite de envios do plano
- [ ] JSON enviado está correto (comparar com WinDev)
- [ ] Activity Feed mostra "Delivered"

## Logs Úteis

Com a versão atualizada, você verá:

```
📧 Enviando email via SendGrid
📤 JSON enviado para SendGrid
📩 Resposta da API SendGrid
✅ Email aceito pelo SendGrid (status 202)
ℹ️  Para rastrear entrega, acesse SendGrid Activity Feed
```

## Comandos Úteis

```bash
# Ver últimos envios SendGrid
tail -50 log/icrmsenderemail_$(date +%Y%m%d).log | grep SendGrid

# Ver JSON enviado
tail -200 log/icrmsenderemail_$(date +%Y%m%d).log | grep "JSON enviado"

# Ver Message IDs
tail -100 log/icrmsenderemail_$(date +%Y%m%d).log | grep "message_id"
```

## Links Úteis

- Activity Feed: https://app.sendgrid.com/email_activity
- Sender Authentication: https://app.sendgrid.com/settings/sender_auth/senders
- Suppressions: https://app.sendgrid.com/suppressions
- API Keys: https://app.sendgrid.com/settings/api_keys
- Mail Settings: https://app.sendgrid.com/settings/mail_settings
- Documentação API: https://docs.sendgrid.com/api-reference/mail-send/mail-send

## Contato Suporte SendGrid

Se após todas as verificações o problema persistir:
- Suporte: https://support.sendgrid.com/
- Status: https://status.sendgrid.com/
