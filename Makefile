# Makefile para o projeto CMT

# Variáveis
BINARY_NAME=cmt
VERSION=$(shell git tag -l --sort=-v:refname | head -n 1 | sed 's/^v//' 2>/dev/null || echo "dev")
GIT_COMMIT=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"

# Detecção do sistema operacional
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	OS := darwin
else ifeq ($(UNAME_S),Linux)
	OS := linux
else
	OS := windows
endif

# Detecção da arquitetura
UNAME_M := $(shell uname -m)
ifeq ($(UNAME_M),x86_64)
	ARCH := amd64
else ifeq ($(UNAME_M),arm64)
	ARCH := arm64
else ifeq ($(UNAME_M),aarch64)
	ARCH := arm64
else
	ARCH := 386
endif

.PHONY: all build clean test run install uninstall release snapshot help

# Alvos padrão
all: build

# Compilação local para desenvolvimento
build:
	@echo "Building $(BINARY_NAME) $(VERSION) ($(GIT_COMMIT)) for $(OS)/$(ARCH)..."
	@go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/cmt
	@echo "Build completed successfully!"

# Limpar arquivos de build
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf dist/
	@go clean
	@echo "Clean completed successfully!"

# Executar testes
test:
	@echo "Running tests..."
	@go test ./...
	@echo "Tests completed successfully!"

# Executar o programa (para desenvolvimento)
run:
	@go run $(LDFLAGS) ./cmd/cmt $(ARGS)

# Instalação local
install: build
	@echo "Installing $(BINARY_NAME)..."
	@go install $(LDFLAGS) ./cmd/cmt
	@echo "Installation completed successfully!"

# Desinstalação
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(shell which $(BINARY_NAME))
	@echo "Uninstallation completed successfully!"

# Processo interativo para criar uma nova versão e release
release:
	@echo "Iniciando processo de release com GoReleaser..."
	@if [ -f .env ]; then \
		echo "Carregando variáveis de ambiente do arquivo .env..."; \
	fi
	@echo "Executando GoReleaser..."
	@set -a; [ -f .env ] && source .env; set +a; goreleaser release --clean
	@echo "Release concluída com sucesso!"

# Criar uma release em modo snapshot (para testes)
snapshot:
	@echo "Creating snapshot release..."
	@set -a; [ -f .env ] && source .env; set +a; goreleaser release --snapshot --clean

# Exibir informações de ajuda
help:
	@echo "Makefile para $(BINARY_NAME) - Comandos disponíveis:"
	@echo ""
	@echo "  make                - Equivalente a 'make build'"
	@echo "  make build          - Compilar o projeto para o ambiente local"
	@echo "  make clean          - Remover arquivos temporários e de build"
	@echo "  make test           - Executar testes"
	@echo "  make run ARGS=\"\"   - Executar aplicação (passar argumentos em ARGS)"
	@echo "  make install        - Instalar o binário localmente"
	@echo "  make uninstall      - Desinstalar o binário"
	@echo "  make release        - Processo interativo para criar nova versão e release"
	@echo "  make snapshot       - Criar uma release de teste (não publicada)"
	@echo "  make help           - Exibir esta mensagem de ajuda"
	@echo ""
	@echo "Exemplos:"
	@echo "  make run ARGS=\"--help\""
	@echo "  make run ARGS=\"--version\""
	@echo ""
