package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"oc-ai/internal/ai"

	"github.com/spf13/cobra"
)

var explainCmd = &cobra.Command{
	Use:   "explain [command]",
	Short: "Explain what an oc/kubectl command does",
	Long: `Explain what a command does in the current context.
Provides detailed information about command behavior, risks, and alternatives.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		command := strings.Join(args, " ")

		// Remove duplicate CLI tool name if present
		if strings.HasPrefix(command, activeTool+" ") {
			command = strings.TrimPrefix(command, activeTool+" ")
		}

		aiClient := ai.NewClient(cfg.OpenAIKey, activeTool, cmd.Flag("ai-model").Value.String())

		// Get current context for more accurate explanation
		ctx, err := cliClient.GetContext()
		if err != nil {
			fmt.Printf("Warning: Could not get cluster context: %v\n", err)
			ctx = make(map[string]string)
		}

		explanation, err := aiClient.ExplainCommand(command)
		if err != nil {
			return fmt.Errorf("failed to explain command: %w", err)
		}

		// Create formatted output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "\nCommand Explanation: %s %s\n", activeTool, command)
		if ctx["cluster"] != "" {
			fmt.Fprintf(w, "Context:\t%s\n", ctx["cluster"])
		}
		if ctx["namespace"] != "" {
			fmt.Fprintf(w, "Namespace:\t%s\n", ctx["namespace"])
		}
		fmt.Fprintf(w, "\n%s\n", explanation)
		w.Flush()

		return nil
	},
}

func init() {
	explainCmd.Flags().Bool("verbose", false, "Show detailed explanation with examples")
	rootCmd.AddCommand(explainCmd)
}
