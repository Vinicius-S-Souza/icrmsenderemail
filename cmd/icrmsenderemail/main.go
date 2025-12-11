package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/cliente"
	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/config"
	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/control"
	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/dashboard"
	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/database"
	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/email"
	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/health"
	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/logger"
	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/manual"
	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/message"
	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/metrics"
	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/service"
	"github.com/Vinicius-S-Souza/icrmsenderemail/pkg/version"
	"go.uber.org/zap"
)

func main() {
	// Inicializar logger básico (funciona tanto para serviço quanto para execução normal)
	log := logger.CreateLogger()

	// Configuração do serviço
	svcConfig := service.Config{
		Name:        "icrmsenderemail",
		DisplayName: "iCRM Sender Email",
		Description: "Serviço de envio de e-mail usando SMTP, SendGrid, Zenvia, Pontaltech ou mock",
		Logger:      log,
		RunFunc:     runApplication,
	}

	// Verifica se há argumentos de linha de comando
	if len(os.Args) > 1 {
		serviceCommand := os.Args[1]

		// Comando --version ou -v: Exibe versão e sai
		if serviceCommand == "--version" || serviceCommand == "-v" || serviceCommand == "version" {
			version.PrintVersion()
			return
		}

		// Comandos de controle de serviço (install, uninstall, start, stop, restart)
		if serviceCommand == "install" || serviceCommand == "uninstall" ||
			serviceCommand == "start" || serviceCommand == "stop" || serviceCommand == "restart" {

			if err := service.Run(svcConfig, os.Args[1:]); err != nil {
				log.Fatal("Erro ao executar comando de serviço",
					zap.String("command", serviceCommand),
					zap.Error(err))
			}
			return
		}

		// Comando desconhecido
		fmt.Printf("Comando desconhecido: %s\n", serviceCommand)
		fmt.Println("\nUso: icrmsenderemail [comando]")
		fmt.Println("\nComandos disponíveis:")
		fmt.Println("  install    - Instala o serviço")
		fmt.Println("  uninstall  - Desinstala o serviço")
		fmt.Println("  start      - Inicia o serviço")
		fmt.Println("  stop       - Para o serviço")
		fmt.Println("  restart    - Reinicia o serviço")
		fmt.Println("  version    - Exibe informações de versão")
		fmt.Println("  -v         - Exibe informações de versão")
		fmt.Println("  --version  - Exibe informações de versão")
		fmt.Println("\nSem argumentos: executa a aplicação em modo normal (foreground)")
		return
	}

	// Tenta executar como serviço (se foi iniciado pelo Windows Service Manager)
	// Se falhar, executa normalmente
	isService, err := service.IsWindowsService()
	if err == nil && isService {
		// Rodando como serviço do Windows
		if err := service.Run(svcConfig, []string{}); err != nil {
			log.Fatal("Erro ao executar como serviço", zap.Error(err))
		}
		return
	}

	// Execução normal (modo interativo)
	if err := runApplication(context.Background()); err != nil {
		fmt.Printf("Erro na execução: %v\n", err)
		os.Exit(1)
	}
}

// getBaseDir retorna o diretório base da aplicação
// Para "go run": retorna o diretório de trabalho atual (raiz do projeto)
// Para executável: retorna o diretório onde o executável está
func getBaseDir() (string, bool, error) {
	// Obter diretório de trabalho atual
	workDir, err := os.Getwd()
	if err != nil {
		return "", false, fmt.Errorf("erro ao obter diretório de trabalho: %w", err)
	}

	// Obter diretório do executável
	exePath, err := os.Executable()
	if err != nil {
		return "", false, fmt.Errorf("erro ao obter caminho do executável: %w", err)
	}
	exeDir := filepath.Dir(exePath)

	// Verificar se está executando via "go run" (executável no cache do Go)
	isGoRun := strings.Contains(exeDir, "go-build") || strings.Contains(exePath, "go-build")

	if isGoRun {
		// Quando executado via "go run", usar o diretório de trabalho atual
		return workDir, true, nil
	}

	// Quando executado via binário, usar o diretório do executável
	return exeDir, false, nil
}

// findConfigFile procura o arquivo dbinit.ini no diretório base apropriado
func findConfigFile() (string, error) {
	configFileName := "dbinit.ini"

	baseDir, isGoRun, err := getBaseDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(baseDir, configFileName)
	if _, err := os.Stat(configPath); err == nil {
		return configPath, nil
	}

	runMode := "executável"
	if isGoRun {
		runMode = "go run"
	}

	return "", fmt.Errorf("arquivo %s não encontrado no diretório base (%s) [modo: %s]",
		configFileName, baseDir, runMode)
}

