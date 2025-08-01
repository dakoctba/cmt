# Arquitetura do Projeto cmt

Este documento descreve a arquitetura e organização do projeto `cmt` seguindo as convenções e boas práticas da comunidade Go.

## Estrutura de Diretórios

```
cmt/
├── cmd/cmt/           # Ponto de entrada da aplicação
│   ├── main.go        # Função main e configuração do CLI
│   └── main_test.go   # Testes do ponto de entrada
├── internal/          # Código interno da aplicação (não exportável)
│   ├── commit/        # Lógica de geração de mensagens de commit
│   │   ├── commit.go  # Função principal de commit
│   │   └── commit_test.go
│   ├── config/        # Gerenciamento de configuração
│   │   ├── config.go  # Inicialização e gerenciamento de config
│   │   └── config_test.go
│   └── ollama/        # Integração com Ollama
│       └── ollama.go  # Funções de comunicação com Ollama
├── pkg/               # Bibliotecas públicas (reutilizáveis)
│   ├── git/           # Operações Git
│   │   └── git.go     # Funções para verificar repo e obter diff
│   └── spinner/       # Utilitários de spinner
│       └── spinner.go # Implementação do spinner de loading
├── docs/              # Documentação
│   ├── README.md      # README original
│   ├── examples.md    # Exemplos de uso
│   └── ARCHITECTURE.md # Este arquivo
├── tests/             # Testes de integração
│   └── integration_test.go
├── build/             # Artefatos de build e configuração
│   ├── Makefile       # Makefile para build e desenvolvimento
│   └── .goreleaser.yml # Configuração do GoReleaser
├── scripts/           # Scripts de automação
│   └── build.sh       # Script de build multiplataforma
├── go.mod             # Dependências do Go
├── go.sum             # Checksums das dependências
├── README.md          # README principal
├── LICENSE            # Licença do projeto
└── .gitignore         # Arquivos ignorados pelo Git
```

## Princípios de Design

### 1. Separação de Responsabilidades

- **cmd/cmt/**: Responsável apenas pelo ponto de entrada e configuração do CLI
- **internal/**: Contém toda a lógica de negócio da aplicação
- **pkg/**: Contém utilitários genéricos que podem ser reutilizados por outros projetos

### 2. Visibilidade de Pacotes

- **internal/**: Pacotes privados, não podem ser importados por projetos externos
- **pkg/**: Pacotes públicos, podem ser importados por outros projetos
- **cmd/**: Ponto de entrada da aplicação

### 3. Organização por Funcionalidade

- **commit/**: Toda a lógica relacionada à geração de mensagens de commit
- **config/**: Gerenciamento de configuração usando Viper
- **ollama/**: Integração específica com a API do Ollama
- **git/**: Operações Git genéricas
- **spinner/**: Utilitário de UI para feedback visual

## Fluxo de Dados

```
main.go (cmd/cmt/)
    ↓
config.InitConfig() (internal/config/)
    ↓
commit.RunCommit() (internal/commit/)
    ↓
ollama.CheckInstallation() (internal/ollama/)
git.CheckRepo() (internal/git/)
git.GetStagedDiff() (internal/git/)
    ↓
spinner.New().Start() (internal/spinner/)
    ↓
ollama.GenerateCommitMessage() (internal/ollama/)
    ↓
spinner.Stop() (internal/spinner/)
```

## Convenções de Nomenclatura

### Arquivos
- `main.go`: Ponto de entrada da aplicação
- `*_test.go`: Arquivos de teste
- `*.go`: Arquivos de código fonte

### Pacotes
- Nomes em minúsculas, sem underscores
- Nomes descritivos que indicam a funcionalidade
- Um pacote por diretório

### Funções
- Nomes em PascalCase para funções públicas
- Nomes em camelCase para funções privadas
- Nomes descritivos que indicam a ação

## Testes

### Estrutura de Testes
- **Testes unitários**: Junto com o código que testam (`*_test.go`)
- **Testes de integração**: Na pasta `tests/`
- **Cobertura**: Alvo de pelo menos 80% de cobertura

### Convenções de Teste
- Nomes de teste descritivos: `TestFunctionName_Scenario_ExpectedResult`
- Uso de subtests com `t.Run()`
- Setup e cleanup adequados
- Mocks para dependências externas

## Build e Deploy

### Makefile
- `make build`: Build para a plataforma atual
- `make build-all`: Build para múltiplas plataformas
- `make test`: Executa todos os testes
- `make install`: Instala globalmente

### Scripts
- `scripts/build.sh`: Script de build automatizado
- Suporte para múltiplas plataformas (Linux, macOS, Windows)
- Inclusão de informações de versão via Git tags

## Dependências

### Principais
- **cobra**: Framework CLI
- **viper**: Gerenciamento de configuração
- **spf13**: Pacotes auxiliares

### Desenvolvimento
- **golangci-lint**: Linter
- **gosec**: Análise de segurança
- **godoc**: Geração de documentação

## Configuração

### Arquivo de Configuração
- Localização: `~/.cmt.yaml`
- Formato: YAML
- Configurações: modelo padrão do Ollama

### Variáveis de Ambiente
- Suporte para configuração via variáveis de ambiente
- Precedência: flags > variáveis de ambiente > arquivo de config > padrões

## Monitoramento e Logs

### Logs
- Uso de `fmt` para output do usuário
- Logs de erro para `os.Stderr`
- Mensagens informativas para `os.Stdout`

### Tratamento de Erros
- Erros descritivos e úteis
- Uso de `fmt.Errorf` para wrapping de erros
- Códigos de saída apropriados

## Segurança

### Considerações
- Validação de entrada do usuário
- Sanitização de dados antes de enviar para Ollama
- Uso de comandos seguros para operações Git

### Boas Práticas
- Não exposição de informações sensíveis
- Uso de timeouts para operações externas
- Tratamento adequado de erros de rede

## Performance

### Otimizações
- Uso de goroutines para operações assíncronas
- Timeouts para operações externas
- Reutilização de conexões quando possível

### Monitoramento
- Testes de performance incluídos
- Métricas de tempo de resposta
- Testes de concorrência

## Manutenibilidade

### Código Limpo
- Funções pequenas e focadas
- Nomes descritivos
- Documentação adequada
- Comentários explicativos quando necessário

### Refatoração
- Separação clara de responsabilidades
- Interfaces bem definidas
- Baixo acoplamento entre módulos
- Alta coesão dentro dos módulos
