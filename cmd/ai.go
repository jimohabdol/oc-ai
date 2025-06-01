package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"oc-ai/internal/ai"

	"github.com/spf13/cobra"
)

var aiCmd = &cobra.Command{
	Use:   "ai [prompt]",
	Short: "Generate and execute commands from natural language",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		prompt := strings.Join(args, " ")
		aiClient := ai.NewClient(cfg.OpenAIKey, activeTool, cmd.Flag("ai-model").Value.String())

		// Get current context
		ctx, err := cliClient.GetContext()
		if err != nil {
			return fmt.Errorf("failed to get cluster context: %w", err)
		}

		// Generate command
		command, explanation, safety, err := aiClient.GenerateCommand(prompt, ctx)
		if err != nil {
			return err
		}

		// Remove duplicate CLI tool name if present
		if strings.HasPrefix(command, activeTool+" ") {
			command = strings.TrimPrefix(command, activeTool+" ")
		}

		fmt.Printf("\nExplanation: %s\n", explanation)
		fmt.Printf("Safety Level: %s/5\n", safety)
		fmt.Printf("Command: %s %s\n\n", activeTool, command)

		// Safety confirmation
		safetyLevel, err := strconv.Atoi(safety)
		if err != nil {
			return fmt.Errorf("invalid safety level %q: %w", safety, err)
		}

		if safetyLevel < 1 || safetyLevel > 5 {
			return fmt.Errorf("safety level must be between 1 and 5, got %d", safetyLevel)
		}

		if cmd.Flag("yes").Value.String() == "false" && safetyLevel >= 3 {
			fmt.Printf("⚠️ Warning: This command may be destructive (Safety Level: %d/5)\n", safetyLevel)
			fmt.Print("Confirm execution? [y/N]: ")
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			if strings.ToLower(strings.TrimSpace(response)) != "y" {
				fmt.Println("Command cancelled")
				return nil
			}
		}

		// Dry run check
		if cmd.Flag("dry-run").Value.String() == "true" {
			fmt.Println("Dry run - command not executed")
			return nil
		}

		// Execute command
		output, err := cliClient.Execute(command)
		if err != nil {
			return fmt.Errorf("error executing command: %v\nOutput: %s", err, output)
		}

		if output != "" {
			fmt.Println("Command output:")
			fmt.Println(output)
		}

		// Add to history
		if historyCmd := findHistoryCommand(); historyCmd != nil {
			if err := historyCmd.AddToHistory(command); err != nil {
				fmt.Printf("Warning: Failed to save command to history: %v\n", err)
			}
		}

		return nil
	},
}

var historyManager *HistoryCommand

func init() {
	rootCmd.AddCommand(aiCmd)
	var err error
	historyManager, err = NewHistoryCommand()
	if err != nil {
		fmt.Printf("Warning: Failed to initialize history: %v\n", err)
	}
}

func findHistoryCommand() *HistoryCommand {
	return historyManager
}
