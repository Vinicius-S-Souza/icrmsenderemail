package service

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/kardianos/service"
	"go.uber.org/zap"
)

// AppService representa o serviço da aplicação
type AppService struct {
	Logger      *zap.Logger
	RunFunc     func(context.Context) error
	ctx         context.Context
	cancel      context.CancelFunc
	serviceName string
	displayName string
	description string
}

// Config contém as configurações do serviço
type Config struct {
	Name        string
	DisplayName string
	Description string
	Logger      *zap.Logger
	RunFunc     func(context.Context) error
}

// NewAppService cria uma nova instância do serviço
func NewAppService(cfg Config) *AppService {
	return &AppService{
		Logger:      cfg.Logger,
		RunFunc:     cfg.RunFunc,
		serviceName: cfg.Name,
		displayName: cfg.DisplayName,
		description: cfg.Description,
	}
}

// Start implementa service.Interface
func (app *AppService) Start(s service.Service) error {
	// Cria contexto imediatamente
	app.ctx, app.cancel = context.WithCancel(context.Background())

	// Log pode falhar em modo serviço, então usamos defer para garantir que sempre retorna nil
	defer func() {
		if r := recover(); r != nil {
			// Se houver panic no log, ignora e continua
		}
	}()

	if app.Logger != nil {
		app.Logger.Info("Iniciando serviço", zap.String("name", app.serviceName))
	}

	// Inicia o serviço em uma goroutine imediatamente
	go app.run()

	// Retorna imediatamente para o Windows saber que o serviço iniciou
	return nil
}

// run executa a lógica principal do serviço
func (app *AppService) run() {
	// Cria um log de debug para troubleshooting em modo serviço
	// Este arquivo sempre é criado, independente da configuração do logger
	debugLog, err := os.OpenFile("service_debug.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err == nil {
		defer debugLog.Close()
		debugLog.WriteString(fmt.Sprintf("[%s] Serviço iniciado - goroutine executando\n", time.Now().Format("2006-01-02 15:04:05")))
		debugLog.WriteString(fmt.Sprintf("[%s] Diretório de trabalho: %s\n", time.Now().Format("2006-01-02 15:04:05"), func() string {
			dir, _ := os.Getwd()
			return dir
		}()))
	}

	// Protege contra panics que possam ocorrer durante a execução
	defer func() {
		if r := recover(); r != nil {
			if debugLog != nil {
				debugLog.WriteString(fmt.Sprintf("[%s] PANIC: %v\n", time.Now().Format("2006-01-02 15:04:05"), r))
			}
			if app.Logger != nil {
				app.Logger.Error("Panic durante execução do serviço",
					zap.Any("panic", r))
			}
		}
	}()

	if app.Logger != nil {
		app.Logger.Info("Serviço em execução", zap.String("name", app.serviceName))
	}

	if debugLog != nil {
		debugLog.WriteString(fmt.Sprintf("[%s] Chamando RunFunc...\n", time.Now().Format("2006-01-02 15:04:05")))
	}

	if err := app.RunFunc(app.ctx); err != nil && err != context.Canceled {
		if debugLog != nil {
			debugLog.WriteString(fmt.Sprintf("[%s] ERRO: %v\n", time.Now().Format("2006-01-02 15:04:05"), err))
		}
		if app.Logger != nil {
			app.Logger.Error("Erro na execução do serviço", zap.Error(err))
		}
	}

	if debugLog != nil {
		debugLog.WriteString(fmt.Sprintf("[%s] RunFunc finalizado\n", time.Now().Format("2006-01-02 15:04:05")))
	}
}

// Stop implementa service.Interface
func (app *AppService) Stop(s service.Service) error {
	app.Logger.Info("Parando serviço", zap.String("name", app.serviceName))

	if app.cancel != nil {
		app.cancel()
	}

	return nil
}

// IsWindowsService verifica se está rodando como serviço do Windows
func IsWindowsService() (bool, error) {
	if runtime.GOOS != "windows" {
		return false, nil
	}

	// Usa a biblioteca kardianos/service para detectar
	// service.Interactive() retorna true se está em modo interativo (não serviço)
	// então invertemos o resultado
	isInteractive := service.Interactive()
	return !isInteractive, nil
}

// Run executa o serviço com base nos argumentos da linha de comando
func Run(cfg Config, arguments []string) error {
	svcConfig := &service.Config{
		Name:        cfg.Name,
		DisplayName: cfg.DisplayName,
		Description: cfg.Description,
	}

	appService := NewAppService(cfg)
	s, err := service.New(appService, svcConfig)
	if err != nil {
		return fmt.Errorf("erro ao criar serviço: %w", err)
	}

	// Se não há argumentos, executa o serviço normalmente
	if len(arguments) == 0 {
		cfg.Logger.Info("Executando aplicação", zap.String("mode", "normal"))
		return s.Run()
	}

	// Processa comandos de controle do serviço
	command := arguments[0]

	switch command {
	case "install":
		err = s.Install()
		if err != nil {
			return fmt.Errorf("erro ao instalar serviço: %w", err)
		}
		cfg.Logger.Info("Serviço instalado com sucesso", zap.String("name", cfg.Name))
		return nil

	case "uninstall":
		err = s.Uninstall()
		if err != nil {
			return fmt.Errorf("erro ao desinstalar serviço: %w", err)
		}
		cfg.Logger.Info("Serviço desinstalado com sucesso", zap.String("name", cfg.Name))
		return nil

	case "start":
		err = s.Start()
		if err != nil {
			return fmt.Errorf("erro ao iniciar serviço: %w", err)
		}
		cfg.Logger.Info("Serviço iniciado com sucesso", zap.String("name", cfg.Name))
		return nil

	case "stop":
		err = s.Stop()
		if err != nil {
			return fmt.Errorf("erro ao parar serviço: %w", err)
		}
		cfg.Logger.Info("Serviço parado com sucesso", zap.String("name", cfg.Name))
		return nil

	case "restart":
		err = s.Restart()
		if err != nil {
			return fmt.Errorf("erro ao reiniciar serviço: %w", err)
		}
		cfg.Logger.Info("Serviço reiniciado com sucesso", zap.String("name", cfg.Name))
		return nil

	default:
		return fmt.Errorf("comando desconhecido: %s (use: install, uninstall, start, stop, restart)", command)
	}
}
