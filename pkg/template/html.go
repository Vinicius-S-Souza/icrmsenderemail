package template

import "net/http"

// ServeTemplateList serve a página de listagem de templates
func (h *Handler) ServeTemplateList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(templateListHTML))
}

// ServeTemplateEditor serve a página de edição/criação de template
func (h *Handler) ServeTemplateEditor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(templateEditorHTML))
}

// templateListHTML contém o HTML da página de listagem de templates
// Data de criação: 12/12/2025 18:45
const templateListHTML = `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Templates de E-mail - ICRMSenderEmail</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: #333;
            padding: 20px;
            min-height: 100vh;
        }

        .container {
            max-width: 1400px;
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

        .header-left h1 {
            color: #667eea;
            font-size: 2em;
            margin-bottom: 5px;
        }

        .subtitle {
            color: #666;
            font-size: 1em;
        }

        .header-actions {
            display: flex;
            gap: 10px;
        }

        .btn {
            display: inline-block;
            padding: 12px 24px;
            border-radius: 8px;
            font-weight: 600;
            text-decoration: none;
            cursor: pointer;
            border: none;
            font-size: 14px;
            transition: all 0.3s;
        }

        .btn-primary {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
        }

        .btn-primary:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 16px rgba(102, 126, 234, 0.4);
        }

        .btn-secondary {
            background: #f3f4f6;
            color: #667eea;
        }

        .btn-secondary:hover {
            background: #e5e7eb;
        }

        .btn-danger {
            background: #ef4444;
            color: white;
        }

        .btn-danger:hover {
            background: #dc2626;
        }

        .btn-small {
            padding: 6px 12px;
            font-size: 12px;
        }

        .content-card {
            background: white;
            border-radius: 15px;
            box-shadow: 0 5px 20px rgba(0,0,0,0.1);
            padding: 30px;
        }

        .search-bar {
            margin-bottom: 20px;
            display: flex;
            gap: 10px;
        }

        .search-bar input {
            flex: 1;
            padding: 12px;
            border: 1px solid #e5e7eb;
            border-radius: 8px;
            font-size: 14px;
        }

        .search-bar input:focus {
            outline: none;
            border-color: #667eea;
        }

        table {
            width: 100%;
            border-collapse: collapse;
        }

        thead {
            background: #f9fafb;
        }

        th {
            text-align: left;
            padding: 12px;
            font-weight: 600;
            color: #666;
            border-bottom: 2px solid #e5e7eb;
        }

        td {
            padding: 12px;
            border-bottom: 1px solid #e5e7eb;
        }

        tr:hover {
            background: #f9fafb;
        }

        .status-badge {
            display: inline-block;
            padding: 4px 12px;
            border-radius: 20px;
            font-size: 12px;
            font-weight: 600;
        }

        .status-active {
            background: #d1fae5;
            color: #065f46;
        }

        .status-inactive {
            background: #fee2e2;
            color: #991b1b;
        }

        .actions {
            display: flex;
            gap: 5px;
        }

        .pagination {
            display: flex;
            justify-content: center;
            align-items: center;
            gap: 10px;
            margin-top: 20px;
        }

        .pagination button {
            padding: 8px 16px;
            border: 1px solid #e5e7eb;
            background: white;
            border-radius: 6px;
            cursor: pointer;
        }

        .pagination button:disabled {
            opacity: 0.5;
            cursor: not-allowed;
        }

        .pagination span {
            color: #666;
        }

        .empty-state {
            text-align: center;
            padding: 60px 20px;
            color: #999;
        }

        .empty-state-icon {
            font-size: 64px;
            margin-bottom: 20px;
        }

        .loading {
            text-align: center;
            padding: 40px;
            color: #667eea;
        }

        .error-message {
            background: #fee2e2;
            color: #991b1b;
            padding: 12px;
            border-radius: 8px;
            margin-bottom: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <div class="header-left">
                <h1>📝 Templates de E-mail</h1>
                <p class="subtitle">Gerencie templates HTML para seus e-mails</p>
            </div>
            <div class="header-actions">
                <a href="/" class="btn btn-secondary">← Dashboard</a>
                <a href="/templates/novo" class="btn btn-primary">+ Novo Template</a>
            </div>
        </header>

        <div class="content-card">
            <div class="search-bar">
                <input type="text" id="searchInput" placeholder="Buscar templates..." onkeyup="handleSearch()">
                <button class="btn btn-secondary" onclick="loadTemplates()">🔄 Atualizar</button>
            </div>

            <div id="errorMessage" style="display: none;" class="error-message"></div>
            <div id="loadingState" class="loading">Carregando templates...</div>

            <div id="tableContainer" style="display: none;">
                <table id="templatesTable">
                    <thead>
                        <tr>
                            <th>Nome</th>
                            <th>Descrição</th>
                            <th>Status</th>
                            <th>Data Criação</th>
                            <th>Ações</th>
                        </tr>
                    </thead>
                    <tbody id="templatesBody">
                    </tbody>
                </table>

                <div class="pagination">
                    <button id="prevPage" onclick="previousPage()" disabled>← Anterior</button>
                    <span id="pageInfo">Página 1 de 1</span>
                    <button id="nextPage" onclick="nextPage()" disabled>Próxima →</button>
                </div>
            </div>

            <div id="emptyState" style="display: none;" class="empty-state">
                <div class="empty-state-icon">📭</div>
                <h2>Nenhum template encontrado</h2>
                <p>Crie seu primeiro template para começar!</p>
                <br>
                <a href="/templates/novo" class="btn btn-primary">+ Criar Template</a>
            </div>
        </div>
    </div>

    <script>
        let currentPage = 1;
        const limit = 10;
        let totalPages = 1;
        let searchTimeout = null;

        // Carregar templates ao iniciar
        document.addEventListener('DOMContentLoaded', function() {
            loadTemplates();
        });

        function loadTemplates() {
            const searchTerm = document.getElementById('searchInput').value;

            document.getElementById('loadingState').style.display = 'block';
            document.getElementById('tableContainer').style.display = 'none';
            document.getElementById('emptyState').style.display = 'none';
            document.getElementById('errorMessage').style.display = 'none';

            fetch('/api/templates?page=' + currentPage + '&limit=' + limit + '&search=' + encodeURIComponent(searchTerm))
                .then(response => response.json())
                .then(data => {
                    document.getElementById('loadingState').style.display = 'none';

                    if (!data.success) {
                        showError(data.error || 'Erro ao carregar templates');
                        return;
                    }

                    if (!data.data || data.data.length === 0) {
                        document.getElementById('emptyState').style.display = 'block';
                        return;
                    }

                    totalPages = data.totalPages || 1;
                    renderTemplates(data.data);
                    updatePagination();
                })
                .catch(error => {
                    console.error('Erro:', error);
                    document.getElementById('loadingState').style.display = 'none';
                    showError('Erro ao conectar com o servidor');
                });
        }

        function renderTemplates(templates) {
            const tbody = document.getElementById('templatesBody');
            tbody.innerHTML = '';

            templates.forEach(template => {
                const tr = document.createElement('tr');

                const statusClass = template.ativo ? 'status-active' : 'status-inactive';
                const statusText = template.ativo ? 'Ativo' : 'Inativo';

                tr.innerHTML = '<td><strong>' + template.nome + '</strong></td>' +
                    '<td>' + (template.descricao || '-') + '</td>' +
                    '<td><span class="status-badge ' + statusClass + '">' + statusText + '</span></td>' +
                    '<td>' + template.dataCriacao + '</td>' +
                    '<td>' +
                        '<div class="actions">' +
                            '<button class="btn btn-primary btn-small" onclick="editTemplate(' + template.id + ')">✏️ Editar</button>' +
                            '<button class="btn btn-secondary btn-small" onclick="duplicateTemplate(' + template.id + ', \'' + template.nome + '\')">📋 Duplicar</button>' +
                            '<button class="btn btn-danger btn-small" onclick="deleteTemplate(' + template.id + ', \'' + template.nome + '\')">🗑️ Excluir</button>' +
                        '</div>' +
                    '</td>';

                tbody.appendChild(tr);
            });

            document.getElementById('tableContainer').style.display = 'block';
        }

        function updatePagination() {
            document.getElementById('pageInfo').textContent = 'Página ' + currentPage + ' de ' + totalPages;
            document.getElementById('prevPage').disabled = currentPage === 1;
            document.getElementById('nextPage').disabled = currentPage >= totalPages;
        }

        function previousPage() {
            if (currentPage > 1) {
                currentPage--;
                loadTemplates();
            }
        }

        function nextPage() {
            if (currentPage < totalPages) {
                currentPage++;
                loadTemplates();
            }
        }

        function handleSearch() {
            clearTimeout(searchTimeout);
            searchTimeout = setTimeout(() => {
                currentPage = 1;
                loadTemplates();
            }, 500);
        }

        function editTemplate(id) {
            window.location.href = '/templates/' + id + '/editar';
        }

        function duplicateTemplate(id, nome) {
            const newName = prompt('Digite o nome do novo template:', nome + ' (cópia)');
            if (!newName) return;

            fetch('/api/templates/' + id + '/duplicate', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    newName: newName
                })
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    alert('Template duplicado com sucesso!');
                    loadTemplates();
                } else {
                    alert('Erro ao duplicar template: ' + (data.error || 'Erro desconhecido'));
                }
            })
            .catch(error => {
                console.error('Erro:', error);
                alert('Erro ao duplicar template');
            });
        }

        function deleteTemplate(id, nome) {
            if (!confirm('Tem certeza que deseja excluir o template "' + nome + '"?')) {
                return;
            }

            fetch('/api/templates/' + id, {
                method: 'DELETE'
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    alert('Template excluído com sucesso!');
                    loadTemplates();
                } else {
                    alert('Erro ao excluir template: ' + (data.error || 'Erro desconhecido'));
                }
            })
            .catch(error => {
                console.error('Erro:', error);
                alert('Erro ao excluir template');
            });
        }

        function showError(message) {
            const errorDiv = document.getElementById('errorMessage');
            errorDiv.textContent = message;
            errorDiv.style.display = 'block';
        }
    </script>
</body>
</html>
`

