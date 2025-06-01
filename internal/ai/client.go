package ai

import (
    "context"
    "fmt"
    "strings"

    "github.com/sashabaranov/go-openai"
)

type Client struct {
    client *openai.Client
    tool   string
    model  string
}

func NewClient(apiKey, tool, model string) *Client {
    return &Client{
        client: openai.NewClient(apiKey),
        tool:   tool,
        model:  model,
    }
}

func (c *Client) GenerateCommand(prompt string, ctx map[string]string) (string, string, string, error) {
    systemPrompt := fmt.Sprintf(`You are an expert %s administrator.
Current Context:
- Cluster: %s
- Namespace: %s
- User: %s
- Server: %s

Rules:
1. Respond ONLY with: COMMAND|||EXPLANATION|||SAFETY_LEVEL(1-5)
2. SAFETY_LEVEL: 1=Safe, 3=Caution, 5=Dangerous
3. Generate %s commands only
4. Include all required flags
5. Never include destructive commands without confirmation`,
        strings.ToUpper(c.tool), ctx["cluster"], ctx["namespace"], ctx["user"], ctx["server"], c.tool)

    resp, err := c.client.CreateChatCompletion(
        context.Background(),
        openai.ChatCompletionRequest{
            Model: c.model,
            Messages: []openai.ChatCompletionMessage{
                {
                    Role:    openai.ChatMessageRoleSystem,
                    Content: systemPrompt,
                },
                {
                    Role:    openai.ChatMessageRoleUser,
                    Content: prompt,
                },
            },
            Temperature: 0.3,
        },
    )

    if err != nil {
        return "", "", "", fmt.Errorf("AI error: %w", err)
    }

    if len(resp.Choices) == 0 {
        return "", "", "", fmt.Errorf("no response from AI")
    }

    response := resp.Choices[0].Message.Content
    parts := strings.Split(response, "|||")
    if len(parts) < 3 {
        return "", "", "", fmt.Errorf("invalid response format")
    }

    return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), strings.TrimSpace(parts[2]), nil
}

func (c *Client) ExplainCommand(command string) (string, error) {
    resp, err := c.client.CreateChatCompletion(
        context.Background(),
        openai.ChatCompletionRequest{
            Model: c.model,
            Messages: []openai.ChatCompletionMessage{
                {
                    Role:    openai.ChatMessageRoleSystem,
                    Content: "Explain this command in simple terms. Include potential risks if any.",
                },
                {
                    Role:    openai.ChatMessageRoleUser,
                    Content: command,
                },
            },
            Temperature: 0.7,
        },
    )

    if err != nil {
        return "", err
    }

    if len(resp.Choices) == 0 {
        return "", fmt.Errorf("no response from AI")
    }

    return resp.Choices[0].Message.Content, nil
}