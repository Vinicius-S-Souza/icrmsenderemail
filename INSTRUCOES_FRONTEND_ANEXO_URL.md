# Instruções para Adicionar Campo de URL de Anexo na Página Manual

**Data:** 11/12/2025 16:45
**Versão:** 1.0.0

## ✅ Backend Implementado

O backend já está pronto para receber anexos por URL! As seguintes modificações foram feitas:

1. ✅ Estrutura `Attachment` agora suporta campo `URL`
2. ✅ Handler aceita parâmetro `attachmentUrl` no JSON
3. ✅ Validação de URL implementada (deve começar com http:// ou https://)
4. ✅ Processor detecta automaticamente se é URL ou base64
5. ✅ Provider Zenvia usa URL corretamente via campo `fileUrl`

## 📝 Modificações Necessárias no Frontend

Você precisa adicionar o seguinte ao arquivo `pkg/manual/html.go`:

### 1. Adicionar campo de URL de anexo (após linha 428)

```html
<div class="form-group" id="anexoUrlGroup" style="display: none;">
    <label for="anexoUrl">OU URL do Anexo (somente Zenvia)</label>
    <input type="url" id="anexoUrl" placeholder="https://exemplo.com/arquivo.pdf" onchange="handleUrlInput()">
    <div class="hint" id="anexoUrlInfo">Informe a URL pública do arquivo a ser anexado</div>
</div>
```

### 2. Adicionar variável global (após linha 443)

```javascript
let selectedAttachmentUrl = "";
let currentProvider = "";
```

### 3. Modificar função `carregarProviderInfo()` para detectar provider (linha ~795)

Adicione após `document.getElementById('provider-name').textContent = displayName;`:

```javascript
currentProvider = data.providerName;
toggleAnexoFields(currentProvider);
```

### 4. Adicionar nova função para alternar campos de anexo

```javascript
function toggleAnexoFields(providerName) {
    const fileGroup = document.getElementById('anexo').parentElement;
    const urlGroup = document.getElementById('anexoUrlGroup');

    if (providerName === 'zenvia') {
        // Zenvia: mostrar campo de URL, esconder upload de arquivo
        fileGroup.style.display = 'none';
        urlGroup.style.display = 'block';
    } else {
        // Outros providers: mostrar upload de arquivo, esconder URL
        fileGroup.style.display = 'block';
        urlGroup.style.display = 'none';
    }
}
```

### 5. Adicionar função de validação de URL

```javascript
function handleUrlInput() {
    const url = document.getElementById('anexoUrl').value.trim();
    const anexoUrlInfo = document.getElementById('anexoUrlInfo');

    if (!url) {
        selectedAttachmentUrl = "";
        anexoUrlInfo.textContent = 'Informe a URL pública do arquivo a ser anexado';
        anexoUrlInfo.style.color = '';
        return;
    }

    // Validar URL
    if (!url.startsWith('http://') && !url.startsWith('https://')) {
        anexoUrlInfo.textContent = '❌ URL inválida. Deve começar com http:// ou https://';
        anexoUrlInfo.style.color = '#f44336';
        selectedAttachmentUrl = "";
        return;
    }

    // Extrair nome do arquivo da URL
    const fileName = url.substring(url.lastIndexOf('/') + 1) || 'arquivo';

    selectedAttachmentUrl = url;
    anexoUrlInfo.textContent = '✓ URL válida: ' + fileName;
    anexoUrlInfo.style.color = '#4caf50';

    console.log('URL de anexo selecionada:', url);
}
```

### 6. Modificar função `dispararEmail()` para enviar URL (linha ~650)

Substituir o bloco de `body: JSON.stringify(...)` por:

```javascript
body: JSON.stringify({
    cliCodigo: clienteValidado.cliCodigo,
    email: emailDestinatario,
    assunto: assunto,
    mensagem: mensagem,
    isHtml: isHtml,
    attachmentData: selectedAttachment?.data || "",
    attachmentName: selectedAttachment?.name || "",
    attachmentType: selectedAttachment?.type || "",
    attachmentUrl: selectedAttachmentUrl || ""  // ← NOVO CAMPO
})
```

### 7. Limpar URL após envio (após linha 676)

Adicionar:

```javascript
// Limpar URL de anexo
document.getElementById('anexoUrl').value = '';
selectedAttachmentUrl = "";
const anexoUrlInfo = document.getElementById('anexoUrlInfo');
anexoUrlInfo.textContent = 'Informe a URL pública do arquivo a ser anexado';
anexoUrlInfo.style.color = '';
```

## 🎯 Resultado Esperado

### Para Zenvia:
- Campo "URL do Anexo" será exibido
- Campo de upload de arquivo será ocultado
- Usuário digita URL pública do arquivo
- URL é validada e enviada ao backend

### Para SendGrid/Pontaltech:
- Campo de upload de arquivo será exibido
- Campo "URL do Anexo" será ocultado
- Funciona como antes (base64)

## 🔍 Testando

1. Inicie o serviço com Zenvia configurado
2. Abra a página de disparo manual
3. Verifique que o campo "URL do Anexo" aparece
4. Digite uma URL válida (ex: https://exemplo.com/teste.pdf)
5. Envie o e-mail
6. Verifique nos logs que a URL foi recebida corretamente

## ⚠️ Importante

- A URL **DEVE** ser pública e acessível pela internet
- Zenvia **NÃO aceita** anexos em base64
- Se usar arquivo em vez de URL com Zenvia, o anexo será ignorado
- SendGrid e Pontaltech **NÃO precisam** de URL, usam base64

## 📊 Status

- ✅ Backend: 100% implementado
- ⏳ Frontend: Aguardando implementação das modificações acima
