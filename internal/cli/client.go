package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type CLI interface {
	Execute(command string) (string, error)
	GetContext() (map[string]string, error)
	GetVersion() (string, error)
	Supports(feature string) bool
}

type BaseCLI struct {
	command    string
	kubeconfig string
}

const defaultTimeout = 30 * time.Second

func (c *BaseCLI) Execute(cmd string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	args := ParseCommand(cmd)
	command := exec.CommandContext(ctx, c.command, args...)

	// Set up environment
	env := os.Environ()
	if c.kubeconfig != "" {
		env = append(env, fmt.Sprintf("KUBECONFIG=%s", c.kubeconfig))
	}
	command.Env = env

	output, err := command.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return string(output), fmt.Errorf("command timed out after %s: %w", defaultTimeout, err)
		}
		return string(output), err
	}

	return string(output), nil
}

func (c *BaseCLI) GetContext() (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []string{"config", "view", "-o", "json"}
	command := exec.CommandContext(ctx, c.command, args...)

	if c.kubeconfig != "" {
		command.Env = append(os.Environ(), fmt.Sprintf("KUBECONFIG=%s", c.kubeconfig))
	}

	output, err := command.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("context fetch timed out after 5s: %w", err)
		}
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	var config struct {
		CurrentContext string `json:"current-context"`
		Contexts       []struct {
			Name    string `json:"name"`
			Context struct {
				Cluster   string `json:"cluster"`
				Namespace string `json:"namespace"`
				User      string `json:"user"`
			} `json:"context"`
		} `json:"contexts"`
		Clusters []struct {
			Name    string `json:"name"`
			Cluster struct {
				Server string `json:"server"`
			} `json:"cluster"`
		} `json:"clusters"`
	}

	if err := json.Unmarshal(output, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	result := make(map[string]string)
	for _, c := range config.Contexts {
		if c.Name == config.CurrentContext {
			result["namespace"] = c.Context.Namespace
			result["user"] = c.Context.User
			result["cluster"] = c.Context.Cluster
			break
		}
	}

	for _, cluster := range config.Clusters {
		if cluster.Name == result["cluster"] {
			result["server"] = cluster.Cluster.Server
			break
		}
	}

	return result, nil
}

func (c *BaseCLI) GetVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	command := exec.CommandContext(ctx, c.command, "version")
	if c.kubeconfig != "" {
		command.Env = append(os.Environ(), fmt.Sprintf("KUBECONFIG=%s", c.kubeconfig))
	}

	output, err := command.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("version check timed out after 5s: %w", err)
		}
		return "", err
	}
	return string(output), nil
}

func (c *BaseCLI) Supports(feature string) bool {
	return false
}
