package manual

// manualSendHTML contém o HTML da página de disparo manual de e-mail
// Data de criação: 11/12/2025
const manualSendHTML = `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ICRMSenderEmail - Disparo Manual de E-mail</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }

        .container {
            max-width: 900px;
            margin: 0 auto;
        }

        header {
            background: white;
            padding: 30px;
            border-radius: 15px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.1);
            margin-bottom: 30px;
            display: flex;
            justify-content: space-between;
            align-items: center;
            flex-wrap: wrap;
        }

        .header-left {
            flex: 1;
        }

        header h1 {
            color: #667eea;
            font-size: 2.5em;
            margin-bottom: 10px;
        }

        .subtitle {
            color: #666;
            font-size: 1.1em;
        }

        .subtitle strong {
            color: #667eea;
        }

        .status {
            display: inline-block;
            padding: 5px 15px;
            background: #10b981;
            color: white;
            border-radius: 20px;
            font-size: 0.9em;
            margin-left: 15px;
            animation: pulse 2s infinite;
        }

        @keyframes pulse {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.7; }
        }

        .dashboard-btn {
            display: inline-block;
            padding: 12px 24px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            text-decoration: none;
            border-radius: 8px;
            font-weight: 600;
            transition: all 0.3s;
            box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
        }

        .dashboard-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 16px rgba(102, 126, 234, 0.4);
        }

        .card {
            background: white;
            border-radius: 12px;
            padding: 30px;
            box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
            margin-bottom: 20px;
        }

        .form-group {
            margin-bottom: 20px;
        }

        label {
            display: block;
            font-weight: 600;
            color: #333;
            margin-bottom: 8px;
            font-size: 0.95rem;
        }

        input[type="text"],
        input[type="email"],
        textarea {
            width: 100%;
            padding: 12px;
            border: 2px solid #e0e0e0;
            border-radius: 8px;
            font-size: 1rem;
            transition: border-color 0.3s;
        }

        input[type="text"]:focus,
        input[type="email"]:focus,
        textarea:focus {
            outline: none;
            border-color: #667eea;
        }

        textarea {
            resize: vertical;
            min-height: 150px;
            font-family: inherit;
        }

        .checkbox-group {
            display: flex;
            align-items: center;
            gap: 8px;
        }

        input[type="checkbox"] {
            width: 18px;
            height: 18px;
            cursor: pointer;
        }

        .button-group {
            display: flex;
            gap: 10px;
            margin-top: 20px;
        }

        button {
            flex: 1;
            padding: 12px 24px;
            border: none;
            border-radius: 8px;
            font-size: 1rem;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s;
        }

        .btn-primary {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
        }

        .btn-primary:hover:not(:disabled) {
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
        }

        .btn-secondary {
            background: #f5f5f5;
            color: #333;
        }

        .btn-secondary:hover:not(:disabled) {
            background: #e0e0e0;
        }

        button:disabled {
            opacity: 0.6;
            cursor: not-allowed;
        }

        .alert {
            padding: 15px;
            border-radius: 8px;
            margin-bottom: 20px;
            display: none;
        }

        .alert.show {
            display: block;
            animation: slideDown 0.3s ease;
        }

        @keyframes slideDown {
            from {
                opacity: 0;
                transform: translateY(-10px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }

        .alert-success {
            background: #e8f5e9;
            color: #2e7d32;
            border-left: 4px solid #4caf50;
        }

        .alert-error {
            background: #ffebee;
            color: #c62828;
            border-left: 4px solid #f44336;
        }

        .alert-info {
            background: #e3f2fd;
            color: #1565c0;
            border-left: 4px solid #2196f3;
        }

        .client-info {
            background: #f5f5f5;
            padding: 20px;
            border-radius: 8px;
            margin-top: 20px;
            display: none;
        }

        .client-info.show {
            display: block;
            animation: fadeIn 0.3s ease;
        }

        @keyframes fadeIn {
            from { opacity: 0; }
            to { opacity: 1; }
        }

        .client-info h3 {
            color: #667eea;
            margin-bottom: 15px;
            font-size: 1.1rem;
        }

        .info-row {
            display: flex;
            padding: 8px 0;
            border-bottom: 1px solid #e0e0e0;
        }

        .info-row:last-child {
            border-bottom: none;
        }

        .info-label {
            font-weight: 600;
            color: #666;
            min-width: 150px;
        }

        .info-value {
            color: #333;
            flex: 1;
        }

        .status-tracking {
            background: #f5f5f5;
            padding: 20px;
            border-radius: 8px;
            margin-top: 20px;
            display: none;
        }

        .status-tracking.show {
            display: block;
            animation: fadeIn 0.3s ease;
        }

        .status-tracking h3 {
            color: #667eea;
            margin-bottom: 15px;
            font-size: 1.1rem;
        }

        .status-badge {
            display: inline-block;
            padding: 6px 12px;
            border-radius: 20px;
            font-size: 0.85rem;
            font-weight: 600;
            margin-top: 5px;
        }

        .status-pendente {
            background: #fff3e0;
            color: #e65100;
        }

        .status-enviado {
            background: #e8f5e9;
            color: #2e7d32;
        }

        .status-erro {
            background: #ffebee;
            color: #c62828;
        }

        .status-invalido {
            background: #fce4ec;
            color: #880e4f;
        }

        .spinner {
            border: 3px solid #f3f3f3;
            border-top: 3px solid #667eea;
            border-radius: 50%;
            width: 20px;
            height: 20px;
            animation: spin 1s linear infinite;
            display: inline-block;
            margin-left: 10px;
            vertical-align: middle;
        }

        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }

        .hint {
            font-size: 0.85rem;
            color: #666;
            margin-top: 5px;
            font-style: italic;
        }

        @media (max-width: 768px) {
            .button-group {
                flex-direction: column;
            }

            header h1 {
                font-size: 1.5rem;
            }

            .card {
                padding: 20px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <div class="header-left">
                <h1>📧 ICRMSenderEmail</h1>
                <span class="subtitle">Disparo Manual de E-mail</span>
                <span class="subtitle"> | Provedor: <strong id="provider-name">Carregando...</strong></span>
                <span class="status" id="connection-status">● Conectado</span>
            </div>
            <a href="/" class="dashboard-btn">📊 Dashboard</a>
        </header>

        <div id="alertContainer"></div>

        <div class="card">
            <h2 style="margin-bottom: 20px; color: #333;">1. Validar Cliente</h2>

            <div class="form-group">
                <label for="cliCodigo">Código do Cliente</label>
                <input type="text" id="cliCodigo" placeholder="Ex: 12345">
            </div>

            <div class="form-group">
                <label for="cliCpfCnpj">OU CPF/CNPJ do Cliente</label>
                <input type="text" id="cliCpfCnpj" placeholder="Ex: 123.456.789-00 ou 12.345.678/0001-90">
            </div>

            <button class="btn-primary" onclick="validarCliente()">
                Validar Cliente
            </button>

            <div id="clientInfo" class="client-info"></div>
        </div>

        <div class="card" id="messageCard" style="display: none;">
            <h2 style="margin-bottom: 20px; color: #333;">2. Compor E-mail</h2>

            <div class="form-group">
                <label for="emailDestinatario">E-mail do Destinatário</label>
                <input type="email" id="emailDestinatario" placeholder="exemplo@email.com">
            </div>

            <div class="form-group">
                <label for="assunto">Assunto</label>
                <input type="text" id="assunto" placeholder="Digite o assunto do e-mail" maxlength="500">
            </div>

            <div class="form-group">
                <label for="mensagem">Mensagem</label>
                <textarea id="mensagem" placeholder="Digite sua mensagem aqui..."></textarea>
            </div>

            <div class="form-group">
                <div class="checkbox-group">
                    <input type="checkbox" id="isHtml" onchange="toggleHtmlMode()">
                    <label for="isHtml" style="margin-bottom: 0;">Enviar como HTML</label>
                </div>
                <div class="hint">Marque esta opção se sua mensagem contiver formatação HTML</div>
            </div>

            <div class="form-group" id="anexoFileGroup">
                <label for="anexo">Anexo (opcional)</label>
                <input type="file" id="anexo" onchange="handleFileSelect(event)" accept="*/*">
                <div class="hint" id="anexoInfo">Tamanho máximo: 10MB</div>
            </div>

            <div class="form-group" id="anexoUrlGroup" style="display: none;">
                <label for="anexoUrl">URL do Anexo (somente Zenvia)</label>
                <input type="url" id="anexoUrl" placeholder="https://exemplo.com/arquivo.pdf" onchange="handleUrlInput()">
                <div class="hint" id="anexoUrlInfo">Informe a URL pública do arquivo a ser anexado</div>
            </div>

            <button class="btn-primary" onclick="dispararEmail()" id="btnEnviar">
                Enviar E-mail
            </button>

            <div id="statusTracking" class="status-tracking"></div>
        </div>
    </div>

    <script>
        let clienteValidado = null;
        let emailIdEnviado = null;
        let statusCheckInterval = null;
        let validandoCliente = false;
        let selectedAttachment = null;
        let selectedAttachmentUrl = "";
        let currentProvider = "";

        function showAlert(message, type) {
            const container = document.getElementById('alertContainer');
            const alert = document.createElement('div');
            alert.className = 'alert alert-' + type + ' show';
            alert.textContent = message;
            container.innerHTML = '';
            container.appendChild(alert);

            setTimeout(() => {
                alert.classList.remove('show');
                setTimeout(() => alert.remove(), 300);
            }, 5000);
        }

        async function validarCliente() {
            if (validandoCliente) {
                console.log('Validação já em andamento, ignorando chamada duplicada');
                return;
            }

            if (!verificarConexao()) {
                return;
            }

            validandoCliente = true;

            const cliCodigo = document.getElementById('cliCodigo').value.trim();
            const cliCpfCnpj = document.getElementById('cliCpfCnpj').value.trim();

            if (!cliCodigo && !cliCpfCnpj) {
                showAlert('Por favor, informe o código ou CPF/CNPJ do cliente', 'error');
                validandoCliente = false;
                return;
            }

            try {
                const response = await fetch('/api/manual/validar-cliente', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        cliCodigo: cliCodigo,
                        cliCpfCnpj: cliCpfCnpj
                    })
                });

                const data = await response.json();

                if (data && data.success) {
                    clienteValidado = data;
                    mostrarInfoCliente(data);

                    if (data.emailValido) {
                        document.getElementById('messageCard').style.display = 'block';
                        document.getElementById('emailDestinatario').value = data.email;
                        showAlert('Cliente validado com sucesso!', 'success');
                    } else {
                        showAlert('Cliente encontrado, mas o e-mail é inválido', 'error');
                    }
                } else {
                    showAlert(data && data.error ? data.error : 'Cliente não encontrado', 'error');
                    clienteValidado = null;
                    document.getElementById('clientInfo').classList.remove('show');
                    document.getElementById('messageCard').style.display = 'none';
                }
            } catch (error) {
                console.error('Erro ao validar cliente:', error);
                showAlert('Erro ao validar cliente: ' + error.message, 'error');
            } finally {
                validandoCliente = false;
            }
        }

        function mostrarInfoCliente(data) {
            const infoDiv = document.getElementById('clientInfo');
            const emailStatus = data.emailValido
                ? '<span style="color: #4caf50;">&#x2713; Válido</span>'
                : '<span style="color: #f44336;">&#x2717; Inválido</span>';

            const htmlContent = '<h3>&#x2713; Cliente Encontrado</h3>' +
                '<div class="info-row">' +
                    '<div class="info-label">Código:</div>' +
                    '<div class="info-value">' + data.cliCodigo + '</div>' +
                '</div>' +
                '<div class="info-row">' +
                    '<div class="info-label">CPF/CNPJ:</div>' +
                    '<div class="info-value">' + data.cliCpfCnpj + '</div>' +
                '</div>' +
                '<div class="info-row">' +
                    '<div class="info-label">Nome:</div>' +
                    '<div class="info-value">' + data.cliNome + '</div>' +
                '</div>' +
                '<div class="info-row">' +
                    '<div class="info-label">E-mail:</div>' +
                    '<div class="info-value">' + (data.email || 'Não cadastrado') + ' ' + emailStatus + '</div>' +
                '</div>';

            infoDiv.innerHTML = htmlContent;
            infoDiv.classList.add('show');
        }

        function toggleHtmlMode() {
            const isHtml = document.getElementById('isHtml').checked;
            const mensagemField = document.getElementById('mensagem');

            if (isHtml) {
                mensagemField.placeholder = 'Digite HTML aqui, ex: <h1>Título</h1><p>Parágrafo</p>';
            } else {
                mensagemField.placeholder = 'Digite sua mensagem aqui...';
            }
        }

        function handleFileSelect(event) {
            const file = event.target.files[0];
            const anexoInfo = document.getElementById('anexoInfo');

            if (!file) {
                selectedAttachment = null;
                anexoInfo.textContent = 'Tamanho máximo: 10MB';
                anexoInfo.style.color = '';
                return;
            }

            // Validar tamanho (10MB = 10 * 1024 * 1024 bytes)
            const maxSize = 10 * 1024 * 1024;
            if (file.size > maxSize) {
                showAlert('Arquivo muito grande. Tamanho máximo: 10MB', 'error');
                event.target.value = '';
                selectedAttachment = null;
                anexoInfo.textContent = 'Tamanho máximo: 10MB';
                anexoInfo.style.color = '';
                return;
            }

            // Ler arquivo e converter para base64
            const reader = new FileReader();
            reader.onload = function(e) {
                const base64Data = e.target.result.split(',')[1]; // Remover prefixo "data:...;base64,"

                selectedAttachment = {
                    data: base64Data,
                    name: file.name,
                    type: file.type || 'application/octet-stream'
                };

                // Atualizar informação visual
                const sizeKB = (file.size / 1024).toFixed(2);
                anexoInfo.textContent = '✓ ' + file.name + ' (' + sizeKB + ' KB)';
                anexoInfo.style.color = '#4caf50';

                console.log('Arquivo selecionado:', file.name, 'Tipo:', file.type, 'Tamanho:', sizeKB + 'KB');
            };

            reader.onerror = function() {
                showAlert('Erro ao ler arquivo', 'error');
                event.target.value = '';
                selectedAttachment = null;
                anexoInfo.textContent = 'Tamanho máximo: 10MB';
                anexoInfo.style.color = '';
            };

            reader.readAsDataURL(file);
        }

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

        async function dispararEmail() {
            if (!verificarConexao()) {
                return;
            }

            if (!clienteValidado || !clienteValidado.emailValido) {
                showAlert('Por favor, valide um cliente com e-mail válido primeiro', 'error');
                return;
            }

            const emailDestinatario = document.getElementById('emailDestinatario').value.trim();
            const assunto = document.getElementById('assunto').value.trim();
            const mensagem = document.getElementById('mensagem').value.trim();
            const isHtml = document.getElementById('isHtml').checked;

            if (!emailDestinatario) {
                showAlert('Por favor, informe o e-mail do destinatário', 'error');
                return;
            }

            if (!assunto) {
                showAlert('Por favor, informe o assunto', 'error');
                return;
            }

            if (!mensagem) {
                showAlert('Por favor, digite uma mensagem', 'error');
                return;
            }

            const btnEnviar = document.getElementById('btnEnviar');
            btnEnviar.disabled = true;
            btnEnviar.innerHTML = 'Enviando... <span class="spinner"></span>';

            try {
                const response = await fetch('/api/manual/disparar', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        cliCodigo: clienteValidado.cliCodigo,
                        email: emailDestinatario,
                        assunto: assunto,
                        mensagem: mensagem,
                        isHtml: isHtml,
                        attachmentData: selectedAttachment?.data || "",
                        attachmentName: selectedAttachment?.name || "",
                        attachmentType: selectedAttachment?.type || "",
                        attachmentUrl: selectedAttachmentUrl || ""
                    })
                });

                const data = await response.json();

                if (data.success) {
                    emailIdEnviado = data.emailId;
                    showAlert(data.message, 'success');
                    document.getElementById('mensagem').value = '';
                    document.getElementById('assunto').value = '';

                    // Limpar anexo selecionado
                    document.getElementById('anexo').value = '';
                    selectedAttachment = null;
                    const anexoInfo = document.getElementById('anexoInfo');
                    anexoInfo.textContent = 'Tamanho máximo: 10MB';
                    anexoInfo.style.color = '';

                    // Limpar URL de anexo
                    document.getElementById('anexoUrl').value = '';
                    selectedAttachmentUrl = "";
                    const anexoUrlInfo = document.getElementById('anexoUrlInfo');
                    anexoUrlInfo.textContent = 'Informe a URL pública do arquivo a ser anexado';
                    anexoUrlInfo.style.color = '';

                    iniciarAcompanhamentoStatus(data.emailId);
                } else {
                    showAlert(data.error || 'Erro ao enviar e-mail', 'error');
                    btnEnviar.disabled = false;
                    btnEnviar.textContent = 'Enviar E-mail';
                }
            } catch (error) {
                console.error('Erro ao disparar e-mail:', error);
                showAlert('Erro ao disparar e-mail. Tente novamente.', 'error');
                btnEnviar.disabled = false;
                btnEnviar.textContent = 'Enviar E-mail';
            }
        }

        function iniciarAcompanhamentoStatus(emailId) {
            const trackingDiv = document.getElementById('statusTracking');
            trackingDiv.innerHTML = '<h3>&#x1F4CA; Acompanhamento do Envio</h3>' +
                '<div class="info-row">' +
                    '<div class="info-label">ID do E-mail:</div>' +
                    '<div class="info-value">' + emailId + '</div>' +
                '</div>' +
                '<div class="info-row">' +
                    '<div class="info-label">Status:</div>' +
                    '<div class="info-value" id="statusValue">' +
                        '<span class="status-badge status-pendente">Aguardando processamento...</span>' +
                    '</div>' +
                '</div>';
            trackingDiv.classList.add('show');

            statusCheckInterval = setInterval(() => consultarStatus(emailId), 2000);

            setTimeout(() => {
                if (statusCheckInterval) {
                    clearInterval(statusCheckInterval);
                    statusCheckInterval = null;
                }
            }, 60000);
        }

        async function consultarStatus(emailId) {
            try {
                const response = await fetch('/api/manual/status?id=' + emailId);
                const data = await response.json();

                if (data.success) {
                    atualizarStatusDisplay(data);

                    // Para a consulta se o status for final
                    if (data.status === 2 || data.status === 4 || data.status === 125) {
                        if (statusCheckInterval) {
                            clearInterval(statusCheckInterval);
                            statusCheckInterval = null;
                        }

                        const btnEnviar = document.getElementById('btnEnviar');
                        btnEnviar.disabled = false;
                        btnEnviar.textContent = 'Enviar E-mail';
                    }
                }
            } catch (error) {
                console.error('Erro ao consultar status:', error);
            }
        }

        function atualizarStatusDisplay(data) {
            const statusValue = document.getElementById('statusValue');
            let statusHTML = '';
            let badgeClass = 'status-pendente';

            switch(data.status) {
                case 0:
                    badgeClass = 'status-pendente';
                    statusHTML = 'Pendente';
                    break;
                case 2:
                    badgeClass = 'status-enviado';
                    statusHTML = '✓ Enviado com sucesso';
                    break;
                case 125:
                    badgeClass = 'status-invalido';
                    statusHTML = '✗ E-mail inválido';
                    break;
                case 3:
                    badgeClass = 'status-erro';
                    statusHTML = '⚠ Erro (tentando novamente...)';
                    break;
                case 4:
                    badgeClass = 'status-erro';
                    statusHTML = '✗ Falha permanente';
                    break;
                default:
                    statusHTML = data.statusDesc;
            }

            let html = '<span class="status-badge ' + badgeClass + '">' + statusHTML + '</span>';

            if (data.dataEnvio) {
                html += '<br><small style="color: #666;">Enviado em: ' + data.dataEnvio + '</small>';
            }

            if (data.tentativas > 0) {
                html += '<br><small style="color: #666;">Tentativas: ' + data.tentativas + '</small>';
            }

            if (data.erroMsg) {
                html += '<br><small style="color: #f44336;">Erro: ' + data.erroMsg + '</small>';
            }

            if (data.idProvedor) {
                html += '<br><small style="color: #666;">ID Provedor: ' + data.idProvedor + '</small>';
            }

            statusValue.innerHTML = html;
        }

        let serviceConnected = true;
        let providerCheckInterval = null;

        async function carregarProviderInfo() {
            try {
                const response = await fetch('/api/manual/provider-info');
                if (response.ok) {
                    const data = await response.json();

                    const providerNames = {
                        'mock': '🧪 Mock (Teste)',
                        'smtp': '📧 SMTP',
                        'sendgrid': '📨 SendGrid',
                        'zenvia': '🇧🇷 Zenvia',
                        'pontaltech': '📡 Pontaltech'
                    };

                    const displayName = providerNames[data.providerName] || data.providerName.toUpperCase();
                    document.getElementById('provider-name').textContent = displayName;

                    // Atualizar provider atual e alternar campos de anexo
                    currentProvider = data.providerName;
                    toggleAnexoFields(currentProvider);

                    const statusElement = document.getElementById('connection-status');
                    if (data.status === 'online') {
                        statusElement.textContent = '● Conectado';
                        statusElement.style.background = '#10b981';
                        serviceConnected = true;
                    } else {
                        statusElement.textContent = '● Desconectado';
                        statusElement.style.background = '#ef4444';
                        serviceConnected = false;
                    }
                    return true;
                } else {
                    throw new Error('Falha ao carregar informações');
                }
            } catch (error) {
                console.error('Erro ao carregar informações do provider:', error);
                document.getElementById('provider-name').textContent = 'Erro ao carregar';
                document.getElementById('connection-status').textContent = '● Desconectado';
                document.getElementById('connection-status').style.background = '#ef4444';
                serviceConnected = false;
                return false;
            }
        }

        function toggleAnexoFields(providerName) {
            const fileGroup = document.getElementById('anexoFileGroup');
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

        function verificarConexao() {
            if (!serviceConnected) {
                showAlert('⚠️ Serviço Desconectado: O serviço de e-mail está desconectado no momento. Por favor, aguarde a reconexão.', 'error');
                return false;
            }
            return true;
        }

        window.addEventListener('beforeunload', () => {
            if (statusCheckInterval) {
                clearInterval(statusCheckInterval);
            }
            if (providerCheckInterval) {
                clearInterval(providerCheckInterval);
            }
        });

        carregarProviderInfo();
        providerCheckInterval = setInterval(carregarProviderInfo, 5000);
    </script>
</body>
</html>
`
