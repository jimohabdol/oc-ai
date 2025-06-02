package ai

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/sashabaranov/go-openai"
)

type Client struct {
	client *openai.Client
	tool   string
	model  string
	cache  *promptCache
}

type promptCache struct {
	sync.RWMutex
	responses map[string]cachedResponse
}

type cachedResponse struct {
	command     string
	explanation string
	safety      string
	timestamp   time.Time
}

const cacheDuration = 5 * time.Minute

func NewClient(apiKey, tool, model string) *Client {
	return &Client{
		client: openai.NewClient(apiKey),
		tool:   tool,
		model:  model,
		cache: &promptCache{
			responses: make(map[string]cachedResponse),
		},
	}
}

func (c *Client) GenerateCommand(prompt string, ctx map[string]string) (string, string, string, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("%s:%s:%s:%s:%s",
		prompt,
		ctx["cluster"],
		ctx["namespace"],
		ctx["user"],
		ctx["server"])

	if resp, ok := c.cache.get(cacheKey); ok {
		return resp.command, resp.explanation, resp.safety, nil
	}

	// Set timeout for API call
	apiCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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
		apiCtx,
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

	command := strings.TrimSpace(parts[0])
	explanation := strings.TrimSpace(parts[1])
	safety := strings.TrimSpace(parts[2])

	// Cache the response
	c.cache.set(cacheKey, cachedResponse{
		command:     command,
		explanation: explanation,
		safety:      safety,
		timestamp:   time.Now(),
	})

	return command, explanation, safety, nil
}

func (c *Client) ExplainCommand(command string) (string, error) {
	// Set timeout for API call
	apiCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.client.CreateChatCompletion(
		apiCtx,
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

func (pc *promptCache) get(key string) (cachedResponse, bool) {
	pc.RLock()
	defer pc.RUnlock()

	if resp, ok := pc.responses[key]; ok {
		if time.Since(resp.timestamp) < cacheDuration {
			return resp, true
		}
		delete(pc.responses, key)
	}
	return cachedResponse{}, false
}

func (pc *promptCache) set(key string, resp cachedResponse) {
	pc.Lock()
	defer pc.Unlock()

	pc.responses[key] = resp

	// Clean up old entries
	now := time.Now()
	for k, v := range pc.responses {
		if now.Sub(v.timestamp) > cacheDuration {
			delete(pc.responses, k)
		}
	}
}
