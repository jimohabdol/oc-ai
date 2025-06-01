package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"oc-ai/internal/ai"

	"github.com/spf13/cobra"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Start interactive AI session",
	RunE: func(cmd *cobra.Command, args []string) error {
		aiClient := ai.NewClient(cfg.OpenAIKey, activeTool, cmd.Flag("ai-model").Value.String())

		fmt.Println("Entering interactive mode. Type 'exit' to quit.")
		fmt.Printf("Using %s with %s\n\n", activeTool, cmd.Flag("ai-model").Value.String())

		for {
			// Get current context
			ctx, err := cliClient.GetContext()
			if err != nil {
				fmt.Printf("Warning: Couldn't get context: %v\n", err)
				ctx = make(map[string]string)
			}

			// Show current context
			if ns := ctx["namespace"]; ns != "" {
				fmt.Printf("[%s/%s] ", ctx["cluster"], ns)
			}

			// Read input
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "exit" || input == "quit" {
				break
			}

			if input == "" {
				continue
			}

			// Process command
			command, explanation, safety, err := aiClient.GenerateCommand(input, ctx)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

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
				output, err := cliClient.Execute(command)
				if err != nil {
					fmt.Printf("Error: %v\nOutput: %s\n", err, output)
				} else if output != "" {
					fmt.Println("Output:")
					fmt.Println(output)
				}
			case "r":
				// Command revision
				fmt.Print("Enter revised command: ")
				revised, _ := reader.ReadString('\n')
				revised = strings.TrimSpace(revised)

				output, err := cliClient.Execute(revised)
				if err != nil {
					fmt.Printf("Error: %v\nOutput: %s\n", err, output)
				} else if output != "" {
					fmt.Println("Output:")
					fmt.Println(output)
				}
			default:
				fmt.Println("Command not executed")
			}

			fmt.Println()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
}
