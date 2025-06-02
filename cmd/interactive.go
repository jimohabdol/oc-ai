package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"oc-ai/internal/ai"

	"github.com/spf13/cobra"
)

// CommandCache holds recently used commands and their results
type CommandCache struct {
	sync.RWMutex
	cache map[string]cachedCommand
}

type cachedCommand struct {
	result    string
	timestamp time.Time
}

var (
	cmdCache = &CommandCache{
		cache: make(map[string]cachedCommand),
	}
	cacheDuration = 5 * time.Minute
)

func (c *CommandCache) Get(key string) (string, bool) {
	c.RLock()
	defer c.RUnlock()
	if cmd, ok := c.cache[key]; ok {
		if time.Since(cmd.timestamp) < cacheDuration {
			return cmd.result, true
		}
		delete(c.cache, key)
	}
	return "", false
}

func (c *CommandCache) Set(key, value string) {
	c.Lock()
	defer c.Unlock()
	c.cache[key] = cachedCommand{
		result:    value,
		timestamp: time.Now(),
	}
}

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Start interactive AI session",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		aiClient := ai.NewClient(cfg.OpenAIKey, activeTool, cmd.Flag("ai-model").Value.String())
		var lastContext map[string]string

		fmt.Println("Entering interactive mode. Type 'exit' to quit.")
		fmt.Printf("Using %s with %s\n\n", activeTool, cmd.Flag("ai-model").Value.String())

		// Start context update goroutine
		contextChan := make(chan map[string]string)
		go func() {
			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					if ctx, err := cliClient.GetContext(); err == nil {
						contextChan <- ctx
					}
				}
			}
		}()

		// Command execution channel
		type execResult struct {
			output string
			err    error
		}
		execChan := make(chan execResult)

		reader := bufio.NewReader(os.Stdin)
		for {
			// Non-blocking context update
			select {
			case ctx := <-contextChan:
				lastContext = ctx
			default:
				if lastContext == nil {
					var err error
					lastContext, err = cliClient.GetContext()
					if err != nil {
						lastContext = make(map[string]string)
					}
				}
			}

			// Show current context
			if ns := lastContext["namespace"]; ns != "" {
				fmt.Printf("[%s/%s] ", lastContext["cluster"], ns)
			}

			// Read input
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "exit" || input == "quit" {
				break
			}

			if input == "" {
				continue
			}

			// Process command concurrently
			commandChan := make(chan string)
			explanationChan := make(chan string)
			safetyChan := make(chan string)
			errChan := make(chan error)

			go func() {
				command, explanation, safety, err := aiClient.GenerateCommand(input, lastContext)
				if err != nil {
					errChan <- err
					return
				}
				commandChan <- command
				explanationChan <- explanation
				safetyChan <- safety
			}()

			// Wait for command generation with timeout
			select {
			case err := <-errChan:
				fmt.Printf("Error: %v\n", err)
				continue
			case command := <-commandChan:
				explanation := <-explanationChan
				safety := <-safetyChan

				// Remove duplicate CLI tool name if present
				if strings.HasPrefix(command, activeTool+" ") {
					command = strings.TrimPrefix(command, activeTool+" ")
				}

				fmt.Printf("\nCommand: %s %s\n", activeTool, command)
				fmt.Printf("Explanation: %s\n", explanation)
				fmt.Printf("Safety Level: %s/5\n", safety)

				// Get confirmation
				fmt.Print("\nExecute? [y/N/r (run/revise)]: ")
				response, _ := reader.ReadString('\n')
				response = strings.ToLower(strings.TrimSpace(response))

				switch response {
				case "y":
					// Check cache first
					if cachedOutput, found := cmdCache.Get(command); found {
						fmt.Println("Output (cached):")
						fmt.Println(cachedOutput)
						continue
					}

					// Execute command concurrently
					go func() {
						output, err := cliClient.Execute(command)
						execChan <- execResult{output, err}
					}()

					// Wait for execution
					result := <-execChan
					if result.err != nil {
						fmt.Printf("Error: %v\nOutput: %s\n", result.err, result.output)
					} else if result.output != "" {
						fmt.Println("Output:")
						fmt.Println(result.output)
						cmdCache.Set(command, result.output)
					}

					// Asynchronously save to history
					go func(cmd string) {
						if historyCmd := findHistoryCommand(); historyCmd != nil {
							if err := historyCmd.AddToHistory(cmd); err != nil {
								fmt.Printf("Warning: Failed to save command to history: %v\n", err)
							}
						}
					}(command)

				case "r":
					fmt.Print("Enter revised command: ")
					revised, _ := reader.ReadString('\n')
					revised = strings.TrimSpace(revised)

					// Execute revised command concurrently
					go func() {
						output, err := cliClient.Execute(revised)
						execChan <- execResult{output, err}
					}()

					// Wait for execution
					result := <-execChan
					if result.err != nil {
						fmt.Printf("Error: %v\nOutput: %s\n", result.err, result.output)
					} else if result.output != "" {
						fmt.Println("Output:")
						fmt.Println(result.output)
					}

					// Asynchronously save to history
					go func(cmd string) {
						if historyCmd := findHistoryCommand(); historyCmd != nil {
							if err := historyCmd.AddToHistory(cmd); err != nil {
								fmt.Printf("Warning: Failed to save command to history: %v\n", err)
							}
						}
					}(revised)

				default:
					fmt.Println("Command not executed")
				}
			case <-time.After(15 * time.Second):
				fmt.Println("Command generation timed out")
				continue
			}

			fmt.Println()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
}
