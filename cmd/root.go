package cmd

import (
	"fmt"
	"os"
	"strings"

	"oc-ai/cmd/compat"
	"oc-ai/internal/cli"
	"oc-ai/internal/config"

	"github.com/spf13/cobra"
)

var (
	cfg        *config.Config
	cliClient  cli.CLI
	activeTool string
)

var rootCmd = &cobra.Command{
	Use:   "oc-ai",
	Short: "AI-powered wrapper for oc/kubectl",
	Long: `An intelligent wrapper for OpenShift/Kubernetes CLI that converts natural language
to commands with safety checks and interactive features. Works with both 'oc' and 'kubectl'.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		activeTool, cliClient, err = cli.DetectCLI()
		if err != nil {
			return fmt.Errorf("no oc or kubectl found in PATH")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no subcommand, pass through to underlying CLI
		if len(args) > 0 {
			output, err := cliClient.Execute(strings.Join(args, " "))
			if err != nil {
				return err
			}
			fmt.Print(output)
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Common flags
	rootCmd.PersistentFlags().BoolP("yes", "y", false, "Auto-confirm command execution")
	rootCmd.PersistentFlags().Bool("dry-run", false, "Show command without executing")
	rootCmd.PersistentFlags().String("ai-model", "gpt-4-turbo", "AI model to use")

	// Inherited flags from oc/kubectl
	rootCmd.PersistentFlags().StringP("namespace", "n", "", "Namespace to use")
	rootCmd.PersistentFlags().String("context", "", "Context to use")
	rootCmd.PersistentFlags().Bool("insecure-skip-tls-verify", false, "Skip TLS verification")
}

func addCompatibilityFlags() {
	if activeTool == "oc" {
		compat.AddOCFlags(rootCmd)
	} else {
		compat.AddKubectlFlags(rootCmd)
	}
}
