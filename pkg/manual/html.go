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
        select,
        textarea {
            width: 100%;
            padding: 12px;
            border: 2px solid #e0e0e0;
            border-radius: 8px;
            font-size: 1rem;
            transition: border-color 0.3s;
        }

        select {
            cursor: pointer;
            background-color: white;
        }

        input[type="text"]:focus,
        input[type="email"]:focus,
        select:focus,
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

        /* Estilos para campos de anexo - aumentados */
        input[type="file"],
        input[type="url"] {
            width: 100%;
            padding: 16px;
            border: 2px dashed #667eea;
            border-radius: 8px;
            font-size: 1rem;
            background: #f8f9ff;
            cursor: pointer;
            transition: all 0.3s;
        }

        input[type="file"]:hover,
        input[type="url"]:hover {
            border-color: #764ba2;
            background: #f0f2ff;
        }

        input[type="file"]:focus,
        input[type="url"]:focus {
            outline: none;
            border-color: #667eea;
            border-style: solid;
            background: white;
        }

        /* Melhorar visibilidade do campo de anexo */
        #anexoFileGroup,
        #anexoUrlGroup {
            background: #fafbff;
            padding: 20px;
            border-radius: 10px;
            border: 1px solid #e8eaf6;
        }

        #anexoFileGroup label,
        #anexoUrlGroup label {
            font-size: 1.1rem;
            color: #667eea;
            font-weight: 600;
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

        /* Modal de Preview */
        .modal {
            display: none;
            position: fixed;
            z-index: 1000;
            left: 0;
            top: 0;
            width: 100%;
            height: 100%;
            background-color: rgba(0,0,0,0.6);
            overflow: auto;
        }

        .modal-content {
            background-color: #fefefe;
            margin: 2% auto;
            padding: 0;
            border-radius: 12px;
            width: 90%;
            max-width: 1000px;
            max-height: 90vh;
            box-shadow: 0 10px 50px rgba(0,0,0,0.3);
            display: flex;
            flex-direction: column;
        }

        .modal-header {
            padding: 20px 30px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border-radius: 12px 12px 0 0;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        .modal-header h3 {
            margin: 0;
            font-size: 1.5rem;
        }

        .modal-close {
            color: white;
            font-size: 32px;
            font-weight: bold;
            cursor: pointer;
            background: rgba(255, 255, 255, 0.2);
            border: 2px solid rgba(255, 255, 255, 0.3);
            border-radius: 8px;
            padding: 0;
            width: 40px;
            height: 40px;
            display: flex;
            align-items: center;
            justify-content: center;
            transition: all 0.3s;
            line-height: 1;
            position: relative;
        }

        .modal-close:hover {
            background: rgba(255, 255, 255, 0.3);
            border-color: rgba(255, 255, 255, 0.5);
            transform: scale(1.1);
        }

        .modal-close:active {
            transform: scale(0.95);
        }

        .modal-body {
            padding: 30px;
            overflow-y: auto;
            flex: 1;
        }

        .preview-section {
            margin-bottom: 25px;
        }

        .preview-section h4 {
            color: #667eea;
            margin-bottom: 10px;
            font-size: 1.1rem;
        }

        .preview-content {
            background: #f8f9fa;
            padding: 15px;
            border-radius: 8px;
            border: 1px solid #e0e0e0;
            min-height: 100px;
            max-height: 400px;
            overflow-y: auto;
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
                <label for="templateSelect">Modo de Composição</label>
                <select id="templateSelect" onchange="handleTemplateChange()">
                    <option value="0">✍️ Compor manualmente</option>
                </select>
                <div class="hint">Escolha um template ou componha seu e-mail manualmente</div>
            </div>

            <div id="templatePreviewGroup" class="form-group" style="display: none;">
                <button type="button" class="btn-secondary" onclick="previewTemplate()" id="btnPreview">
                    👁️ Visualizar Preview do Template
                </button>
            </div>

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

    <!-- Modal de Preview -->
    <div id="previewModal" class="modal">
        <div class="modal-content">
            <div class="modal-header">
                <h3>📧 Preview do E-mail</h3>
                <button class="modal-close" onclick="closePreviewModal()">&times;</button>
            </div>
            <div class="modal-body">
                <div class="preview-section">
                    <h4>Assunto:</h4>
                    <div class="preview-content" id="previewAssunto"></div>
                </div>
                <div class="preview-section">
                    <h4>Corpo do E-mail:</h4>
                    <div class="preview-content" id="previewCorpo"></div>
                </div>
            </div>
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
        let selectedTemplateId = 0;
        let availableTemplates = [];
        let serviceConnected = false;
        let providerCheckInterval = null;

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

            // Se não está usando template, validar campos manuais
            if (selectedTemplateId === 0) {
                if (!assunto) {
                    showAlert('Por favor, informe o assunto', 'error');
                    return;
                }

                if (!mensagem) {
                    showAlert('Por favor, digite uma mensagem', 'error');
                    return;
                }
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
                        templateId: selectedTemplateId,
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

        // Funções de Template
        async function carregarTemplates() {
            try {
                const response = await fetch('/api/templates?activeOnly=true&limit=100');
                if (response.ok) {
                    const data = await response.json();
                    availableTemplates = data.data || [];
                    preencherSeletorTemplates();
                    console.log('Templates carregados:', availableTemplates.length);
                } else {
                    console.error('Erro ao carregar templates');
                }
            } catch (error) {
                console.error('Erro ao carregar templates:', error);
            }
        }

        function preencherSeletorTemplates() {
            const select = document.getElementById('templateSelect');
            select.innerHTML = '<option value="0">✍️ Compor manualmente</option>';

            availableTemplates.forEach(template => {
                const option = document.createElement('option');
                option.value = template.id;
                option.textContent = '📧 ' + template.nome;
                select.appendChild(option);
            });
        }

        function handleTemplateChange() {
            const select = document.getElementById('templateSelect');
            selectedTemplateId = parseInt(select.value);

            const previewGroup = document.getElementById('templatePreviewGroup');
            const assuntoField = document.getElementById('assunto');
            const mensagemField = document.getElementById('mensagem');
            const isHtmlCheckbox = document.getElementById('isHtml');

            if (selectedTemplateId > 0) {
                // Modo template
                previewGroup.style.display = 'block';

                // Buscar template selecionado
                const template = availableTemplates.find(t => t.id === selectedTemplateId);
                if (template) {
                    assuntoField.value = template.assuntoPadrao || '';
                    assuntoField.readOnly = true;
                    mensagemField.value = '[Template será processado automaticamente com dados do cliente]';
                    mensagemField.readOnly = true;
                    isHtmlCheckbox.checked = true;
                    isHtmlCheckbox.disabled = true;
                }
            } else {
                // Modo manual
                previewGroup.style.display = 'none';
                assuntoField.readOnly = false;
                mensagemField.readOnly = false;
                mensagemField.value = '';
                isHtmlCheckbox.disabled = false;
            }
        }

        async function previewTemplate() {
            console.log('=== Preview Template ===');
            console.log('Cliente validado:', clienteValidado);
            console.log('Template ID selecionado:', selectedTemplateId);

            if (!clienteValidado) {
                showAlert('⚠️ Por favor, valide um cliente primeiro', 'error');
                return;
            }

            if (selectedTemplateId === 0) {
                showAlert('⚠️ Nenhum template selecionado', 'error');
                return;
            }

            const btnPreview = document.getElementById('btnPreview');
            btnPreview.disabled = true;
            btnPreview.textContent = 'Carregando preview...';

            const requestData = {
                templateId: selectedTemplateId,
                cliCodigo: clienteValidado.cliCodigo
            };
            console.log('Request data:', requestData);

            try {
                console.log('Fazendo requisição para /api/manual/preview-template...');
                const response = await fetch('/api/manual/preview-template', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(requestData)
                });

                console.log('Response status:', response.status);
                const data = await response.json();
                console.log('Response data:', data);

                if (data.success) {
                    console.log('Preview gerado com sucesso!');
                    console.log('Assunto:', data.assunto);
                    console.log('Corpo (primeiros 200 chars):', data.corpo.substring(0, 200));

                    document.getElementById('previewAssunto').textContent = data.assunto;
                    document.getElementById('previewCorpo').innerHTML = data.corpo;
                    document.getElementById('previewModal').style.display = 'block';
                } else {
                    console.error('Erro no preview:', data.error);
                    showAlert('❌ Erro ao gerar preview: ' + (data.error || 'Erro desconhecido'), 'error');
                }
            } catch (error) {
                console.error('Erro ao fazer preview:', error);
                showAlert('❌ Erro ao gerar preview do template: ' + error.message, 'error');
            } finally {
                btnPreview.disabled = false;
                btnPreview.textContent = '👁️ Visualizar Preview do Template';
            }
        }

        function closePreviewModal() {
            document.getElementById('previewModal').style.display = 'none';
        }

        // Fechar modal ao clicar fora
        window.onclick = function(event) {
            const modal = document.getElementById('previewModal');
            if (event.target === modal) {
                closePreviewModal();
            }
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
        carregarTemplates();
    </script>
</body>
</html>
`
