package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-ini/ini"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConfig struct {
	LogDir        string
	ConsoleOutput bool
	LogLevel      string
	RetentionDays int
}

var logger *zap.Logger
var verboseMode bool
var currentLogFile *os.File
var currentConfig *LogConfig
var currentDate string

// SetVerboseMode define o modo verbose globalmente
func SetVerboseMode(enabled bool) {
	verboseMode = enabled
}

// IsVerbose retorna se o modo verbose está habilitado
func IsVerbose() bool {
	return verboseMode
}

func getLogFileName() string {
	return fmt.Sprintf("icrmsenderemail_%s.log", time.Now().Format("20060102"))
}

func CreateLogger() *zap.Logger {
	// Configuração padrão se o logger não foi inicializado
	config := &LogConfig{
		LogDir:        "log",
		ConsoleOutput: true,
		LogLevel:      "info",
		RetentionDays: 30,
	}

	// Tentar carregar configurações do arquivo
	if cfg, err := ini.Load("dbinit.ini"); err == nil {
		loggerSection := cfg.Section("logger")
		config = &LogConfig{
			LogDir:        loggerSection.Key("log_dir").MustString("log"),
			ConsoleOutput: loggerSection.Key("console_output").MustBool(true),
			LogLevel:      loggerSection.Key("log_level").MustString("info"),
			RetentionDays: loggerSection.Key("retention_days").MustInt(30),
		}
	}

	if err := InitLogger(config); err != nil {
		// Em caso de erro, criar um logger básico para console
		logger, _ = zap.NewProduction()
	}

	return logger
}

func InitLogger(config *LogConfig) error {
	// Criar diretório de logs se não existir
	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório de logs: %v", err)
	}

	// Salvar configuração para rotação
	currentConfig = config
	currentDate = time.Now().Format("20060102")

	// Configurar nível de log
	level := zap.InfoLevel
	if config.LogLevel != "" {
		var err error
		level, err = zapcore.ParseLevel(config.LogLevel)
		if err != nil {
			return fmt.Errorf("nível de log inválido: %v", err)
		}
	}

	// Nome do arquivo de log com data
	logFileName := getLogFileName()
	logFilePath := filepath.Join(config.LogDir, logFileName)

	// Fechar arquivo anterior se existir
	if currentLogFile != nil {
		currentLogFile.Close()
	}

	// Abrir arquivo de log
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo de log: %v", err)
	}
	currentLogFile = logFile

	// Configurar encoders
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02T15:04:05.000Z07:00"))
	}

	// Criar core do arquivo
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
	fileCore := zapcore.NewCore(
		fileEncoder,
		zapcore.AddSync(logFile),
		level,
	)

	// Array de cores para o logger
	var cores []zapcore.Core
	cores = append(cores, fileCore)

	// Adicionar saída para console se configurado
	if config.ConsoleOutput {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// Criar logger com todos os cores
	core := zapcore.NewTee(cores...)
	logger = zap.New(core)

	// Executar limpeza de logs antigos em goroutine
	go cleanOldLogs(config.LogDir, config.RetentionDays)

	// Iniciar rotação automática diária
	go rotateDailyLog()

	return nil
}

// rotateDailyLog verifica a cada minuto se o dia mudou e rotaciona o log
func rotateDailyLog() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		newDate := time.Now().Format("20060102")

		// Verificar se o dia mudou
		if newDate != currentDate && currentConfig != nil {
			fmt.Printf("[LOG ROTATION] Dia alterado de %s para %s, rotacionando arquivo de log...\n", currentDate, newDate)

			// Sincronizar logs pendentes
			if logger != nil {
				logger.Sync()
			}

			// Reinicializar logger com novo arquivo
			if err := InitLogger(currentConfig); err != nil {
				fmt.Printf("[LOG ROTATION] Erro ao rotacionar log: %v\n", err)
			} else {
				fmt.Printf("[LOG ROTATION] Arquivo de log rotacionado com sucesso: icrmsenderemail_%s.log\n", newDate)
			}
		}
	}
}

func cleanOldLogs(logDir string, retentionDays int) {
	if retentionDays <= 0 {
		return
	}

	for {
		files, err := os.ReadDir(logDir)
		if err != nil {
			fmt.Printf("erro ao ler diretório de logs: %v\n", err)
			time.Sleep(24 * time.Hour)
			continue
		}

		threshold := time.Now().AddDate(0, 0, -retentionDays)

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			info, err := file.Info()
			if err != nil {
				continue
			}

			if info.ModTime().Before(threshold) {
				filePath := filepath.Join(logDir, file.Name())
				if err := os.Remove(filePath); err != nil {
					fmt.Printf("erro ao remover log antigo %s: %v\n", filePath, err)
				}
			}
		}

		// Verificar novamente em 24 horas
		time.Sleep(24 * time.Hour)
	}
}

func GetLogger() *zap.Logger {
	if logger == nil {
		// Configuração padrão se o logger não foi inicializado
		config := &LogConfig{
			LogDir:        "log",
			ConsoleOutput: true,
			LogLevel:      "info",
			RetentionDays: 30,
		}
		if err := InitLogger(config); err != nil {
			// Em caso de erro, criar um logger básico para console
			logger, _ = zap.NewProduction()
		}
	}
	return logger
}

func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}