// templateEditorHTML contém o HTML da página de edição/criação de templates
// Data de criação: 12/12/2025 18:50
const templateEditorHTML = `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Editor de Template - ICRMSenderEmail</title>

    <!-- Quill.js -->
    <link href="https://cdn.quilljs.com/1.3.6/quill.snow.css" rel="stylesheet">
    <script src="https://cdn.quilljs.com/1.3.6/quill.js"></script>

    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: #333;
            padding: 20px;
            min-height: 100vh;
        }

        .container {
            max-width: 1600px;
            margin: 0 auto;
        }

        header {
            background: white;
            padding: 20px 30px;
            border-radius: 15px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.1);
            margin-bottom: 20px;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        h1 {
            color: #667eea;
            font-size: 1.8em;
        }

        .header-actions {
            display: flex;
            gap: 10px;
        }

        .btn {
            padding: 10px 20px;
            border-radius: 8px;
            font-weight: 600;
            text-decoration: none;
            cursor: pointer;
            border: none;
            font-size: 14px;
            transition: all 0.3s;
        }

        .btn-primary {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
        }

        .btn-primary:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 16px rgba(102, 126, 234, 0.4);
        }

        .btn-secondary {
            background: #f3f4f6;
            color: #667eea;
        }

        .btn-secondary:hover {
            background: #e5e7eb;
        }

        .editor-layout {
            display: grid;
            grid-template-columns: 1fr 400px;
            gap: 20px;
        }

        .editor-main {
            background: white;
            border-radius: 15px;
            padding: 30px;
            box-shadow: 0 5px 20px rgba(0,0,0,0.1);
        }

        .editor-sidebar {
            background: white;
            border-radius: 15px;
            padding: 20px;
            box-shadow: 0 5px 20px rgba(0,0,0,0.1);
            max-height: calc(100vh - 120px);
            overflow-y: auto;
        }

        .form-group {
            margin-bottom: 20px;
        }

        label {
            display: block;
            margin-bottom: 8px;
            font-weight: 600;
            color: #333;
        }

        input[type="text"],
        textarea {
            width: 100%;
            padding: 10px;
            border: 1px solid #e5e7eb;
            border-radius: 6px;
            font-size: 14px;
            font-family: inherit;
        }

        input:focus,
        textarea:focus {
            outline: none;
            border-color: #667eea;
        }

        .tabs {
            display: flex;
            border-bottom: 2px solid #e5e7eb;
            margin-bottom: 20px;
        }

        .tab {
            padding: 12px 24px;
            cursor: pointer;
            border: none;
            background: none;
            font-size: 14px;
            font-weight: 600;
            color: #666;
            border-bottom: 3px solid transparent;
            margin-bottom: -2px;
        }

        .tab.active {
            color: #667eea;
            border-bottom-color: #667eea;
        }

        .tab-content {
            display: none;
        }

        .tab-content.active {
            display: block;
        }

        .editor-container {
            background: white;
            border: 1px solid #e5e7eb;
            border-radius: 8px;
            min-height: 300px;
        }

        .ql-editor {
            min-height: 250px;
        }

        .toggle-switch {
            display: flex;
            align-items: center;
            gap: 10px;
        }

        .switch {
            position: relative;
            display: inline-block;
            width: 50px;
            height: 24px;
        }

        .switch input {
            opacity: 0;
            width: 0;
            height: 0;
        }

        .slider {
            position: absolute;
            cursor: pointer;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background-color: #ccc;
            transition: .4s;
            border-radius: 24px;
        }

        .slider:before {
            position: absolute;
            content: "";
            height: 18px;
            width: 18px;
            left: 3px;
            bottom: 3px;
            background-color: white;
            transition: .4s;
            border-radius: 50%;
        }

        input:checked + .slider {
            background-color: #667eea;
        }

        input:checked + .slider:before {
            transform: translateX(26px);
        }

        .macro-list {
            display: flex;
            flex-direction: column;
            gap: 8px;
        }

        .macro-item {
            padding: 8px 12px;
            background: #f9fafb;
            border-radius: 6px;
            cursor: pointer;
            transition: all 0.2s;
            border: 1px solid #e5e7eb;
        }

        .macro-item:hover {
            background: #667eea;
            color: white;
        }

        .macro-key {
            font-family: monospace;
            font-weight: 600;
            font-size: 13px;
        }

        .macro-desc {
            font-size: 11px;
            color: #666;
            margin-top: 4px;
        }

        .macro-item:hover .macro-desc {
            color: rgba(255,255,255,0.8);
        }

        .preview-btn {
            width: 100%;
            margin-top: 10px;
        }

        /* Estilos para indicadores de tamanho */
        .size-indicator {
            margin-top: 15px;
            padding: 15px;
            background: #f9fafb;
            border-radius: 8px;
            border: 1px solid #e5e7eb;
        }

        .size-stat {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 10px;
            padding: 8px;
            background: white;
            border-radius: 6px;
        }

        .size-stat:last-child {
            margin-bottom: 0;
        }

        .size-label {
            font-size: 13px;
            color: #666;
        }

        .size-value {
            font-weight: 700;
            font-size: 14px;
        }

        .size-value.ok {
            color: #10b981;
        }

        .size-value.warning {
            color: #f59e0b;
        }

        .size-value.danger {
            color: #ef4444;
        }

        .size-bar-container {
            margin-top: 8px;
            height: 6px;
            background: #e5e7eb;
            border-radius: 3px;
            overflow: hidden;
        }

        .size-bar {
            height: 100%;
            transition: width 0.3s, background-color 0.3s;
            border-radius: 3px;
        }

        .size-bar.ok {
            background: linear-gradient(90deg, #10b981, #34d399);
        }

        .size-bar.warning {
            background: linear-gradient(90deg, #f59e0b, #fbbf24);
        }

        .size-bar.danger {
            background: linear-gradient(90deg, #ef4444, #f87171);
        }

        .size-warning {
            margin-top: 10px;
            padding: 10px;
            background: #fef3c7;
            border-left: 4px solid #f59e0b;
            border-radius: 4px;
            font-size: 12px;
            color: #92400e;
            display: none;
        }

        .size-warning.show {
            display: block;
        }

        .size-error {
            margin-top: 10px;
            padding: 10px;
            background: #fee2e2;
            border-left: 4px solid #ef4444;
            border-radius: 4px;
            font-size: 12px;
            color: #991b1b;
            display: none;
        }

        .size-error.show {
            display: block;
        }

        .alert {
            padding: 12px;
            border-radius: 8px;
            margin-bottom: 20px;
        }

        .alert-success {
            background: #d1fae5;
            color: #065f46;
        }

        .alert-error {
            background: #fee2e2;
            color: #991b1b;
        }

        .sidebar-section {
            margin-bottom: 24px;
            padding-bottom: 24px;
            border-bottom: 1px solid #e5e7eb;
        }

        .sidebar-section:last-child {
            border-bottom: none;
        }

        .sidebar-section h3 {
            font-size: 14px;
            margin-bottom: 12px;
            color: #667eea;
        }

        .hint {
            font-size: 12px;
            color: #999;
            margin-top: 4px;
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1 id="pageTitle">Novo Template</h1>
            <div class="header-actions">
                <a href="/templates" class="btn btn-secondary">← Voltar</a>
                <button class="btn btn-primary" onclick="saveTemplate()">💾 Salvar</button>
            </div>
        </header>

        <div id="alertContainer"></div>

        <div class="editor-layout">
            <!-- Editor Principal -->
            <div class="editor-main">
                <!-- Informações Básicas -->
                <div class="form-group">
                    <label for="templateName">Nome do Template *</label>
                    <input type="text" id="templateName" placeholder="Ex: Boas-vindas, Newsletter Mensal..." required>
                </div>

                <div class="form-group">
                    <label for="templateDesc">Descrição</label>
                    <textarea id="templateDesc" rows="2" placeholder="Breve descrição do template..."></textarea>
                </div>

                <div class="form-group">
                    <label for="templateSubject">Assunto Padrão</label>
                    <input type="text" id="templateSubject" placeholder="Ex: Bem-vindo à {{empresa}}!">
                    <div class="hint">Você pode usar macros no assunto</div>
                </div>

                <!-- Tabs de Seções -->
                <div class="tabs">
                    <button class="tab active" onclick="switchTab('header')">📄 Header</button>
                    <button class="tab" onclick="switchTab('body')">📝 Body</button>
                    <button class="tab" onclick="switchTab('footer')">📋 Footer</button>
                </div>

                <!-- Header Tab -->
                <div id="tab-header" class="tab-content active">
                    <div class="form-group">
                        <label>Cabeçalho do E-mail (opcional)</label>
                        <div id="headerEditor" class="editor-container"></div>
                    </div>
                </div>

                <!-- Body Tab -->
                <div id="tab-body" class="tab-content">
                    <div class="form-group">
                        <label>Corpo do E-mail *</label>
                        <div id="bodyEditor" class="editor-container"></div>
                    </div>
                </div>

                <!-- Footer Tab -->
                <div id="tab-footer" class="tab-content">
                    <div class="form-group">
                        <label>Rodapé do E-mail (opcional)</label>
                        <div id="footerEditor" class="editor-container"></div>
                    </div>
                </div>
            </div>

            <!-- Sidebar -->
            <div class="editor-sidebar">
                <!-- Status -->
                <div class="sidebar-section">
                    <h3>Status</h3>
                    <div class="toggle-switch">
                        <label class="switch">
                            <input type="checkbox" id="templateActive" checked>
                            <span class="slider"></span>
                        </label>
                        <span id="statusLabel">Ativo</span>
                    </div>
                </div>

                <!-- Macros Disponíveis -->
                <div class="sidebar-section">
                    <h3>Macros Disponíveis</h3>
                    <div class="hint" style="margin-bottom: 12px;">Clique para inserir no editor ativo</div>
                    <div id="macrosList" class="macro-list"></div>
                </div>

                <!-- Estatísticas de Tamanho -->
                <div class="sidebar-section">
                    <h3>📊 Tamanho do Template</h3>
                    <div class="size-indicator">
                        <div class="size-stat">
                            <span class="size-label">Tamanho Total</span>
                            <span class="size-value ok" id="totalSize">0 KB</span>
                        </div>
                        <div class="size-bar-container">
                            <div class="size-bar ok" id="totalSizeBar" style="width: 0%"></div>
                        </div>

                        <div class="size-stat" style="margin-top: 10px;">
                            <span class="size-label">Imagens Base64</span>
                            <span class="size-value ok" id="imageSize">0 KB</span>
                        </div>
                        <div class="size-bar-container">
                            <div class="size-bar ok" id="imageSizeBar" style="width: 0%"></div>
                        </div>

                        <div class="size-stat" style="margin-top: 10px;">
                            <span class="size-label">Limite Zenvia</span>
                            <span class="size-value" style="color: #666;">65 KB</span>
                        </div>

                        <div class="size-warning" id="sizeWarning">
                            ⚠️ <strong>Atenção:</strong> Template próximo do limite. Imagens base64 serão removidas automaticamente no envio via Zenvia.
                        </div>

                        <div class="size-error" id="sizeError">
                            ❌ <strong>Erro:</strong> Template excede o limite mesmo sem imagens. Reduza o conteúdo HTML.
                        </div>
                    </div>
                </div>

                <!-- Ações -->
                <div class="sidebar-section">
                    <h3>Ações</h3>
                    <button class="btn btn-secondary preview-btn" onclick="showPreview()">👁️ Visualizar Preview</button>
                </div>
            </div>
        </div>
    </div>

    <script>
        let headerEditor, bodyEditor, footerEditor;
        let currentEditor = 'header';
        let templateId = null;
        let macros = [];

        // Extrair ID da URL se estiver editando
        const pathParts = window.location.pathname.split('/');
        if (pathParts[2] && pathParts[2] !== 'novo') {
            templateId = parseInt(pathParts[2]);
        }

        // Inicializar editores Quill
        document.addEventListener('DOMContentLoaded', function() {
            const toolbarOptions = [
                ['bold', 'italic', 'underline', 'strike'],
                ['blockquote', 'code-block'],
                [{ 'header': 1 }, { 'header': 2 }],
                [{ 'list': 'ordered'}, { 'list': 'bullet' }],
                [{ 'align': [] }],
                ['link', 'image'],
                [{ 'color': [] }, { 'background': [] }],
                ['clean']
            ];

            // Handler customizado para imagens
            const imageHandler = function() {
                const editor = this.quill;

                // Perguntar ao usuário qual método usar
                const useUrl = confirm(
                    '📸 Inserir Imagem\n\n' +
                    'Clique em OK para inserir URL da imagem (recomendado)\n' +
                    'Clique em Cancelar para fazer upload (será convertido em base64)'
                );

                if (useUrl) {
                    // Usar URL
                    const url = prompt(
                        '🔗 Insira a URL da imagem:\n\n' +
                        'Exemplo: https://exemplo.com/imagem.jpg\n\n' +
                        '⚠️ A imagem deve estar hospedada em um servidor público e acessível.'
                    );

                    if (url) {
                        // Validar URL
                        if (!url.startsWith('http://') && !url.startsWith('https://')) {
                            alert('❌ URL inválida!\n\nA URL deve começar com http:// ou https://');
                            return;
                        }

                        // Inserir imagem com URL
                        const range = editor.getSelection(true);
                        editor.insertEmbed(range.index, 'image', url);
                        editor.setSelection(range.index + 1);

                        showAlert('✅ Imagem inserida via URL (não afeta limite de tamanho)', 'success');

                        // Atualizar estatísticas
                        setTimeout(() => updateSizeStats(), 100);
                    }
                } else {
                    // Upload de arquivo (base64)
                    const input = document.createElement('input');
                    input.setAttribute('type', 'file');
                    input.setAttribute('accept', 'image/*');
                    input.click();

                    input.onchange = () => {
                        const file = input.files[0];
                        if (file) {
                            // Verificar tamanho do arquivo
                            const maxSize = 2 * 1024 * 1024; // 2MB
                            if (file.size > maxSize) {
                                alert('❌ Arquivo muito grande!\n\nTamanho máximo: 2MB\nTamanho do arquivo: ' + (file.size / 1024 / 1024).toFixed(2) + 'MB\n\n💡 Use uma URL de imagem hospedada para imagens grandes.');
                                return;
                            }

                            const reader = new FileReader();
                            reader.onload = (e) => {
                                const range = editor.getSelection(true);
                                editor.insertEmbed(range.index, 'image', e.target.result);
                                editor.setSelection(range.index + 1);

                                showAlert('⚠️ Imagem convertida em base64 (aumenta o tamanho do template)', 'warning');

                                // Atualizar estatísticas
                                setTimeout(() => updateSizeStats(), 100);
                            };
                            reader.readAsDataURL(file);
                        }
                    };
                }
            };

            headerEditor = new Quill('#headerEditor', {
                theme: 'snow',
                modules: {
                    toolbar: {
                        container: toolbarOptions,
                        handlers: {
                            image: imageHandler
                        }
                    }
                },
                placeholder: 'Digite o HTML do cabeçalho...'
            });

            bodyEditor = new Quill('#bodyEditor', {
                theme: 'snow',
                modules: {
                    toolbar: {
                        container: toolbarOptions,
                        handlers: {
                            image: imageHandler
                        }
                    }
                },
                placeholder: 'Digite o conteúdo principal do e-mail...'
            });

            footerEditor = new Quill('#footerEditor', {
                theme: 'snow',
                modules: {
                    toolbar: {
                        container: toolbarOptions,
                        handlers: {
                            image: imageHandler
                        }
                    }
                },
                placeholder: 'Digite o rodapé do e-mail...'
            });

            // Monitorar mudanças nos editores para atualizar estatísticas de tamanho
            headerEditor.on('text-change', updateSizeStats);
            bodyEditor.on('text-change', updateSizeStats);
            footerEditor.on('text-change', updateSizeStats);

            // Toggle status label
            document.getElementById('templateActive').addEventListener('change', function(e) {
                document.getElementById('statusLabel').textContent = e.target.checked ? 'Ativo' : 'Inativo';
            });

            // Carregar macros
            loadMacros();

            // Atualizar estatísticas iniciais
            setTimeout(() => {
                updateSizeStats();
            }, 100);

            // Se está editando, carregar dados do template
            if (templateId) {
                document.getElementById('pageTitle').textContent = 'Editar Template';
                loadTemplate();
            }
        });

        function switchTab(tab) {
            // Atualizar tabs
            document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
            document.querySelectorAll('.tab-content').forEach(tc => tc.classList.remove('active'));

            event.target.classList.add('active');
            document.getElementById('tab-' + tab).classList.add('active');

            currentEditor = tab;
        }

        function loadMacros() {
            fetch('/api/templates/macros')
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        macros = data.data;
                        renderMacros();
                    }
                })
                .catch(error => console.error('Erro ao carregar macros:', error));
        }

        function renderMacros() {
            const container = document.getElementById('macrosList');
            container.innerHTML = '';

            macros.forEach(macro => {
                const div = document.createElement('div');
                div.className = 'macro-item';
                div.onclick = () => insertMacro(macro.key);
                div.innerHTML = '<div class="macro-key">' + macro.key + '</div>' +
                    '<div class="macro-desc">' + macro.description + '</div>';
                container.appendChild(div);
            });
        }

        function insertMacro(macroKey) {
            let editor;
            switch(currentEditor) {
                case 'header': editor = headerEditor; break;
                case 'body': editor = bodyEditor; break;
                case 'footer': editor = footerEditor; break;
            }

            const range = editor.getSelection(true);
            editor.insertText(range.index, macroKey);
        }

        function loadTemplate() {
            console.log('Carregando template ID:', templateId);

            fetch('/api/templates/' + templateId)
                .then(response => {
                    console.log('Response status:', response.status);
                    return response.json();
                })
                .then(data => {
                    console.log('Dados recebidos:', data);

                    if (data.success) {
                        const template = data.data;
                        console.log('Template:', template);

                        // Preencher campos de texto
                        document.getElementById('templateName').value = template.nome || '';
                        document.getElementById('templateDesc').value = template.descricao || '';
                        document.getElementById('templateSubject').value = template.assuntoPadrao || '';
                        document.getElementById('templateActive').checked = template.ativo;
                        document.getElementById('statusLabel').textContent = template.ativo ? 'Ativo' : 'Inativo';

                        // Carregar HTML nos editores Quill
                        // Limpar editores primeiro
                        headerEditor.setText('');
                        bodyEditor.setText('');
                        footerEditor.setText('');

                        // Carregar conteúdo usando root.innerHTML após limpar
                        if (template.headerHtml && template.headerHtml !== '') {
                            console.log('Carregando header:', template.headerHtml.substring(0, 100));
                            headerEditor.root.innerHTML = template.headerHtml;
                        }

                        if (template.bodyHtml && template.bodyHtml !== '') {
                            console.log('Carregando body:', template.bodyHtml.substring(0, 100));
                            bodyEditor.root.innerHTML = template.bodyHtml;
                        }

                        if (template.footerHtml && template.footerHtml !== '') {
                            console.log('Carregando footer:', template.footerHtml.substring(0, 100));
                            footerEditor.root.innerHTML = template.footerHtml;
                        }

                        console.log('Template carregado com sucesso');

                        // Atualizar estatísticas de tamanho após carregar
                        setTimeout(() => {
                            updateSizeStats();
                        }, 100);

                        showAlert('Template carregado com sucesso!', 'success');
                    } else {
                        console.error('Erro na resposta:', data.error);
                        showAlert('Erro ao carregar template: ' + data.error, 'error');
                    }
                })
                .catch(error => {
                    console.error('Erro ao carregar template:', error);
                    showAlert('Erro ao carregar template: ' + error.message, 'error');
                });
        }

        function saveTemplate() {
            const nome = document.getElementById('templateName').value.trim();
            const descricao = document.getElementById('templateDesc').value.trim();
            const assuntoPadrao = document.getElementById('templateSubject').value.trim();
            const ativo = document.getElementById('templateActive').checked;

            const headerHtml = headerEditor.root.innerHTML;
            const bodyHtml = bodyEditor.root.innerHTML;
            const footerHtml = footerEditor.root.innerHTML;

            // Validações
            if (!nome) {
                showAlert('Nome do template é obrigatório', 'error');
                return;
            }

            if (!bodyHtml || bodyHtml === '<p><br></p>') {
                showAlert('Corpo do template é obrigatório', 'error');
                return;
            }

            // Validação de tamanho
            const fullHtml = headerHtml + bodyHtml + footerHtml;
            const totalBytes = getByteSize(fullHtml);
            const maxBytes = 65000;

            // Avisar se exceder o limite
            if (totalBytes > maxBytes) {
                const htmlWithoutImages = removeBase64Images(fullHtml);
                const htmlWithoutImagesBytes = getByteSize(htmlWithoutImages);

                if (htmlWithoutImagesBytes > maxBytes) {
                    // Erro crítico: excede mesmo sem imagens
                    if (!confirm('⚠️ ATENÇÃO: O template excede ' + formatBytes(totalBytes) + ' KB (limite: 65 KB).\n\n' +
                        'Mesmo removendo as imagens base64, o conteúdo HTML excede o limite da API Zenvia.\n\n' +
                        'O envio via Zenvia FALHARÁ. Deseja salvar mesmo assim?')) {
                        return;
                    }
                } else {
                    // Aviso: excede mas as imagens serão removidas
                    if (!confirm('⚠️ AVISO: O template tem ' + formatBytes(totalBytes) + ' KB (limite: 65 KB).\n\n' +
                        'As imagens base64 serão REMOVIDAS AUTOMATICAMENTE ao enviar via Zenvia.\n\n' +
                        'Deseja continuar?')) {
                        return;
                    }
                }
            } else if (totalBytes > maxBytes * 0.8) {
                // Aviso suave: próximo do limite
                const imageStats = countBase64Images(fullHtml);
                if (imageStats.count > 0) {
                    showAlert('⚠️ Template próximo do limite (80%). Cuidado ao adicionar mais conteúdo.', 'warning');
                }
            }

            const payload = {
                nome: nome,
                descricao: descricao,
                headerHtml: headerHtml === '<p><br></p>' ? '' : headerHtml,
                bodyHtml: bodyHtml,
                footerHtml: footerHtml === '<p><br></p>' ? '' : footerHtml,
                assuntoPadrao: assuntoPadrao,
                ativo: ativo,
                criadoPor: 'sistema'
            };

            const url = templateId ? '/api/templates/' + templateId : '/api/templates';
            const method = templateId ? 'PUT' : 'POST';

            fetch(url, {
                method: method,
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(payload)
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    showAlert('Template salvo com sucesso!', 'success');
                    setTimeout(() => {
                        window.location.href = '/templates';
                    }, 1500);
                } else {
                    showAlert('Erro ao salvar: ' + (data.error || 'Erro desconhecido'), 'error');
                }
            })
            .catch(error => {
                console.error('Erro:', error);
                showAlert('Erro ao salvar template', 'error');
            });
        }

        function showPreview() {
            const headerHtml = headerEditor.root.innerHTML;
            const bodyHtml = bodyEditor.root.innerHTML;
            const footerHtml = footerEditor.root.innerHTML;

            fetch('/api/templates/preview', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    headerHtml: headerHtml === '<p><br></p>' ? '' : headerHtml,
                    bodyHtml: bodyHtml,
                    footerHtml: footerHtml === '<p><br></p>' ? '' : footerHtml,
                    useSampleData: true
                })
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    const previewWindow = window.open('', 'Preview', 'width=800,height=600');
                    previewWindow.document.write(data.html);
                    previewWindow.document.close();
                } else {
                    showAlert('Erro ao gerar preview: ' + data.error, 'error');
                }
            })
            .catch(error => {
                console.error('Erro:', error);
                showAlert('Erro ao gerar preview', 'error');
            });
        }

        // Função para calcular tamanho de string em bytes
        function getByteSize(str) {
            return new Blob([str]).size;
        }

        // Função para formatar bytes em KB
        function formatBytes(bytes) {
            return (bytes / 1024).toFixed(2);
        }

        // Função para contar e calcular tamanho de imagens base64
        function countBase64Images(html) {
            const regex = /<img[^>]*src="data:image\/[^"]*"[^>]*>/g;
            const matches = html.match(regex);
            if (!matches) return { count: 0, size: 0 };

            let totalSize = 0;
            matches.forEach(match => {
                totalSize += getByteSize(match);
            });

            return { count: matches.length, size: totalSize };
        }

        // Função para remover imagens base64 do HTML
        function removeBase64Images(html) {
            const regex = /<img[^>]*src="data:image\/[^"]*"[^>]*>/g;
            return html.replace(regex, '');
        }

        // Função para atualizar estatísticas de tamanho
        function updateSizeStats() {
            const headerHtml = headerEditor.root.innerHTML;
            const bodyHtml = bodyEditor.root.innerHTML;
            const footerHtml = footerEditor.root.innerHTML;

            const fullHtml = headerHtml + bodyHtml + footerHtml;
            const totalBytes = getByteSize(fullHtml);
            const imageStats = countBase64Images(fullHtml);

            const htmlWithoutImages = removeBase64Images(fullHtml);
            const htmlWithoutImagesBytes = getByteSize(htmlWithoutImages);

            const maxBytes = 65000; // Limite da Zenvia
            const totalKB = formatBytes(totalBytes);
            const imageKB = formatBytes(imageStats.size);
            const percentage = (totalBytes / maxBytes) * 100;
            const imagePercentage = (imageStats.size / maxBytes) * 100;

            // Atualizar valores
            document.getElementById('totalSize').textContent = totalKB + ' KB';
            document.getElementById('imageSize').textContent = imageKB + ' KB';
            if (imageStats.count > 0) {
                document.getElementById('imageSize').textContent += ' (' + imageStats.count + ' img)';
            }

            // Atualizar barra de progresso total
            const totalBar = document.getElementById('totalSizeBar');
            totalBar.style.width = Math.min(percentage, 100) + '%';

            // Atualizar barra de progresso de imagens
            const imageBar = document.getElementById('imageSizeBar');
            imageBar.style.width = Math.min(imagePercentage, 100) + '%';

            // Atualizar classes de cor baseado no tamanho
            const totalValue = document.getElementById('totalSize');
            const imageValue = document.getElementById('imageSize');

            // Limpar classes
            totalValue.classList.remove('ok', 'warning', 'danger');
            totalBar.classList.remove('ok', 'warning', 'danger');
            imageValue.classList.remove('ok', 'warning', 'danger');
            imageBar.classList.remove('ok', 'warning', 'danger');

            // Definir cores baseado no tamanho total
            if (totalBytes > maxBytes) {
                // Excede o limite
                totalValue.classList.add('danger');
                totalBar.classList.add('danger');
            } else if (totalBytes > maxBytes * 0.8) {
                // 80% ou mais do limite
                totalValue.classList.add('warning');
                totalBar.classList.add('warning');
            } else {
                // OK
                totalValue.classList.add('ok');
                totalBar.classList.add('ok');
            }

            // Definir cores para imagens
            if (imageStats.size > maxBytes * 0.5) {
                imageValue.classList.add('danger');
                imageBar.classList.add('danger');
            } else if (imageStats.size > maxBytes * 0.3) {
                imageValue.classList.add('warning');
                imageBar.classList.add('warning');
            } else {
                imageValue.classList.add('ok');
                imageBar.classList.add('ok');
            }

            // Mostrar/esconder avisos
            const warningDiv = document.getElementById('sizeWarning');
            const errorDiv = document.getElementById('sizeError');

            warningDiv.classList.remove('show');
            errorDiv.classList.remove('show');

            if (totalBytes > maxBytes) {
                // Verificar se mesmo sem imagens excede o limite
                if (htmlWithoutImagesBytes > maxBytes) {
                    errorDiv.classList.add('show');
                } else {
                    warningDiv.classList.add('show');
                }
            } else if (totalBytes > maxBytes * 0.8) {
                // Próximo do limite
                if (imageStats.count > 0) {
                    warningDiv.classList.add('show');
                }
            }
        }

        function showAlert(message, type) {
            const container = document.getElementById('alertContainer');
            const alert = document.createElement('div');
            alert.className = 'alert alert-' + type;
            alert.textContent = message;
            container.appendChild(alert);

            setTimeout(() => {
                alert.remove();
            }, 5000);
        }
    </script>
</body>
</html>
`
