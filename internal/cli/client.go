package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
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

func (c *BaseCLI) Execute(cmd string) (string, error) {
	args := strings.Split(cmd, " ")

	// Set kubeconfig as environment variable if specified
	command := exec.Command(c.command, args...)
	if c.kubeconfig != "" {
		command.Env = append(os.Environ(), fmt.Sprintf("KUBECONFIG=%s", c.kubeconfig))
	}

	output, err := command.CombinedOutput()
	return string(output), err
}

func (c *BaseCLI) GetContext() (map[string]string, error) {
	args := []string{"config", "view", "-o", "json"}
	if c.kubeconfig != "" {
		args = append([]string{"--kubeconfig", c.kubeconfig}, args...)
	}
	output, err := exec.Command(c.command, args...).CombinedOutput()
	if err != nil {
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

	ctx := make(map[string]string)
	for _, c := range config.Contexts {
		if c.Name == config.CurrentContext {
			ctx["namespace"] = c.Context.Namespace
			ctx["user"] = c.Context.User
			ctx["cluster"] = c.Context.Cluster
			break
		}
	}

	for _, cluster := range config.Clusters {
		if cluster.Name == ctx["cluster"] {
			ctx["server"] = cluster.Cluster.Server
			break
		}
	}

	return ctx, nil
}

func (c *BaseCLI) GetVersion() (string, error) {
	output, err := exec.Command(c.command, "version").CombinedOutput()
	return string(output), err
}

func (c *BaseCLI) Supports(feature string) bool {
	// Implementation varies between oc and kubectl
	return false
}
