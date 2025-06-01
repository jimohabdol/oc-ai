package compat

import "github.com/spf13/cobra"

func AddOCFlags(cmd *cobra.Command) {
	// OpenShift specific flags
	cmd.Flags().Bool("insecure", false, "Ignore TLS verification")
	cmd.Flags().String("as", "", "Username to impersonate")
	cmd.Flags().String("as-group", "", "Group to impersonate")
	cmd.Flags().String("as-uid", "", "UID to impersonate")
	cmd.Flags().String("certificate-authority", "", "Path to cert authority")
	cmd.Flags().String("request-timeout", "0", "Request timeout")
	cmd.Flags().Bool("loglevel", false, "Set log level")
	
	// Project specific flags
	cmd.Flags().StringP("project", "p", "", "Project name")
	cmd.Flags().Bool("all-projects", false, "All projects")
	
	// Output flags
	cmd.Flags().StringP("output", "o", "", "Output format")
}

func IsOCSpecificFlag(flag string) bool {
	ocFlags := map[string]bool{
		"as":            true,
		"as-group":      true,
		"as-uid":        true,
		"project":       true,
		"all-projects":  true,
		"loglevel":      true,
	}
	return ocFlags[flag]
}