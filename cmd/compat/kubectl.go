package compat

import "github.com/spf13/cobra"

func AddKubectlFlags(cmd *cobra.Command) {
	// Kubernetes specific flags
	cmd.Flags().Bool("validate", true, "Validate against schema")
	cmd.Flags().String("field-manager", "kubectl", "Field manager name")
	cmd.Flags().String("chunk-size", "500", "Return large lists in chunks")
	cmd.Flags().Bool("show-kind", false, "Show resource kind")
	cmd.Flags().Bool("show-labels", false, "Show labels in output")
	
	// Output flags
	cmd.Flags().StringP("output", "o", "", "Output format")
	cmd.Flags().Bool("watch", false, "Watch for changes")
}

func IsKubectlSpecificFlag(flag string) bool {
	kubeFlags := map[string]bool{
		"validate":      true,
		"field-manager": true,
		"chunk-size":    true,
		"show-kind":     true,
		"show-labels":   true,
		"watch":        true,
	}
	return kubeFlags[flag]
}