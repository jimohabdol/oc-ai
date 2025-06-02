package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"oc-ai/internal/cli"
	"time"
)

type ClusterContext struct {
	Cluster   string `json:"cluster"`
	Namespace string `json:"namespace"`
	User      string `json:"user"`
	Server    string `json:"server"`
}

type ContextManager struct {
	cli   cli.CLI
	cache ClusterContext
}

func NewContextManager(cli cli.CLI) *ContextManager {
	return &ContextManager{cli: cli}
}

func (cm *ContextManager) GetCurrentContext() ClusterContext {
	cm.UpdateContext()
	return cm.cache
}

func (cm *ContextManager) UpdateContext() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	output, err := cm.cli.Execute("config view -o json")
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Printf("Warning: Context update timed out after 5s\n")
		}
		return
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

	if err := json.Unmarshal([]byte(output), &config); err != nil {
		fmt.Printf("Warning: Failed to parse context: %v\n", err)
		return
	}

	for _, ctx := range config.Contexts {
		if ctx.Name == config.CurrentContext {
			// Find cluster server
			var server string
			for _, cluster := range config.Clusters {
				if cluster.Name == ctx.Context.Cluster {
					server = cluster.Cluster.Server
					break
				}
			}

			cm.cache = ClusterContext{
				Cluster:   ctx.Context.Cluster,
				Namespace: ctx.Context.Namespace,
				User:      ctx.Context.User,
				Server:    server,
			}
			break
		}
	}
}
