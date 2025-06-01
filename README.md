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

## 📋 Use Cases and Examples

### 1. Basic Resource Management

```bash
# List resources
oc-ai ai "show all pods in the development namespace"
> Command: oc get pods -n development
> Safety: 1/5 (Safe - Read-only)

# Resource creation
oc-ai ai "create a new deployment named frontend with nginx image and 3 replicas"
> Command: oc create deployment frontend --image=nginx --replicas=3
> Safety: 3/5 (Medium - Resource modification)

# Resource modification
oc-ai ai "scale the frontend deployment to 5 replicas"
> Command: oc scale deployment frontend --replicas=5
> Safety: 3/5 (Medium - Resource modification)
```

### 2. Troubleshooting and Debugging

```bash
# Pod issues
oc-ai ai "show logs from crashed pods in the last hour"
> Command: oc get pods --field-selector=status.phase=Failed --sort-by=.status.startTime | tail
> Safety: 1/5 (Safe - Read-only)

# Resource issues
oc-ai ai "find pods that are pending due to insufficient resources"
> Command: oc get pods --field-selector=status.phase=Pending -o json | jq '.items[] | select(.status.conditions[] | select(.reason=="Unschedulable"))'
> Safety: 1/5 (Safe - Read-only)
```

### 3. Configuration Management

```bash
# ConfigMaps and Secrets
oc-ai ai "create a configmap from file config.properties"
> Command: oc create configmap app-config --from-file=config.properties
> Safety: 2/5 (Low - Resource creation)

# Resource quotas
oc-ai ai "show resource quotas in dev namespace"
> Command: oc get resourcequota -n dev
> Safety: 1/5 (Safe - Read-only)
```

### 4. Security and Access Control

```bash
# RBAC management
oc-ai ai "create role for read-only access to pods"
> Command: oc create role pod-reader --verb=get,list,watch --resource=pods
> Safety: 3/5 (Medium - Security modification)

# Context switching
oc-ai ai "switch to production context"
> Command: oc config use-context production
> Safety: 2/5 (Low - Context change)
```

### 5. Advanced Operations

```bash
# Rolling updates
oc-ai ai "rollout new version of frontend deployment"
> Command: oc rollout restart deployment/frontend
> Safety: 3/5 (Medium - Service impact)

# Port forwarding
oc-ai ai "forward local port 8080 to service frontend"
> Command: oc port-forward svc/frontend 8080:80
> Safety: 1/5 (Safe - Local only)
```

### 6. Using Multiple Kubeconfig Files

```bash
# Using specific kubeconfig
oc-ai --kubeconfig=/path/to/kubeconfig ai "list all pods"
> Command: oc get pods --all-namespaces
> Safety: 1/5 (Safe - Read-only)

# Interactive mode
oc-ai interactive
> show pods with high memory usage
Command: oc adm top pods --sort-by=memory
Safety: 1/5 (Safe - Read-only)
Execute? [y/N/r]: y
```

### 7. Interactive Mode

```bash
# Start interactive session
oc-ai interactive

# Example session:
> show pods with high memory usage
Command: kubectl top pods --sort-by=memory
Safety: 1/5 (Safe - Read-only)
Execute? [Y/n]: y

> scale frontend deployment
Command: kubectl scale deployment frontend --replicas=3
Safety: 3/5 (Medium - Resource modification)
Execute? [Y/n/r]: r
Enter revised command: kubectl scale deployment frontend --replicas=5

> exit
```

### 8. Template Management

```bash
# List available templates
oc-ai template list

# Show template details
oc-ai template show deploy-app

# Run template with parameters
oc-ai template run deploy-app --name=myapp --replicas=3
```

## 🔧 Configuration

### Configuration File Locations

The configuration file (`config.yaml`) can be placed in:

- Linux/macOS: `~/.config/oc-ai/config.yaml`
- Windows: `%APPDATA%\oc-ai\config.yaml`
- Current directory: `./config.yaml`

### Example Configuration
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

# Windows PowerShell
$env:OPENAI_API_KEY="your-api-key"
$env:KUBECONFIG="C:\path\to\kubeconfig"
$env:OC_AI_DEFAULT_MODEL="gpt-4-turbo"
```

## 🛡️ Safety Levels

Every command is assigned a safety level from 1 to 5:

| Level | Risk | Description | Example | Behavior |
|-------|------|-------------|----------|-----------|
| 1 | Safe | Read-only operations | `get pods` | Auto-execute |
| 2 | Low | Non-destructive changes | `label pod` | Auto-execute |
| 3 | Medium | Resource modifications | `scale deployment` | Confirm |
| 4 | High | Resource deletion | `delete pod` | Warning + Confirm |
| 5 | Critical | Cluster-wide impact | `delete namespace` | Double confirm |

## 🔍 Debugging Tips

1. Use `--dry-run` flag to see commands without executing them:
```bash
oc-ai --dry-run ai "scale frontend deployment to 3 replicas"
```

2. Use `-y` flag to auto-confirm commands (use with caution):
```bash
oc-ai -y ai "restart all pods in namespace"
```

3. Check command history:
```bash
oc-ai history
```

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

