package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

type HistoryCommand struct {
	filePath string
}

func NewHistoryCommand() (*HistoryCommand, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user config directory: %w", err)
	}

	dirPath := filepath.Join(configDir, "oc-ai")
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create history directory: %w", err)
	}

	return &HistoryCommand{
		filePath: filepath.Join(dirPath, "history.json"),
	}, nil
}

func (h *HistoryCommand) AddToHistory(command string) error {
	entries := h.loadHistory()
	entries = append(entries, HistoryEntry{
		Timestamp: time.Now(),
		Command:   command,
		Tool:      activeTool,
	})
	return h.saveHistory(entries)
}

func (h *HistoryCommand) loadHistory() []HistoryEntry {
	data, err := os.ReadFile(h.filePath)
	if err != nil {
		return []HistoryEntry{}
	}

	var entries []HistoryEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		fmt.Printf("Warning: Failed to parse history file: %v\n", err)
		return []HistoryEntry{}
	}
	return entries
}

func (h *HistoryCommand) saveHistory(entries []HistoryEntry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal history entries: %w", err)
	}

	if err := os.WriteFile(h.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write history file: %w", err)
	}
	return nil
}

type HistoryEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Command   string    `json:"command"`
	Tool      string    `json:"tool"`
}

func init() {
	historyCmd := &cobra.Command{
		Use:   "history",
		Short: "Show command history",
		Run: func(cmd *cobra.Command, args []string) {
			hc, err := NewHistoryCommand()
			if err != nil {
				fmt.Printf("Error initializing history: %v\n", err)
				return
			}

			entries := hc.loadHistory()
			fmt.Println("Command History:")
			for i, entry := range entries {
				fmt.Printf("%d. [%s] %s %s\n",
					i+1,
					entry.Timestamp.Format("2006-01-02 15:04:05"),
					entry.Tool,
					entry.Command)
			}
		},
	}

	rootCmd.AddCommand(historyCmd)
}