// runApplication é a lógica principal da aplicação
func runApplication(ctx context.Context) error {
	// Mensagem obrigatória de início (sempre no stdout)
	fmt.Println("[ICRMSENDEREMAIL] Aplicação iniciando...")

	// Obter diretório do executável
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("erro ao obter caminho do executável: %w", err)
	}
	exeDir := filepath.Dir(exePath)

	// Obter diretório de trabalho
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("erro ao obter diretório de trabalho: %w", err)
	}

	fmt.Printf("Diretório do executável: %s\n", exeDir)
	fmt.Printf("Diretório de trabalho: %s\n", workDir)
	fmt.Printf("%s\n", version.GetVersion())
	fmt.Println("========================================")

	// Localizar arquivo de configuração
	configPath, err := findConfigFile()
	if err != nil {
		return fmt.Errorf("erro ao localizar configuração: %w", err)
	}

	fmt.Printf("Usando arquivo de configuração: %s\n", configPath)

	// Obter diretório base da aplicação
	baseDir, isGoRun, err := getBaseDir()
	if err != nil {
		return fmt.Errorf("erro ao obter diretório base: %w", err)
	}

	// Mudar para o diretório base para que caminhos relativos funcionem
	if err := os.Chdir(baseDir); err != nil {
		return fmt.Errorf("erro ao mudar para diretório base: %w", err)
	}

	runMode := "executável"
	if isGoRun {
		runMode = "go run"
	}
	fmt.Printf("Modo de execução: %s\n", runMode)
	fmt.Printf("Diretório base: %s\n", baseDir)

	// Carregar configurações
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("erro ao carregar configurações: %w", err)
	}

	// Validar configurações
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("configuração inválida: %w", err)
	}

	// Inicializar logger (LogDir será relativo ao diretório base)
	logConfig := &logger.LogConfig{
		LogDir:        cfg.Logger.LogDir,
		ConsoleOutput: cfg.Logger.ConsoleOutput,
		LogLevel:      cfg.Logger.LogLevel,
		RetentionDays: cfg.Logger.RetentionDays,
	}

	if err := logger.InitLogger(logConfig); err != nil {
		return fmt.Errorf("erro ao inicializar logger: %w", err)
	}

	log := logger.GetLogger()
	log.Info("Aplicação iniciada",
		zap.String("versao", version.Version),
		zap.String("build_date", version.BuildDate))

	// Conectar ao banco de dados
	log.Info("Conectando ao banco de dados Oracle...")
	db, err := database.ConnectOracle(database.DBConfig{
		Username: cfg.Database.Username,
		Password: cfg.Database.Password,
		TNS:      cfg.Database.TNS,
	})
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados", zap.Error(err))
	}
	defer db.Close()

	log.Info("Conexão com banco de dados estabelecida")

	// Inicializar provedor de e-mail
	log.Info("Inicializando provedor de e-mail",
		zap.String("provider", cfg.Email.Provider))

	var provider email.Provider
	switch cfg.Email.Provider {
	case "mock":
		provider = email.NewMockProvider(log)
	case "smtp":
		provider = email.NewSMTPProvider(
			cfg.Email.SMTPHost,
			cfg.Email.SMTPPort,
			cfg.Email.SMTPUsername,
			cfg.Email.SMTPPassword,
			cfg.Email.SMTPUseTLS,
			log,
		)
	case "sendgrid":
		provider = email.NewSendGridProvider(
			cfg.Email.SendGridAPIKey,
			log,
		)
	case "zenvia":
		provider = email.NewZenviaProvider(
			cfg.Email.ZenviaAPIToken,
			log,
		)
	case "pontaltech":
		provider = email.NewPontaltechProvider(
			cfg.Email.PontaltechUsername,
			cfg.Email.PontaltechPassword,
			cfg.Email.PontaltechAccountID,
			cfg.Email.PontaltechAPIURL,
			cfg.Email.PontaltechCallbackURL,
			log,
		)
	default:
		log.Fatal("Provedor de e-mail não suportado",
			zap.String("provider", cfg.Email.Provider))
	}

	// Criar componentes
	repo := message.NewRepository(db, log)
	sender := email.NewSender(provider, log)
	metricsCollector := metrics.NewPerformanceMetrics()

	// Criar processador
	processor := message.NewProcessor(
		repo,
		sender,
		metricsCollector,
		&cfg.Performance,
		cfg.Email.DefaultFrom,
		log,
	)

	// Iniciar health check server
	if cfg.Health.Enabled {
		healthChecker := health.NewHealthChecker(db, log)
		health.StartHealthServer(cfg.Health.HTTPPort, healthChecker, log)
	}

	// Iniciar dashboard
	var dashboardServer *dashboard.Dashboard
	if cfg.Dashboard.EnableDashboard {
		dashboardConfig := dashboard.Config{
			Port:            cfg.Dashboard.DashboardPort,
			EnableDashboard: true,
			ProviderName:    cfg.Email.Provider,
			DaysOffset:      cfg.Performance.DataDisparoOffset,
			MaxTentativas:   cfg.Performance.MaxTentativas,
		}
		dashboardServer = dashboard.NewDashboard(dashboardConfig, metricsCollector, repo, log)

		// Registrar endpoints de disparo manual
		clienteRepo := cliente.NewRepository(db, log)
		manualHandler := manual.NewHandler(clienteRepo, repo, cfg.Email.Provider)
		dashboardServer.RegisterManualEndpoints(manualHandler)

		go func() {
			if err := dashboardServer.Start(); err != nil && err != http.ErrServerClosed {
				log.Error("Erro no dashboard", zap.Error(err))
			}
		}()

		log.Info("Dashboard e endpoints manuais iniciados",
			zap.Int("port", cfg.Dashboard.DashboardPort))
	}

	// Monitorar sinais do sistema (apenas se não estiver em modo serviço do Windows)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Criar contexto interno que pode ser cancelado por sinais OU pelo contexto do serviço
	appCtx, appCancel := context.WithCancel(ctx)
	defer appCancel()

	// Iniciar monitoramento de arquivo de controle
	go control.WatchStopFile(appCtx, log, appCancel)

	// Iniciar processador
	if err := processor.Start(); err != nil {
		log.Fatal("Erro ao iniciar processador", zap.Error(err))
	}

	log.Info("Serviço iniciado com sucesso",
		zap.Int("workers", cfg.Performance.WorkerCount),
		zap.Int("batch_size", cfg.Performance.BatchSize),
		zap.Bool("dashboard", cfg.Dashboard.EnableDashboard),
		zap.Int("dashboard_port", cfg.Dashboard.DashboardPort))

	// Exibir estatísticas periodicamente
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-appCtx.Done():
				return
			case <-ticker.C:
				metricsCollector.LogMetrics(log)
				printDBStats(log, repo)
			}
		}
	}()

	// Resetar estatísticas na virada do dia
	go func() {
		currentDay := time.Now().Day()
		ticker := time.NewTicker(1 * time.Minute) // Verifica a cada 1 minuto
		defer ticker.Stop()

		for {
			select {
			case <-appCtx.Done():
				return
			case <-ticker.C:
				now := time.Now()
				if now.Day() != currentDay {
					log.Info("🔄 Virada do dia detectada, resetando estatísticas...",
						zap.Int("dia_anterior", currentDay),
						zap.Int("dia_atual", now.Day()))

					// Log das estatísticas antes do reset
					metricsCollector.LogMetrics(log)

					// Reset das métricas
					metricsCollector.Reset()

					// Atualizar o dia atual
					currentDay = now.Day()

					log.Info("✅ Estatísticas resetadas com sucesso para o novo dia")
				}
			}
		}
	}()

	// Aguardar sinal de parada
	select {
	case sig := <-sigChan:
		log.Info("Sinal recebido, iniciando shutdown...",
			zap.String("signal", sig.String()))
	case <-appCtx.Done():
		log.Info("Contexto cancelado, iniciando shutdown...")
	}

	// Graceful shutdown
	log.Info("Parando processador...")
	if err := processor.Stop(); err != nil {
		log.Error("Erro ao parar processador", zap.Error(err))
	}

	// Parar dashboard
	if dashboardServer != nil {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := dashboardServer.Stop(shutdownCtx); err != nil {
			log.Error("Erro ao parar dashboard", zap.Error(err))
		}
	}

	// Exibir estatísticas finais
	metricsCollector.LogMetrics(log)
	printDBStats(log, repo)

	log.Info("Aplicação finalizada com sucesso")
	logger.GetLogger().Sync()

	return nil
}

func printDBStats(log *zap.Logger, repo *message.Repository) {
	// Estatísticas do banco
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if dbStats, err := repo.GetStats(ctx); err == nil {
		log.Info("📊 Banco de dados (hoje)",
			zap.Int64("status_0_pendentes", dbStats["status_0"]),
			zap.Int64("status_2_enviados", dbStats["status_2"]),
			zap.Int64("status_3_erros", dbStats["status_3"]),
			zap.Int64("status_4_falhas_permanentes", dbStats["status_4"]),
			zap.Int64("status_125_invalidos", dbStats["status_125"]))
	}
}
