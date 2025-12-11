# Limitação de Anexos - Zenvia

**Data:** 11/12/2025 16:15
**Versão:** 1.0.0

## ⚠️ IMPORTANTE: Zenvia não suporta anexos em base64

A API Zenvia Email **NÃO aceita** anexos enviados diretamente em base64 como SendGrid e Pontaltech.

## Como funciona a Zenvia:

A Zenvia **exige** que anexos estejam hospedados em URLs públicas na internet.

### Formato correto (Zenvia):
```json
{
  "from": "remetente@dominio.com",
  "to": "destinatario@exemplo.com",
  "contents": [
    {
      "type": "email",
      "subject": "Assunto",
      "html": "Corpo do email",
      "attachments": [
        {
          "fileUrl": "https://seuservidor.com/arquivos/documento.pdf",
          "fileName": "documento.pdf"
        }
      ]
    }
  ]
}
```

## Status atual:

**Anexos estão DESABILITADOS para o provider Zenvia** até que seja implementada uma solução de hospedagem de arquivos.

## Possíveis soluções futuras:

1. **Servidor de arquivos temporários**
   - Implementar endpoint para upload de arquivos
   - Hospedar arquivos temporariamente (ex: 24 horas)
   - Gerar URL pública para cada arquivo
   - Usar essa URL no campo `fileUrl`

2. **Integração com serviços de armazenamento**
   - AWS S3
   - Google Cloud Storage
   - Azure Blob Storage
   - Gerar URLs assinadas temporárias

3. **Servidor HTTP local público**
   - Expor pasta local via HTTP
   - Usar ngrok ou similar para URL pública
   - **Não recomendado para produção**

## Comparação com outros providers:

| Provider | Suporte a Base64 | Método |
|----------|------------------|---------|
| SendGrid | ✅ Sim | Base64 direto no JSON |
| Pontaltech | ✅ Sim | Base64 direto no JSON |
| Zenvia | ❌ Não | Apenas URL pública (fileUrl) |

## Comportamento atual:

Se tentar enviar um e-mail com anexo via Zenvia:
- ⚠️ Um aviso será logado
- 📧 O e-mail será enviado **SEM o anexo**
- ✅ O envio não falhará, apenas ignorará o anexo

## Exemplo de log:

```
⚠️  AVISO: Zenvia não suporta anexos em base64
filename: documento.pdf
info: Zenvia só aceita anexos via URL pública (fileUrl). O anexo será ignorado.
```

## Recomendação:

Para envio de e-mails com anexo, use:
- **SendGrid** (recomendado) ✅
- **Pontaltech** ✅

Não use Zenvia se precisar de anexos.
