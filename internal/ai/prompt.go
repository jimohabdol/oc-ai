package ai

import (
	"fmt"
	"strings"
)

const (
	SystemPromptTemplate = `You are an expert %s administrator.
Current Context:
- Cluster: %s
- Namespace: %s
- User: %s
- Server: %s

Rules:
1. Respond ONLY with: COMMAND|||EXPLANATION|||SAFETY_LEVEL(1-5)
2. SAFETY_LEVEL: 1=Safe, 3=Caution, 5=Dangerous
3. Generate commands for %s but DO NOT include '%s' at the start of the command
4. Include all required flags
5. Never include destructive commands without confirmation`

	ExplainPromptTemplate = `Explain what this %s command does in simple terms. 
Include:
1. What resources it affects
2. Potential risks
3. Common use cases
4. Any safer alternatives if applicable

Command: %s`
)

func BuildSystemPrompt(tool string, ctx ClusterContext) string {
	return fmt.Sprintf(SystemPromptTemplate,
		strings.ToUpper(tool),
		ctx.Cluster,
		ctx.Namespace,
		ctx.User,
		ctx.Server,
		tool,
		tool)
}

func BuildExplainPrompt(tool, command string) string {
	return fmt.Sprintf(ExplainPromptTemplate, tool, command)
}
