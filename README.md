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

## 🚀 Installation

### Prerequisites

- Go 1.21 or higher
- OpenShift CLI (`oc`) or Kubernetes CLI (`kubectl`)
- OpenAI API key

### Quick Install

```bash
# Clone repository
git clone https://github.com/your-repo/oc-ai.git
cd oc-ai

# Build
go build -o oc-ai

# Set OpenAI API key
export OPENAI_API_KEY="your-api-key"

# Optional: Install globally
sudo mv oc-ai /usr/local/bin/
```

## 📋 Detailed Usage Guide

### Basic Commands

```bash
# Natural language command generation
oc-ai ai "list all pods that have crashed in the last hour"
oc-ai ai "scale the frontend deployment to 3 replicas"

# Command explanation
oc-ai explain "oc delete pod --force --grace-period=0"

# Interactive mode
oc-ai interactive

# View command history
oc-ai history
```

### Safety Levels Explained

Every command is assigned a safety level from 1 to 5:

| Level | Risk | Description | Example | Behavior |
|-------|------|-------------|----------|-----------|
| 1 | Safe | Read-only operations | `get pods` | Auto-execute |
| 2 | Low | Non-destructive changes | `label pod` | Auto-execute |
| 3 | Medium | Resource modifications | `scale deployment` | Confirm |
| 4 | High | Resource deletion | `delete pod` | Warning + Confirm |
| 5 | Critical | Cluster-wide impact | `delete namespace` | Double confirm |

### Interactive Mode Features

```bash
$ oc-ai interactive
🤖 OC-AI Interactive Shell

> show pods with high memory usage
Command: kubectl top pods --sort-by=memory
Safety: 1/5 (Safe - Read-only)
Execute? [Y/n]: y

NAME         CPU    MEMORY
pod-1        120m   1.2Gi
pod-2        85m    800Mi

> restart the pod-1 pod
Command: kubectl delete pod pod-1
Safety: 4/5 (High - Pod will be deleted)
⚠️  Warning: This is a destructive operation
Execute? [y/N]: n

> exit
Goodbye! 👋
```

### Template Management

Create and manage reusable command templates:

```bash
# List available templates
oc-ai template list

# View template details
oc-ai template show deploy-app

# Execute template with parameters
oc-ai template run deploy-app \
  --name=myapp \
  --image=nginx:1.21 \
  --replicas=3 \
  --namespace=production
```

### Configuration

Create `~/.config/oc-ai/config.yaml`:

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

### Error Handling

The tool provides clear error messages and handling:

```bash
# Invalid safety level
Error: invalid safety level "6": safety level must be between 1 and 5

# Missing required template parameter
Error: parameter "name" is required for template "deploy-app"

# History file access error
Warning: Failed to save command to history: permission denied

# Command parsing error
Error: unterminated quoted string in command
```

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
# Development build
go build -o oc-ai

# Production build with version
VERSION=$(git describe --tags)
go build -ldflags="-X main.Version=$VERSION" -o oc-ai

# Run tests
go test ./...
```

### Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

