# cmt

A Go-based tool for generating conventional commit messages using AI models via Ollama.

## Project Structure

This project follows Go community conventions and best practices:

```
cmt/
├── cmd/cmt/           # Main application entry point
├── internal/          # Private application code
│   ├── commit/        # Commit message generation logic
│   ├── config/        # Configuration management
│   ├── git/           # Git operations
│   ├── ollama/        # Ollama integration
│   └── spinner/       # Loading spinner utilities
├── docs/              # Documentation
├── tests/             # Integration tests
├── build/             # Build artifacts (generated)
├── Makefile           # Build configuration
└── .goreleaser.yml    # Release configuration
```

## Features

- Generate conventional commit messages from staged Git changes
- Uses AI models through Ollama
- Configurable model selection
- Follows Unix conventions
- Configuration stored in `~/.cmt.yaml`

## Installation

### Prerequisites

- Go 1.21 or later
- [Ollama](https://ollama.ai/) installed and running
- Git repository

### Build

```bash
git clone https://github.com/dakoctba/cmt.git
cd cmt
go build -o cmt ./cmd/cmt
```

### Install globally (optional)

```bash
# On macOS/Linux
sudo cp cmt /usr/local/bin/

# Or add to your PATH
export PATH=$PATH:$(pwd)
```

## Usage

### Basic usage

```bash
# Stage your changes first
git add .

# Generate commit message
cmt
```

### With specific model

```bash
cmt --model llama3.1
```

### Configuration

The tool automatically creates a configuration file at `~/.cmt.yaml` on first run:

```yaml
model: llama3.1
```

You can edit this file to change the default model.

### Available flags

- `--model`: Specify the model to use (overrides config)
- `--config`: Specify a custom config file path
- `--help`: Show help message
- `--version`: Show version information

## Development

### Running tests

```bash
# Run all tests
go test ./...

# Run specific test package
go test ./internal/commit
go test ./internal/config
go test ./internal/git
go test ./internal/spinner

# Run integration tests
go test ./tests
```

### Building

```bash
# Build for current platform
go build -o cmt ./cmd/cmt

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o cmt ./cmd/cmt
```

## How it works

1. Checks if Ollama is installed and running
2. Verifies you're in a Git repository
3. Gets staged changes using `git diff --cached`
4. Shows an animated loading spinner while the AI model processes
5. Sends the diff to the specified AI model via Ollama
6. Generates a conventional commit message
7. Displays the generated commit command

## Conventional Commits

The tool generates commit messages following the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<optional scope>): <short description>

[optional body]

[optional footer(s)]
```

Common types:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation-only changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Routine tasks

## License

MIT License - see LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

For detailed documentation, see the [docs/](docs/) directory.
