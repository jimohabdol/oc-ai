package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

type CommandTemplate struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Command     string              `json:"command"`
	Parameters  []TemplateParameter `json:"parameters"`
}

type TemplateParameter struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

func loadTemplates() ([]CommandTemplate, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(configDir, "oc-ai", "templates.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var templates []CommandTemplate
	if err := json.Unmarshal(data, &templates); err != nil {
		return nil, err
	}

	return templates, nil
}

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage command templates",
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates",
	Run: func(cmd *cobra.Command, args []string) {
		templates, err := loadTemplates()
		if err != nil {
			fmt.Printf("Error loading templates: %v\n", err)
			return
		}

		fmt.Println("Available Templates:")
		for _, t := range templates {
			fmt.Printf("  %s - %s\n", t.Name, t.Description)
		}
	},
}

var templateShowCmd = &cobra.Command{
	Use:   "show [name]",
	Short: "Show template details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		templates, err := loadTemplates()
		if err != nil {
			fmt.Printf("Error loading templates: %v\n", err)
			return
		}

		for _, t := range templates {
			if t.Name == args[0] {
				fmt.Printf("\nTemplate: %s\n", t.Name)
				fmt.Printf("Description: %s\n", t.Description)
				fmt.Printf("Command: %s\n", t.Command)
				fmt.Println("Parameters:")
				for _, p := range t.Parameters {
					fmt.Printf("  %s (required: %t) - %s\n",
						p.Name, p.Required, p.Description)
				}
				return
			}
		}
		fmt.Printf("Template '%s' not found\n", args[0])
	},
}

var templateRunCmd = &cobra.Command{
	Use:   "run [name]",
	Short: "Execute a template",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		templates, err := loadTemplates()
		if err != nil {
			return fmt.Errorf("error loading templates: %w", err)
		}

		var selected *CommandTemplate
		for _, t := range templates {
			if t.Name == args[0] {
				selected = &t
				break
			}
		}

		if selected == nil {
			return fmt.Errorf("template '%s' not found", args[0])
		}

		params := make(map[string]string)
		for _, p := range selected.Parameters {
			value, err := cmd.Flags().GetString(p.Name)
			if err != nil {
				return fmt.Errorf("error getting parameter %s: %w", p.Name, err)
			}
			if p.Required && value == "" {
				return fmt.Errorf("parameter %s is required", p.Name)
			}
			params[p.Name] = value
		}

		tmpl, err := template.New("command").Parse(selected.Command)
		if err != nil {
			return fmt.Errorf("error parsing template: %w", err)
		}

		var buf strings.Builder
		if err := tmpl.Execute(&buf, params); err != nil {
			return fmt.Errorf("error executing template: %w", err)
		}

		generatedCmd := buf.String()
		fmt.Printf("Generated command: %s %s\n", activeTool, generatedCmd)

		if cmd.Flag("dry-run").Value.String() == "true" {
			fmt.Println("Dry run - command not executed")
			return nil
		}

		output, err := cliClient.Execute(generatedCmd)
		if err != nil {
			return fmt.Errorf("error executing command: %v\nOutput: %s", err, output)
		}

		if output != "" {
			fmt.Println("Command output:")
			fmt.Println(output)
		}

		return nil
	},
}

func init() {
	// Initialize template parameters for run command
	templates, _ := loadTemplates()
	for _, t := range templates {
		for _, p := range t.Parameters {
			templateRunCmd.Flags().String(p.Name, "", p.Description)
		}
	}

	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateShowCmd)
	templateCmd.AddCommand(templateRunCmd)
	rootCmd.AddCommand(templateCmd)
}
