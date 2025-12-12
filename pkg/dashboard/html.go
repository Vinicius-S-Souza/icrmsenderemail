package dashboard

// dashboardHTML contém o HTML da página principal do dashboard
// Data de criação: 11/12/2025
const dashboardHTML = `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ICRMSenderEmail - Dashboard</title>
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

        .header-left {
            flex: 1;
        }

        .header-actions {
            display: flex;
            gap: 10px;
            flex-wrap: wrap;
        }

        .manual-btn, .templates-btn {
            display: inline-block;
            padding: 12px 24px;
            color: white;
            text-decoration: none;
            border-radius: 8px;
            font-weight: 600;
            transition: all 0.3s;
            box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
        }

        .manual-btn {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        }

        .templates-btn {
            background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
            box-shadow: 0 4px 12px rgba(245, 158, 11, 0.3);
        }

        .manual-btn:hover, .templates-btn:hover {
            transform: translateY(-2px);
        }

        .manual-btn:hover {
            box-shadow: 0 6px 16px rgba(102, 126, 234, 0.4);
        }

        .templates-btn:hover {
            box-shadow: 0 6px 16px rgba(245, 158, 11, 0.4);
        }

        h1 {
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

        .grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }

        .card {
            background: white;
            padding: 25px;
            border-radius: 15px;
            box-shadow: 0 5px 20px rgba(0,0,0,0.1);
            transition: transform 0.3s ease, box-shadow 0.3s ease;
        }

        .card:hover {
            transform: translateY(-5px);
            box-shadow: 0 10px 30px rgba(0,0,0,0.15);
        }

        .card-title {
            font-size: 0.9em;
            color: #666;
            text-transform: uppercase;
            letter-spacing: 1px;
            margin-bottom: 10px;
        }

        .card-value {
            font-size: 2.5em;
            font-weight: bold;
            color: #667eea;
            margin-bottom: 10px;
        }

        .card-subtitle {
            font-size: 0.85em;
            color: #999;
        }

        .success { color: #10b981; }
        .error { color: #ef4444; }
        .warning { color: #f59e0b; }
        .info { color: #3b82f6; }

        .chart-container {
            background: white;
            padding: 30px;
            border-radius: 15px;
            box-shadow: 0 5px 20px rgba(0,0,0,0.1);
            margin-bottom: 30px;
        }

        .chart-title {
            font-size: 1.3em;
            color: #333;
            margin-bottom: 20px;
            font-weight: 600;
        }

        canvas {
            max-height: 300px;
        }

        .progress-bar {
            width: 100%;
            height: 8px;
            background: #e5e7eb;
            border-radius: 10px;
            overflow: hidden;
            margin-top: 10px;
        }

        .progress-fill {
            height: 100%;
            background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
            border-radius: 10px;
            transition: width 0.3s ease;
        }

        .timestamp {
            text-align: center;
            color: white;
            margin-top: 20px;
            font-size: 0.9em;
            opacity: 0.8;
        }

        .metrics-row {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
            margin-bottom: 20px;
        }

        .metric-item {
            padding: 15px;
            background: #f9fafb;
            border-radius: 10px;
            border-left: 4px solid #667eea;
        }

        .metric-label {
            font-size: 0.85em;
            color: #666;
            margin-bottom: 5px;
        }

        .metric-value {
            font-size: 1.5em;
            font-weight: bold;
            color: #333;
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <div class="header-left">
                <h1>📧 ICRMSenderEmail</h1>
                <span class="subtitle">Dashboard de Métricas em Tempo Real</span>
                <span class="subtitle"> | Provedor: <strong id="provider-name">Carregando...</strong></span>
                <span class="status" id="connection-status">● Conectado</span>
            </div>
            <div class="header-actions">
                <a href="/templates" class="templates-btn">📝 Templates</a>
                <a href="/manual" class="manual-btn">📨 Disparo Manual</a>
            </div>
        </header>

        <div class="grid">
            <div class="card">
                <div class="card-title">Total Processado</div>
                <div class="card-value" id="total-messages">0</div>
                <div class="card-subtitle">E-mails enviados</div>
                <div class="progress-bar">
                    <div class="progress-fill" style="width: 100%"></div>
                </div>
            </div>

            <div class="card">
                <div class="card-title info">📬 Pendentes</div>
                <div class="card-value info" id="pending-count">0</div>
                <div class="card-subtitle">Mensagens aguardando envio</div>
            </div>

            <div class="card">
                <div class="card-title success">✓ Sucesso</div>
                <div class="card-value success" id="success-count">0</div>
                <div class="card-subtitle" id="success-rate">0% taxa de sucesso</div>
            </div>

            <div class="card">
                <div class="card-title error">✗ Erros</div>
                <div class="card-value error" id="error-count">0</div>
                <div class="card-subtitle" id="error-rate">0% taxa de erro</div>
            </div>

            <div class="card">
                <div class="card-title warning">⚠ E-mails Inválidos</div>
                <div class="card-value warning" id="invalid-count">0</div>
                <div class="card-subtitle" id="invalid-rate">0% taxa inválida</div>
            </div>
        </div>

        <div class="chart-container">
            <div class="chart-title">📊 Taxa de Sucesso vs Erros (Tempo Real)</div>
            <canvas id="successChart"></canvas>
        </div>

        <div class="chart-container">
            <div class="chart-title">⏱️ Tempos Médios de Processamento (ms)</div>
            <div class="metrics-row">
                <div class="metric-item">
                    <div class="metric-label">Tempo de Processamento</div>
                    <div class="metric-value" id="avg-process-time">0ms</div>
                </div>
                <div class="metric-item">
                    <div class="metric-label">Tempo de Envio Email</div>
                    <div class="metric-value" id="avg-send-time">0ms</div>
                </div>
                <div class="metric-item">
                    <div class="metric-label">Tempo de Query</div>
                    <div class="metric-value" id="avg-query-time">0ms</div>
                </div>
            </div>
            <canvas id="timeChart"></canvas>
        </div>

        <div class="chart-container">
            <div class="chart-title">📈 Métricas de Envio de E-mail</div>
            <div class="metrics-row">
                <div class="metric-item">
                    <div class="metric-label">E-mails Enviados com Sucesso</div>
                    <div class="metric-value success" id="email-success-count">0</div>
                </div>
                <div class="metric-item">
                    <div class="metric-label">E-mails com Erro</div>
                    <div class="metric-value error" id="email-error-count">0</div>
                </div>
                <div class="metric-item">
                    <div class="metric-label">Taxa de Sucesso E-mail</div>
                    <div class="metric-value info" id="email-success-rate">0%</div>
                </div>
                <div class="metric-item">
                    <div class="metric-label">Queries Executadas</div>
                    <div class="metric-value info" id="queries-executed">0</div>
                </div>
            </div>
        </div>

        <div class="timestamp" id="last-update">Última atualização: --</div>
    </div>

    <script src="https://cdn.jsdelivr.net/npm/chart.js@4.4.0/dist/chart.umd.min.js"></script>
    <script>
        // Configuração dos gráficos
        const successChartCtx = document.getElementById('successChart').getContext('2d');
        const timeChartCtx = document.getElementById('timeChart').getContext('2d');

        const maxDataPoints = 30;
        const chartData = {
            labels: [],
            success: [],
            errors: [],
            invalidEmails: [],
            processTime: [],
            sendTime: [],
            queryTime: []
        };

        const successChart = new Chart(successChartCtx, {
            type: 'line',
            data: {
                labels: chartData.labels,
                datasets: [
                    {
                        label: 'Sucesso',
                        data: chartData.success,
                        borderColor: '#10b981',
                        backgroundColor: 'rgba(16, 185, 129, 0.1)',
                        tension: 0.4,
                        fill: true
                    },
                    {
                        label: 'Erros',
                        data: chartData.errors,
                        borderColor: '#ef4444',
                        backgroundColor: 'rgba(239, 68, 68, 0.1)',
                        tension: 0.4,
                        fill: true
                    },
                    {
                        label: 'E-mails Inválidos',
                        data: chartData.invalidEmails,
                        borderColor: '#f59e0b',
                        backgroundColor: 'rgba(245, 158, 11, 0.1)',
                        tension: 0.4,
                        fill: true
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: true,
                plugins: {
                    legend: {
                        display: true,
                        position: 'top'
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true
                    }
                }
            }
        });

        const timeChart = new Chart(timeChartCtx, {
            type: 'bar',
            data: {
                labels: chartData.labels,
                datasets: [
                    {
                        label: 'Processamento (ms)',
                        data: chartData.processTime,
                        backgroundColor: 'rgba(102, 126, 234, 0.8)'
                    },
                    {
                        label: 'Envio Email (ms)',
                        data: chartData.sendTime,
                        backgroundColor: 'rgba(118, 75, 162, 0.8)'
                    },
                    {
                        label: 'Query (ms)',
                        data: chartData.queryTime,
                        backgroundColor: 'rgba(59, 130, 246, 0.8)'
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: true,
                plugins: {
                    legend: {
                        display: true,
                        position: 'top'
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true
                    }
                }
            }
        });

        // Conectar ao stream de métricas
        const eventSource = new EventSource('/api/metrics/stream');

        eventSource.onmessage = function(event) {
            const metrics = JSON.parse(event.data);
            updateDashboard(metrics);
        };

        eventSource.onerror = function() {
            document.getElementById('connection-status').textContent = '● Desconectado';
            document.getElementById('connection-status').style.background = '#ef4444';
        };

        eventSource.onopen = function() {
            document.getElementById('connection-status').textContent = '● Conectado';
            document.getElementById('connection-status').style.background = '#10b981';
        };

        function updateDashboard(metrics) {
            // Atualizar nome do provider
            if (metrics.provider_name) {
                const providerNames = {
                    'mock': '🧪 Mock (Teste)',
                    'smtp': '📧 SMTP',
                    'sendgrid': '📨 SendGrid',
                    'zenvia': '🇧🇷 Zenvia',
                    'pontaltech': '📡 Pontaltech'
                };
                const displayName = providerNames[metrics.provider_name] || metrics.provider_name.toUpperCase();
                document.getElementById('provider-name').textContent = displayName;
            }

            // Atualizar cards principais
            document.getElementById('total-messages').textContent = metrics.total_messages_processed.toLocaleString();
            document.getElementById('pending-count').textContent = metrics.pending_messages_count.toLocaleString();
            document.getElementById('success-count').textContent = metrics.success_count.toLocaleString();
            document.getElementById('error-count').textContent = metrics.error_count.toLocaleString();
            document.getElementById('invalid-count').textContent = metrics.invalid_email_count.toLocaleString();

            document.getElementById('success-rate').textContent = metrics.success_rate.toFixed(1) + '% taxa de sucesso';
            document.getElementById('error-rate').textContent = metrics.error_rate.toFixed(1) + '% taxa de erro';
            document.getElementById('invalid-rate').textContent = metrics.invalid_email_rate.toFixed(1) + '% taxa inválida';

            // Atualizar tempos médios
            document.getElementById('avg-process-time').textContent = metrics.avg_process_time_ms.toFixed(2) + 'ms';
            document.getElementById('avg-send-time').textContent = metrics.avg_send_time_ms.toFixed(2) + 'ms';
            document.getElementById('avg-query-time').textContent = metrics.avg_query_time_ms.toFixed(2) + 'ms';

            // Atualizar métricas de email
            document.getElementById('email-success-count').textContent = metrics.email_send_success_count.toLocaleString();
            document.getElementById('email-error-count').textContent = metrics.email_send_error_count.toLocaleString();
            document.getElementById('email-success-rate').textContent = metrics.email_send_success_rate.toFixed(1) + '%';
            document.getElementById('queries-executed').textContent = metrics.queries_executed.toLocaleString();

            // Atualizar timestamp
            const timestamp = new Date(metrics.timestamp);
            document.getElementById('last-update').textContent =
                'Última atualização: ' + timestamp.toLocaleTimeString('pt-BR');

            // Atualizar gráficos
            const timeLabel = timestamp.toLocaleTimeString('pt-BR');

            chartData.labels.push(timeLabel);
            chartData.success.push(metrics.success_count);
            chartData.errors.push(metrics.error_count);
            chartData.invalidEmails.push(metrics.invalid_email_count);
            chartData.processTime.push(metrics.avg_process_time_ms);
            chartData.sendTime.push(metrics.avg_send_time_ms);
            chartData.queryTime.push(metrics.avg_query_time_ms);

            // Limitar número de pontos no gráfico
            if (chartData.labels.length > maxDataPoints) {
                chartData.labels.shift();
                chartData.success.shift();
                chartData.errors.shift();
                chartData.invalidEmails.shift();
                chartData.processTime.shift();
                chartData.sendTime.shift();
                chartData.queryTime.shift();
            }

            successChart.update('none');
            timeChart.update('none');
        }

        // Carregar métricas iniciais
        fetch('/api/metrics')
            .then(response => response.json())
            .then(metrics => updateDashboard(metrics))
            .catch(error => console.error('Erro ao carregar métricas:', error));
    </script>
</body>
</html>
`
