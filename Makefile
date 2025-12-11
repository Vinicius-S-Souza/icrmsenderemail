# Makefile para iCRMSenderEmail
# Data de criação: 11/12/2025
# Versão: 1.0.0

# Variáveis
BINARY_NAME=icrmsenderemail.exe
BUILD_DIR=build
CMD_DIR=cmd/icrmsenderemail
LOG_DIR=log
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Informações de versão
VERSION=1.0.0
BUILD_DATE=$(shell powershell -Command "Get-Date -Format 'yyyy-MM-dd HH:mm:ss'")
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>NUL || echo "dev")

# Flags de build
LDFLAGS=-ldflags "\
	-X 'github.com/Vinicius-S-Souza/icrmsenderemail/pkg/version.Version=$(VERSION)' \
	-X 'github.com/Vinicius-S-Souza/icrmsenderemail/pkg/version.BuildDate=$(BUILD_DATE)' \
	-X 'github.com/Vinicius-S-Souza/icrmsenderemail/pkg/version.GitCommit=$(GIT_COMMIT)'"

.PHONY: all build clean test run install deps help

# Default target
all: clean build

# Build da aplicação
build:
	@echo "Building iCRMSenderEmail v$(VERSION)..."
	@if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

# Limpeza de arquivos de build
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	@if exist $(BUILD_DIR) rmdir /s /q $(BUILD_DIR)
	@echo "Clean completed"

# Executar testes
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Executar aplicação em modo desenvolvimento
run: build
	@echo "Running iCRMSenderEmail..."
	@if not exist $(LOG_DIR) mkdir $(LOG_DIR)
	$(BUILD_DIR)/$(BINARY_NAME)

# Instalar dependências
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies installed"

# Instalar como serviço do Windows
install: build
	@echo "Installing Windows service..."
	$(BUILD_DIR)/$(BINARY_NAME) install
	@echo "Service installed. Use 'make start' to start it."

# Desinstalar serviço
uninstall:
	@echo "Uninstalling Windows service..."
	$(BUILD_DIR)/$(BINARY_NAME) stop
	$(BUILD_DIR)/$(BINARY_NAME) uninstall
	@echo "Service uninstalled"

# Iniciar serviço
start:
	@echo "Starting service..."
	$(BUILD_DIR)/$(BINARY_NAME) start

# Parar serviço
stop:
	@echo "Stopping service..."
	$(BUILD_DIR)/$(BINARY_NAME) stop

# Reiniciar serviço
restart:
	@echo "Restarting service..."
	$(BUILD_DIR)/$(BINARY_NAME) restart

# Exibir versão
version: build
	@$(BUILD_DIR)/$(BINARY_NAME) version

# Criar estrutura de diretórios
dirs:
	@echo "Creating directory structure..."
	@if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)
	@if not exist $(LOG_DIR) mkdir $(LOG_DIR)
	@if not exist sql mkdir sql
	@echo "Directories created"

# Help
help:
	@echo iCRMSenderEmail Makefile Commands:
	@echo.
	@echo   make build      - Compile the application
	@echo   make clean      - Remove build artifacts
	@echo   make test       - Run tests
	@echo   make run        - Build and run in development mode
	@echo   make deps       - Download and tidy dependencies
	@echo   make install    - Install as Windows service
	@echo   make uninstall  - Uninstall Windows service
	@echo   make start      - Start Windows service
	@echo   make stop       - Stop Windows service
	@echo   make restart    - Restart Windows service
	@echo   make version    - Display version information
	@echo   make dirs       - Create directory structure
	@echo   make help       - Show this help message
