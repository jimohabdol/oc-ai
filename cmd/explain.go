package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"oc-ai/internal/ai"
)

var explainCmd = &cobra.Command{
	Use:   "explain [command]",
	Short: "Explain what an oc/kubectl command does",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		command := strings.Join(args, " ")
		aiClient := ai.NewClient(cfg.OpenAIKey, activeTool, cmd.Flag("ai-model").Value.String())
		
		explanation, err := aiClient.ExplainCommand(command)
		if err != nil {
			return fmt.Errorf("failed to explain command: %w", err)
		}

		fmt.Printf("\nCommand Explanation: %s\n\n%s\n", command, explanation)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(explainCmd)
}