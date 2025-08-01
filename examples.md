# Exemplos de Uso - cmt

## Cenários Comuns

### 1. Commit de uma nova funcionalidade

```bash
# Adicionar arquivos
git add .

# Gerar commit message
cmt
```

**Saída esperada:**
```
🤔 ⠋ Thinking with llama3.1 model...
Generated commit message:

git commit -m "feat: add user authentication system" -m "Implements JWT-based authentication with login/logout functionality"
```

### 2. Commit de correção de bug

```bash
# Fazer alterações para corrigir bug
git add .

# Gerar commit message
cmt
```

**Saída esperada:**
```
🤔 ⠋ Thinking with llama3.1 model...
Generated commit message:

git commit -m "fix: resolve memory leak in data processing" -m "Fixes issue where large datasets were not being properly cleaned up"
```

### 3. Usando um modelo específico

```bash
# Usar modelo diferente
cmt --model codellama

# Ou configurar permanentemente
echo "model: codellama" > ~/.cmt.yaml
```

### 4. Commit de documentação

```bash
# Adicionar documentação
git add README.md docs/

# Gerar commit message
cmt
```

**Saída esperada:**
```
🤔 ⠋ Thinking with llama3.1 model...
Generated commit message:

git commit -m "docs: update API documentation" -m "Adds comprehensive examples and usage guidelines for all endpoints"
```

### 5. Commit de refatoração

```bash
# Refatorar código
git add .

# Gerar commit message
cmt
```

**Saída esperada:**
```
🤔 ⠋ Thinking with llama3.1 model...
Generated commit message:

git commit -m "refactor: improve error handling in database layer" -m "Replaces generic error messages with specific database error types"
```

## Fluxo de Trabalho Típico

```bash
# 1. Fazer alterações no código
vim src/main.go

# 2. Verificar mudanças
git status

# 3. Adicionar arquivos ao staging
git add .

# 4. Gerar commit message
cmt

# 5. Copiar e executar o comando gerado
git commit -m "feat: implement new feature" -m "Description here"

# 6. Push para o repositório
git push
```

## Troubleshooting

### Erro: "ollama is not installed"
```bash
# Instalar Ollama
curl -fsSL https://ollama.ai/install.sh | sh

# Iniciar serviço
ollama serve
```

### Erro: "no staged changes found"
```bash
# Adicionar arquivos ao staging primeiro
git add .

# Ou adicionar arquivos específicos
git add src/main.go
```

### Erro: "this is not a Git repository"
```bash
# Inicializar repositório Git
git init

# Ou navegar para um repositório existente
cd /path/to/git/repo
```
