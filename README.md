# OC-AI: AI-Powered OpenShift & Kubernetes CLI Assistant

OC-AI enhances your OpenShift (`oc`) and Kubernetes (`kubectl`) CLI experience with AI capabilities, enabling natural language interaction with your clusters while maintaining safety and reliability.

## ✨ Key Features

- **🔄 Intelligent CLI Detection**: Automatically detects and works with both `oc` and `kubectl`
- **🧠 Natural Language Commands**: Convert English descriptions into precise CLI commands
- **🛡️ Advanced Safety System**: Multi-level risk assessment with smart confirmation handling
- **💬 Interactive Mode**: Rich interactive shell with command suggestions and explanations
- **📚 Command History**: Persistent command history with timestamps and tool tracking
- **🎯 Smart Templates**: Parameterized command templates with validation
- **🔒 Error Handling**: Comprehensive error handling and user feedback
- **⚡ Zero Latency**: Direct command passthrough for native CLI performance
- **🔑 Flexible Authentication**: Support for multiple kubeconfig files and contexts

## 🚀 Installation

### Prerequisites

- Go 1.21 or higher
- OpenShift CLI (`oc`) or Kubernetes CLI (`kubectl`)
- OpenAI API key
- Valid kubeconfig file

### Quick Install

#### Linux/macOS
```bash
# Clone repository
git clone https://github.com/jimohabdol/oc-ai.git
cd oc-ai

# Build
go build -o oc-ai

# Set OpenAI API key
export OPENAI_API_KEY="your-api-key"

# Optional: Install globally
sudo mv oc-ai /usr/local/bin/
```

#### Windows
```powershell
# Clone repository
git clone https://github.com/jimohabdol/oc-ai.git
cd oc-ai

# Build
go build -o oc-ai.exe

# Set OpenAI API key (PowerShell)
$env:OPENAI_API_KEY="your-api-key"

# Optional: Add to PATH
# 1. Create a directory for binaries
mkdir "$env:USERPROFILE\bin"
# 2. Move the executable
move oc-ai.exe "$env:USERPROFILE\bin"
# 3. Add to PATH (requires admin PowerShell)
[Environment]::SetEnvironmentVariable(
    "Path",
    [Environment]::GetEnvironmentVariable("Path", "User") + ";$env:USERPROFILE\bin",
    "User"
)
```

## 📋 Detailed Usage Guide

### Basic Commands

```bash
# Linux/macOS
oc-ai ai "list all pods that have crashed in the last hour"
oc-ai ai "scale the frontend deployment to 3 replicas"

# Windows (PowerShell/CMD)
oc-ai.exe ai "list all pods that have crashed in the last hour"
oc-ai.exe ai "scale the frontend deployment to 3 replicas"
```

### Using Custom Kubeconfig

```bash
# Linux/macOS
oc-ai --kubeconfig=/path/to/kubeconfig ai "list all pods"
export KUBECONFIG=/path/to/kubeconfig
oc-ai ai "list pods"

# Windows (PowerShell)
oc-ai.exe --kubeconfig="C:\path\to\kubeconfig" ai "list all pods"
$env:KUBECONFIG="C:\path\to\kubeconfig"
oc-ai.exe ai "list pods"

# Windows (CMD)
oc-ai.exe --kubeconfig="C:\path\to\kubeconfig" ai "list all pods"
set KUBECONFIG=C:\path\to\kubeconfig
oc-ai.exe ai "list pods"
```

### Configuration File Locations

The configuration file (`config.yaml`) can be placed in:

- Linux/macOS: `~/.config/oc-ai/config.yaml`
- Windows: `%APPDATA%\oc-ai\config.yaml`
- Current directory: `./config.yaml`

Example configuration:
```yaml
# OpenAI Settings
openai_key: "sk-..." # Or use OPENAI_API_KEY env var
default_model: "gpt-4-turbo"

# Command Settings
confirm_execute: true  # Always confirm commands
history_limit: 100     # Number of history entries to keep

# CLI Settings
preferred_cli: "auto"  # "oc", "kubectl", or "auto"
```

### Environment Variables

```bash
# Linux/macOS
export OPENAI_API_KEY="your-api-key"
export KUBECONFIG="/path/to/kubeconfig"
export OC_AI_DEFAULT_MODEL="gpt-4-turbo"

# Windows (PowerShell)
$env:OPENAI_API_KEY="your-api-key"
$env:KUBECONFIG="C:\path\to\kubeconfig"
$env:OC_AI_DEFAULT_MODEL="gpt-4-turbo"

# Windows (CMD)
set OPENAI_API_KEY=your-api-key
set KUBECONFIG=C:\path\to\kubeconfig
set OC_AI_DEFAULT_MODEL=gpt-4-turbo
```

### Interactive Mode

```bash
# Linux/macOS
oc-ai interactive

# Windows
oc-ai.exe interactive
```

Example session:
```
🤖 OC-AI Interactive Shell

> show pods with high memory usage
Command: kubectl top pods --sort-by=memory
Safety: 1/5 (Safe - Read-only)
Execute? [Y/n]: y

NAME         CPU    MEMORY
pod-1        120m   1.2Gi
pod-2        85m    800Mi

> exit
Goodbye! 👋
```

### Template Management

```bash
# Linux/macOS
oc-ai template list
oc-ai template show deploy-app
oc-ai template run deploy-app --name=myapp --replicas=3

# Windows
oc-ai.exe template list
oc-ai.exe template show deploy-app
oc-ai.exe template run deploy-app --name=myapp --replicas=3
```

### Safety Levels

Every command is assigned a safety level from 1 to 5:

| Level | Risk | Description | Example | Behavior |
|-------|------|-------------|----------|-----------|
| 1 | Safe | Read-only operations | `get pods` | Auto-execute |
| 2 | Low | Non-destructive changes | `label pod` | Auto-execute |
| 3 | Medium | Resource modifications | `scale deployment` | Confirm |
| 4 | High | Resource deletion | `delete pod` | Warning + Confirm |
| 5 | Critical | Cluster-wide impact | `delete namespace` | Double confirm |

## 🔧 Development

### Project Structure

```
oc-ai/
├── cmd/                    # Command implementations
│   ├── root.go            # Root command and initialization
│   ├── ai.go              # AI command generation
│   ├── explain.go         # Command explanation
│   ├── interactive.go     # Interactive shell
│   ├── history.go         # Command history
│   ├── template.go        # Template management
│   └── compat/           # CLI compatibility
├── internal/
│   ├── ai/               # AI integration
│   │   ├── client.go     # OpenAI client
│   │   ├── prompt.go     # Prompt templates
│   │   └── context.go    # Context management
│   ├── cli/              # CLI abstraction
│   │   ├── client.go     # Base interface
│   │   ├── detector.go   # CLI detection
│   │   └── executor.go   # Command execution
│   └── config/           # Configuration
└── README.md
```

### Building

```bash
# Linux/macOS
go build -o oc-ai

# Windows
go build -o oc-ai.exe

# Production build with version
$VERSION=$(git describe --tags)
go build -ldflags="-X main.Version=$VERSION" -o oc-ai
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.txt ./...
go tool cover -html=coverage.txt
```

### Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

