package config

import (
	"fmt"
	"time"

	"github.com/go-ini/ini"
)

// Config armazena todas as configurações da aplicação
type Config struct {
	Database    DatabaseConfig
	Email       EmailConfig
	Logger      LoggerConfig
	Health      HealthConfig
	Performance PerformanceConfig
	Dashboard   DashboardConfig
}

// DatabaseConfig configurações do banco de dados
type DatabaseConfig struct {
	Username string
	Password string
	TNS      string
}

// EmailConfig configurações do provedor Email
type EmailConfig struct {
	Provider string // mock, smtp, sendgrid, zenvia, pontaltech

	// SMTP genérico
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPUseTLS   bool

	// SendGrid
	SendGridAPIKey string

	// Zenvia
	ZenviaAPIToken string

	// Pontaltech
	PontaltechUsername    string
	PontaltechPassword    string
	PontaltechAccountID   int
	PontaltechAPIURL      string // URL customizada da API (opcional)
	PontaltechCallbackURL string // URL de callback para notificações (opcional)

	// Comum a todos
	DefaultFrom   string // Remetente padrão
	MaxRetries    int
	RetryInterval time.Duration
}

// LoggerConfig configurações de logging
type LoggerConfig struct {
	LogDir        string
	ConsoleOutput bool
	LogLevel      string
	RetentionDays int
	VerboseMode   bool
}

// HealthConfig configurações do health check
type HealthConfig struct {
	Enabled  bool
	HTTPPort int
}

// PerformanceConfig configurações de performance
type PerformanceConfig struct {
	WorkerCount                  int
	BatchSize                    int
	FetchIntervalSeconds         int
	SendTimeoutSeconds           int
	RetryAttempts                int
	EnableBatching               bool
	EmailRateLimitPerMin         int
	CircuitBreakerThreshold      int
	CircuitBreakerTimeoutSeconds int
	DataDisparoOffset            int // Offset em dias para filtro de DATA_AGENDAMENTO
	MaxTentativas                int // Número máximo de tentativas de envio por Email
}

// DashboardConfig configurações do dashboard
type DashboardConfig struct {
	EnableDashboard bool
	DashboardPort   int
}

// LoadConfig carrega configurações do arquivo INI
func LoadConfig(path string) (*Config, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar arquivo de configuração: %w", err)
	}

	config := &Config{}

	// Database
	dbSection := cfg.Section("oracle")
	config.Database = DatabaseConfig{
		Username: dbSection.Key("username").String(),
		Password: dbSection.Key("password").String(),
		TNS:      dbSection.Key("tns").String(),
	}

	// Email
	emailSection := cfg.Section("email")
	config.Email = EmailConfig{
		Provider: emailSection.Key("provider").MustString("mock"),

		// SMTP
		SMTPHost:     emailSection.Key("smtp_host").String(),
		SMTPPort:     emailSection.Key("smtp_port").MustInt(587),
		SMTPUsername: emailSection.Key("smtp_username").String(),
		SMTPPassword: emailSection.Key("smtp_password").String(),
		SMTPUseTLS:   emailSection.Key("smtp_use_tls").MustBool(true),

		// SendGrid
		SendGridAPIKey: emailSection.Key("sendgrid_api_key").String(),

		// Zenvia
		ZenviaAPIToken: emailSection.Key("zenvia_api_token").String(),

		// Pontaltech
		PontaltechUsername:    emailSection.Key("pontaltech_username").String(),
		PontaltechPassword:    emailSection.Key("pontaltech_password").String(),
		PontaltechAccountID:   emailSection.Key("pontaltech_account_id").MustInt(0),
		PontaltechAPIURL:      emailSection.Key("pontaltech_api_url").String(),
		PontaltechCallbackURL: emailSection.Key("pontaltech_callback_url").String(),

		// Comum
		DefaultFrom:   emailSection.Key("default_from").MustString("noreply@example.com"),
		MaxRetries:    emailSection.Key("max_retries").MustInt(3),
		RetryInterval: time.Duration(emailSection.Key("retry_interval_seconds").MustInt(300)) * time.Second,
	}

	// Logger
	logSection := cfg.Section("logger")
	config.Logger = LoggerConfig{
		LogDir:        logSection.Key("log_dir").MustString("log"),
		ConsoleOutput: logSection.Key("console_output").MustBool(true),
		LogLevel:      logSection.Key("log_level").MustString("info"),
		RetentionDays: logSection.Key("retention_days").MustInt(30),
		VerboseMode:   logSection.Key("verbose_mode").MustBool(false),
	}

	// Health
	healthSection := cfg.Section("health")
	config.Health = HealthConfig{
		Enabled:  healthSection.Key("enable_health_check").MustBool(true),
		HTTPPort: healthSection.Key("http_port").MustInt(8081),
	}

	// Performance
	perfSection := cfg.Section("performance")
	config.Performance = PerformanceConfig{
		WorkerCount:                  perfSection.Key("worker_count").MustInt(5),
		BatchSize:                    perfSection.Key("batch_size").MustInt(20),
		FetchIntervalSeconds:         perfSection.Key("fetch_interval_seconds").MustInt(5),
		SendTimeoutSeconds:           perfSection.Key("send_timeout_seconds").MustInt(30),
		RetryAttempts:                perfSection.Key("retry_attempts").MustInt(3),
		EnableBatching:               perfSection.Key("enable_batching").MustBool(true),
		EmailRateLimitPerMin:         perfSection.Key("email_rate_limit_per_min").MustInt(300),
		CircuitBreakerThreshold:      perfSection.Key("circuit_breaker_threshold").MustInt(10),
		CircuitBreakerTimeoutSeconds: perfSection.Key("circuit_breaker_timeout_seconds").MustInt(30),
		DataDisparoOffset:            perfSection.Key("data_disparo_days_offset").MustInt(0),
		MaxTentativas:                perfSection.Key("max_tentativas").MustInt(5),
	}

	// Dashboard
	dashSection := cfg.Section("dashboard")
	config.Dashboard = DashboardConfig{
		EnableDashboard: dashSection.Key("enable_dashboard").MustBool(true),
		DashboardPort:   dashSection.Key("dashboard_port").MustInt(3101),
	}

	return config, nil
}

// Validate valida as configurações
func (c *Config) Validate() error {
	// Validar database
	if c.Database.Username == "" {
		return fmt.Errorf("database.username não pode ser vazio")
	}
	if c.Database.Password == "" {
		return fmt.Errorf("database.password não pode ser vazio")
	}
	if c.Database.TNS == "" {
		return fmt.Errorf("database.tns não pode ser vazio")
	}

	// Validar Email
	if c.Email.Provider == "" {
		return fmt.Errorf("email.provider não pode ser vazio")
	}

	// Validar performance
	if c.Performance.BatchSize <= 0 {
		return fmt.Errorf("performance.batch_size deve ser maior que 0")
	}
	if c.Performance.WorkerCount <= 0 {
		return fmt.Errorf("performance.worker_count deve ser maior que 0")
	}

	return nil
}
